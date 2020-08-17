package aws

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mitchellh/colorstring"
)

type byTaskName []ecs.Task

func (c byTaskName) Len() int           { return len(c) }
func (c byTaskName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byTaskName) Less(i, j int) bool { return *c[i].TaskDefinitionArn < *c[j].TaskDefinitionArn }

func ListTasks(client *ecs.Client, clusterName, taskFilter string) []ecs.Task {
	req := client.ListTasksRequest(&ecs.ListTasksInput{Cluster: &clusterName})
	p := ecs.NewListTasksPaginator(req)

	taskNames := make([]string, 0)
	for p.Next(context.Background()) {
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

func describeTasks(client *ecs.Client, clusterName string, tasks []string) []ecs.Task {
	params := ecs.DescribeTasksInput{Cluster: &clusterName, Tasks: tasks}
	resp, err := client.DescribeTasksRequest(&params).Send(context.Background())
	if err != nil {
		fmt.Println("Failed to describe services: " + err.Error())
		os.Exit(1)
	}
	return resp.Tasks
}

func PrintTaskDetails(client *ecs.Client, task *ecs.Task, longOutput bool) {
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
	if longOutput == true {
		taskDefinition := ServiceTaskDefinition(client, *task.TaskDefinitionArn)
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
