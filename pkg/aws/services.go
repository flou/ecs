package aws

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/fatih/color"
)

// ListServices describes services running in the ECS cluster filtered by cluster name, service name ans service type
func ListServices(client *ecs.Client, clusterName, serviceFilter, serviceType string) []ecs.Service {
	ecsServices := []ecs.Service{}
	serviceNames := []string{}
	listServicesInput := ecs.ListServicesInput{Cluster: &clusterName}
	if strings.ToLower(serviceType) == "fargate" {
		listServicesInput.LaunchType = "FARGATE"
	} else if strings.ToLower(serviceType) == "ec2" {
		listServicesInput.LaunchType = "EC2"
	}
	req := client.ListServicesRequest(&listServicesInput)

	pager := ecs.NewListServicesPaginator(req)
	for pager.Next(context.Background()) {
		page := pager.CurrentPage()
		serviceNames = append(serviceNames, page.ServiceArns...)
	}
	filteredServicesNames := []string{}
	for _, service := range serviceNames {
		if strings.Contains(service, serviceFilter) {
			filteredServicesNames = append(filteredServicesNames, service)
		}
	}
	sort.Strings(filteredServicesNames)
	for _, services := range chunk(filteredServicesNames, 10) {
		if len(services) > 0 {
			descServices := DescribeServices(client, clusterName, services)
			ecsServices = append(ecsServices, descServices...)
		}
	}
	return ecsServices
}

// Describe a list of services running in the ECS cluster
func DescribeServices(client *ecs.Client, clusterName string, services []string) []ecs.Service {
	params := ecs.DescribeServicesInput{Cluster: &clusterName, Services: services}
	resp, err := client.DescribeServicesRequest(&params).Send(context.Background())
	if err != nil {
		fmt.Println("Failed to describe services: " + err.Error())
		os.Exit(1)
	}
	return resp.Services
}

func FindService(client *ecs.Client, cluster, service string) (ecs.Service, error) {
	var ecsService ecs.Service
	runningServices := DescribeServices(client, cluster, []string{service})
	if len(runningServices) == 0 {
		fmt.Printf("No running service %s in cluster %s\n", service, cluster)
		return ecsService, errors.New("No service in cluster")
	}
	if len(runningServices) > 1 {
		fmt.Printf("Found more than 1 service named %s in cluster %s\n", service, cluster)
		return ecsService, errors.New("No service in cluster")
	}
	return runningServices[0], nil
}

func FindResource(resources []ecs.Resource, name string) ecs.Resource {
	var resource ecs.Resource
	for _, res := range resources {
		if *res.Name == name {
			resource = res
			break
		}
	}
	return resource
}

func FindAttribute(attributes []ecs.Attribute, name string) ecs.Attribute {
	var attribute ecs.Attribute
	for _, attr := range attributes {
		if *attr.Name == name {
			attribute = attr
			break
		}
	}
	return attribute
}

func FindTag(tags []ec2.Tag, name string) ec2.Tag {
	var tag ec2.Tag
	for _, t := range tags {
		if *t.Key == name {
			tag = t
			break
		}
	}
	return tag
}

func serviceUp(service *ecs.Service) bool {
	return *service.DesiredCount == *service.RunningCount &&
		len(service.Events) > 0 &&
		strings.Contains(*service.Events[0].Message, "has reached a steady state")
}

func serviceStatus(service *ecs.Service) string {
	status := color.GreenString("[OK]")
	serviceUp := serviceUp(service)
	if !serviceUp {
		status = color.RedString("[KO]")
	}
	if serviceUp && *service.RunningCount == 0 {
		status = color.YellowString("[WARN]")
	}
	return status
}

func ServiceOk(service *ecs.Service) bool {
	status := serviceStatus(service)
	return strings.Contains(status, "OK")
}

func ServiceTaskDefinition(client *ecs.Client, taskDefinition string) ecs.TaskDefinition {
	resp, err := client.DescribeTaskDefinitionRequest(
		&ecs.DescribeTaskDefinitionInput{TaskDefinition: &taskDefinition}).Send(context.Background())
	if err != nil {
		fmt.Println("Failed to describe task definition: " + err.Error())
		os.Exit(1)
	}
	return *resp.TaskDefinition
}

func PrintServiceDetails(client *ecs.Client, service *ecs.Service, longOutput bool) {
	elbClient := elasticloadbalancingv2.New(client.Config)
	fmt.Printf(
		"%-15s  %-70s %-7s %-8s running %d/%d  (%s)\n",
		serviceStatus(service), color.YellowString(*service.ServiceName),
		service.LaunchType, *service.Status, *service.RunningCount,
		*service.DesiredCount, shortTaskDefinitionName(*service.TaskDefinition),
	)
	if longOutput == true {
		taskDefinition := ServiceTaskDefinition(client, *service.TaskDefinition)
		fmt.Println(linkToConsole(service, clusterNameFromArn(*service.ClusterArn)))
		if taskDefinition.TaskRoleArn != nil {
			fmt.Printf("IAM Role: %s\n", linkToIAM(shortTaskDefinitionName(*taskDefinition.TaskRoleArn)))
		}

		for _, lb := range service.LoadBalancers {
			response, err := elbClient.DescribeTargetGroupsRequest(&elasticloadbalancingv2.DescribeTargetGroupsInput{
				TargetGroupArns: []string{*lb.TargetGroupArn},
			}).Send(context.Background())
			if err != nil {
				fmt.Println("Failed to describe target group: " + err.Error())
				os.Exit(1)
			}
			targetGroup := response.TargetGroups[0]
			fmt.Println("Load Balancing:")
			fmt.Printf("  Target Group: %s\n", *targetGroup.TargetGroupArn)
			if lb.TargetGroupArn != nil {
				fmt.Printf("  Healthcheck: %s %s -> %s(%d)\n", targetGroup.Protocol, *targetGroup.HealthCheckPath, *targetGroup.HealthCheckPort, *targetGroup.Port)
			}
		}
		if service.NetworkConfiguration != nil {
			if service.NetworkConfiguration.AwsvpcConfiguration != nil {
				config := service.NetworkConfiguration.AwsvpcConfiguration
				if len(config.SecurityGroups) > 0 {
					fmt.Printf("Security Group: %s\n", config.SecurityGroups)
				}
				if len(config.Subnets) > 0 {
					fmt.Printf("VPC Subnets: %s\n", config.Subnets)
				}
			}
		}
		for _, container := range taskDefinition.ContainerDefinitions {
			fmt.Printf("- Container: %s\n", color.GreenString(*container.Name))
			fmt.Printf("  Image: %s\n", color.YellowString(*container.Image))
			var (
				containerMemory = "-"
				containerCPU    = "-"
			)
			if container.Cpu != nil {
				containerCPU = strconv.FormatInt(*container.Cpu, 10)
			}
			if container.Memory != nil {
				containerMemory = strconv.FormatInt(*container.Memory, 10)
			}
			fmt.Printf("  Memory: %s / CPU: %s\n", containerMemory, containerCPU)
			if len(container.PortMappings) > 0 {
				fmt.Println("  Ports:")
				for _, port := range container.PortMappings {
					fmt.Printf("   - Host:%d -> Container:%d\n", *port.HostPort, *port.ContainerPort)
				}
			}
			if container.LogConfiguration != nil {
				fmt.Println("  Logs:")
				fmt.Printf("   - log-driver: %s\n", container.LogConfiguration.LogDriver)
				for name, option := range container.LogConfiguration.Options {
					fmt.Printf("   - %s: %s\n", name, option)
				}
			}
			if len(container.Environment) > 0 {
				fmt.Println("  Environment:")
				for _, env := range container.Environment {
					fmt.Printf("   - %s: %s\n", *env.Name, *env.Value)
				}
			}
		}
		fmt.Println()
	}
}
