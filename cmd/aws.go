package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mitchellh/colorstring"
)

func findResource(resources []ecs.Resource, name string) ecs.Resource {
	var resource ecs.Resource
	for _, res := range resources {
		if *res.Name == name {
			resource = res
			break
		}
	}
	return resource
}

func findAttribute(attributes []ecs.Attribute, name string) ecs.Attribute {
	var attribute ecs.Attribute
	for _, attr := range attributes {
		if *attr.Name == name {
			attribute = attr
			break
		}
	}
	return attribute
}

func findTag(tags []ec2.Tag, name string) ec2.Tag {
	var tag ec2.Tag
	for _, t := range tags {
		if *t.Key == name {
			tag = t
			break
		}
	}
	return tag
}

func shortTaskDefinitionName(taskDefinition string) string {
	splitTaskDefinition := strings.Split(taskDefinition, "/")
	return strings.Split(taskDefinition, "/")[len(splitTaskDefinition)-1]
}

func clusterNameFromArn(clusterArn string) string {
	splitClusterArn := strings.Split(clusterArn, "/")
	return strings.Split(clusterArn, "/")[len(splitClusterArn)-1]
}

func listClusters(client *ecs.ECS, filter string) []string {
	clusterNames := []string{}
	listClusterOutput, err := client.ListClustersRequest(&ecs.ListClustersInput{}).Send()
	if err != nil {
		fmt.Println("Failed to list clusters: " + err.Error())
		os.Exit(1)
	}

	if filter == "" {
		clusterNames = listClusterOutput.ClusterArns
	} else {
		for _, cluster := range listClusterOutput.ClusterArns {
			if strings.Contains(strings.ToLower(cluster), strings.ToLower(filter)) {
				clusterNames = append(clusterNames, cluster)
			}
		}
	}
	sort.Strings(clusterNames)
	return clusterNames
}

func describeClusters(client *ecs.ECS, clusters []string) []ecs.Cluster {
	descClusterReq := client.DescribeClustersRequest(&ecs.DescribeClustersInput{Clusters: clusters})
	descClusterOutput, err := descClusterReq.Send()
	if err != nil {
		fmt.Println("Failed to describe clusters: " + err.Error())
		os.Exit(1)
	}
	sort.Slice(descClusterOutput.Clusters, func(i, j int) bool {
		return *descClusterOutput.Clusters[i].ClusterName < *descClusterOutput.Clusters[j].ClusterName
	})
	return descClusterOutput.Clusters
}

// List and describe services running in the ECS cluster
func listServices(client *ecs.ECS, cluster string) []ecs.Service {
	ecsServices := []ecs.Service{}
	serviceNames := []string{}
	req := client.ListServicesRequest(&ecs.ListServicesInput{Cluster: &cluster})
	pager := req.Paginate()
	for pager.Next() {
		page := pager.CurrentPage()
		serviceNames = append(serviceNames, page.ServiceArns...)
	}
	sort.Strings(serviceNames)
	for _, services := range chunk(serviceNames, 10) {
		if len(services) > 0 {
			descServices := describeServices(client, cluster, services)
			ecsServices = append(ecsServices, descServices...)
		}
	}
	return ecsServices
}

// Describe a list of services running in the ECS cluster
func describeServices(client *ecs.ECS, cluster string, services []string) []ecs.Service {
	params := ecs.DescribeServicesInput{Cluster: &cluster, Services: services}
	resp, err := client.DescribeServicesRequest(&params).Send()
	if err != nil {
		fmt.Println("Failed to describe services: " + err.Error())
		os.Exit(1)
	}
	return resp.Services
}

func findService(client *ecs.ECS, cluster, service string) (ecs.Service, error) {
	var ecsService ecs.Service
	runningServices := describeServices(client, cluster, []string{service})
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

// Split a list of strings into a list of smaller lists containing up to `count` items
func chunk(list []string, count int) [][]string {
	newList := make([][]string, len(list)/count+1)
	for index := 0; index < len(list); index += count {
		upperBound := index + count
		if index+count > len(list) {
			upperBound = len(list)
		}
		newIdx := index / count
		newList[newIdx] = list[index:upperBound]
	}
	return newList
}

func loadAWSConfig(region string) aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Println("Failed to load AWS SDK configuration, " + err.Error())
		os.Exit(1)
	}
	if region != "" {
		cfg.Region = region
	}
	return cfg
}

func serviceUp(service *ecs.Service) bool {
	return *service.DesiredCount == *service.RunningCount &&
		len(service.Events) > 0 &&
		strings.Contains(*service.Events[0].Message, "has reached a steady state")
}

func serviceStatus(service *ecs.Service) string {
	status := colorstring.Color("[green][OK]")
	serviceUp := serviceUp(service)
	if !serviceUp {
		status = colorstring.Color("[red][KO]")
	}
	if serviceUp && *service.RunningCount == 0 {
		status = colorstring.Color("[yellow][WARN]")
	}
	return status
}

func serviceOk(service *ecs.Service) bool {
	status := serviceStatus(service)
	return strings.Contains(status, "OK")
}

func serviceTaskDefinition(client *ecs.ECS, service *ecs.Service) ecs.TaskDefinition {
	resp, err := client.DescribeTaskDefinitionRequest(
		&ecs.DescribeTaskDefinitionInput{TaskDefinition: service.TaskDefinition}).Send()
	if err != nil {
		fmt.Println("Failed to describe task definition: " + err.Error())
		os.Exit(1)
	}
	return *resp.TaskDefinition
}

func printServiceDetails(client *ecs.ECS, service *ecs.Service, longOutput bool) {
	colorstring.Printf(
		"%-15s [yellow]%-60s[reset] %-8s running %d/%d  (%s)\n",
		serviceStatus(service), *service.ServiceName,
		*service.Status, *service.RunningCount, *service.DesiredCount,
		shortTaskDefinitionName(*service.TaskDefinition),
	)
	if longOutput == true {
		taskDefinition := serviceTaskDefinition(client, service)
		fmt.Println(linkToConsole(service, clusterNameFromArn(*service.ClusterArn)))
		for _, container := range taskDefinition.ContainerDefinitions {
			portsString := []string{}
			for _, ports := range container.PortMappings {
				portsString = append(portsString, "->"+strconv.FormatInt(*ports.ContainerPort, 10))
			}
			fmt.Printf("- Container: %s\n", *container.Name)
			fmt.Printf("  Image: %s\n", *container.Image)
			fmt.Printf("  Memory: %d / CPU: %d\n", *container.Memory, *container.Cpu)
			fmt.Printf("  Ports: %s\n", strings.Join(portsString, " "))
			fmt.Println("  Environment:")
			for _, env := range container.Environment {
				fmt.Printf("   - %s: %s\n", *env.Name, *env.Value)
			}
		}
		fmt.Println()
	}
}

func linkToConsole(service *ecs.Service, cluster string) string {
	return fmt.Sprintf(
		"https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services/%s/events",
		awsRegion, awsRegion, cluster, *service.ServiceName,
	)
}
