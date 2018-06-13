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
		for _, clusterArn := range listClusterOutput.ClusterArns {
			cluster := clusterNameFromArn(clusterArn)
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
func listServices(client *ecs.ECS, clusterName, serviceFilter string) []ecs.Service {
	ecsServices := []ecs.Service{}
	serviceNames := []string{}
	listServicesInput := ecs.ListServicesInput{Cluster: &clusterName}
	if strings.ToLower(servicesServiceType) == "fargate" {
		listServicesInput.LaunchType = "FARGATE"
	} else if strings.ToLower(servicesServiceType) == "ec2" {
		listServicesInput.LaunchType = "EC2"
	}
	req := client.ListServicesRequest(&listServicesInput)
	pager := req.Paginate()
	for pager.Next() {
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
			descServices := describeServices(client, clusterName, services)
			ecsServices = append(ecsServices, descServices...)
		}
	}
	return ecsServices
}

// Describe a list of services running in the ECS cluster
func describeServices(client *ecs.ECS, clusterName string, services []string) []ecs.Service {
	params := ecs.DescribeServicesInput{Cluster: &clusterName, Services: services}
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

type byTaskName []ecs.Task

func (c byTaskName) Len() int           { return len(c) }
func (c byTaskName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byTaskName) Less(i, j int) bool { return *c[i].TaskDefinitionArn < *c[j].TaskDefinitionArn }

func listTasks(client *ecs.ECS, clusterName, taskFilter string) []ecs.Task {
	req := client.ListTasksRequest(&ecs.ListTasksInput{Cluster: &clusterName})
	p := req.Paginate()

	taskNames := make([]string, 0)
	for p.Next() {
		page := p.CurrentPage()
		taskNames = append(taskNames, page.TaskArns...)
	}

	ecsTasks := make([]ecs.Task, 0)
	for _, tasks := range chunk(taskNames, 100) {
		if len(tasks) > 0 {
			for _, t := range describeTasks(client, clusterName, tasks) {
				if strings.Contains(*t.TaskDefinitionArn, taskFilter) {
					ecsTasks = append(ecsTasks, t)
				}
			}
		}
	}
	sort.Sort(byTaskName(ecsTasks))
	return ecsTasks
}

func describeTasks(client *ecs.ECS, clusterName string, tasks []string) []ecs.Task {
	params := ecs.DescribeTasksInput{Cluster: &clusterName, Tasks: tasks}
	resp, err := client.DescribeTasksRequest(&params).Send()
	if err != nil {
		fmt.Println("Failed to describe services: " + err.Error())
		os.Exit(1)
	}
	return resp.Tasks
}

func printTaskDetails(client *ecs.ECS, task *ecs.Task) {
	fmt.Printf(
		"%-60s  %-10s", shortTaskDefinitionName(*task.TaskDefinitionArn), *task.LastStatus,
	)
	if task.Cpu != nil {
		fmt.Printf("  Cpu: %4s", *task.Cpu)
	}
	if task.Memory != nil {
		fmt.Printf("  Memory: %4s", *task.Memory)
	}
	fmt.Println()
	if tasksLongOutput == true {
		taskDefinition := serviceTaskDefinition(client, *task.TaskDefinitionArn)
		for _, container := range taskDefinition.ContainerDefinitions {
			colorstring.Printf("- Container: [green]%s\n", *container.Name)
			fmt.Printf("  Image: %s\n", *container.Image)
			fmt.Printf("  Memory: %d / CPU: %d\n", *container.Memory, *container.Cpu)
			if len(container.PortMappings) > 0 {
				fmt.Println("  Ports:")
				for _, port := range container.PortMappings {
					fmt.Printf(
						"   - Host:%d -> Container:%d\n",
						*port.HostPort, *port.ContainerPort,
					)
				}
			}
			if len(container.Environment) > 0 {
				fmt.Println("  Environment:")
				for _, env := range container.Environment {
					fmt.Printf("   - %s: %s\n", *env.Name, *env.Value)
				}
			}
			if len(container.Links) > 0 {
				fmt.Printf("  Links: %s\n", strings.Join(container.Links, ","))
			}
			if container.LogConfiguration != nil {
				fmt.Printf("  Logs: %s", container.LogConfiguration.LogDriver)
				switch container.LogConfiguration.LogDriver {
				case "awslogs":
					fmt.Printf(" (%s)\n", container.LogConfiguration.Options["awslogs-group"])
				case "fluentd":
					fmt.Printf(" (tag: %s)\n", container.LogConfiguration.Options["tag"])
				default:
					fmt.Printf("\n")
				}
			}
		}
		fmt.Println()
	}
}

func detailedInstanceOutput(containerInstance *ecs.ContainerInstance) {
	var line string
	instanceAttributes := make([]string, 0)
	capabilities := make([]string, 0)
	for _, attr := range containerInstance.Attributes {
		if strings.Contains(*attr.Name, "ecs.capability.") {
			capability := strings.SplitAfter(*attr.Name, "ecs.capability.")[1]
			if strings.HasPrefix(capability, "docker-remote-api.") {
				continue
			}
			if attr.Value == nil {
				line = fmt.Sprintf(" - %s", capability)
			} else {
				line = fmt.Sprintf(" - %-22s %s", capability, colorstring.Color("[yellow]"+*attr.Value))
			}
			capabilities = append(capabilities, line)
		} else {
			if attr.Value == nil {
			} else {
				line = fmt.Sprintf(" - %s", *attr.Name)
				line = fmt.Sprintf(" - %-22s %s", *attr.Name, colorstring.Color("[yellow]"+*attr.Value))
			}
			instanceAttributes = append(instanceAttributes, line)
		}
	}
	fmt.Println("Attributes:")
	sort.Strings(instanceAttributes)
	for _, attr := range instanceAttributes {
		fmt.Println(attr)
	}
	fmt.Println("Capabilities:")
	sort.Strings(capabilities)
	for _, attr := range capabilities {
		fmt.Println(attr)
	}
	fmt.Println()
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
	defaultRegion := os.Getenv("AWS_DEFAULT_REGION")
	if region != "" {
		cfg.Region = region
	} else if defaultRegion != "" {
		awsRegion = defaultRegion
		cfg.Region = defaultRegion
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

func serviceTaskDefinition(client *ecs.ECS, taskDefinition string) ecs.TaskDefinition {
	resp, err := client.DescribeTaskDefinitionRequest(
		&ecs.DescribeTaskDefinitionInput{TaskDefinition: &taskDefinition}).Send()
	if err != nil {
		fmt.Println("Failed to describe task definition: " + err.Error())
		os.Exit(1)
	}
	return *resp.TaskDefinition
}

func printServiceDetails(client *ecs.ECS, service *ecs.Service, longOutput bool) {
	colorstring.Printf(
		"%-15s [yellow]%-60s[reset] %-7s %-8s running %d/%d  (%s)\n",
		serviceStatus(service), *service.ServiceName,
		service.LaunchType, *service.Status, *service.RunningCount,
		*service.DesiredCount, shortTaskDefinitionName(*service.TaskDefinition),
	)
	if longOutput == true {
		taskDefinition := serviceTaskDefinition(client, *service.TaskDefinition)
		fmt.Println(linkToConsole(service, clusterNameFromArn(*service.ClusterArn)))
		if taskDefinition.TaskRoleArn != nil {
			fmt.Printf("IAM Role: %s\n", linkToIAM(shortTaskDefinitionName(*taskDefinition.TaskRoleArn)))
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
			colorstring.Printf("- Container: [green]%s\n", *container.Name)
			colorstring.Printf("  Image: [yellow]%s\n", *container.Image)
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

func linkToIAM(roleArn string) string {
	return "https://console.aws.amazon.com/iam/home#/roles/" + roleArn
}

func linkToConsole(service *ecs.Service, cluster string) string {
	return fmt.Sprintf(
		"https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services/%s/events",
		awsRegion, awsRegion, cluster, *service.ServiceName,
	)
}
