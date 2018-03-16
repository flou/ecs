package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

const timeFormat = "2018-02-23 10:46:01 +0000 UTC"

var (
	instancesFilter     string
	instancesLongOutput bool
)
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List container instances in your ECS clusters",
	Run:   runCommandInstances,
}

type eInstance struct {
	IPAddress string
	ImageID   string
	Name      string
}

func init() {
	rootCmd.AddCommand(instancesCmd)

	instancesCmd.Flags().StringVarP(&instancesFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	instancesCmd.Flags().BoolVarP(&instancesLongOutput, "long", "l", false, "Enable detailed output of containers instances")
}

func runCommandInstances(cmd *cobra.Command, args []string) {
	cfg := loadAWSConfig(awsRegion)
	client := ecs.New(cfg)
	ec2Client := ec2.New(cfg)

	clusterNames := listClusters(client, instancesFilter)
	for _, cluster := range describeClusters(client, clusterNames) {
		var containerInstanceIds []string
		listContainerResp, err := client.ListContainerInstancesRequest(
			&ecs.ListContainerInstancesInput{Cluster: cluster.ClusterName}).Send()
		if err != nil {
			fmt.Println("Failed to list container instances: " + err.Error())
			os.Exit(1)
		}
		fmt.Printf("--- CLUSTER: %s (%d registered instances)\n",
			*cluster.ClusterName,
			len(listContainerResp.ContainerInstanceArns),
		)
		if len(listContainerResp.ContainerInstanceArns) == 0 {
			fmt.Println()
			continue
		}
		describeContainerInstancesResp, err := client.DescribeContainerInstancesRequest(
			&ecs.DescribeContainerInstancesInput{
				Cluster:            cluster.ClusterName,
				ContainerInstances: listContainerResp.ContainerInstanceArns,
			}).Send()
		if err != nil {
			fmt.Println("Failed to describe container instances: " + err.Error())
			os.Exit(1)
		}

		for _, cinst := range describeContainerInstancesResp.ContainerInstances {
			containerInstanceIds = append(containerInstanceIds, *cinst.Ec2InstanceId)
		}
		describeInstanceResp, err := ec2Client.DescribeInstancesRequest(
			&ec2.DescribeInstancesInput{InstanceIds: containerInstanceIds}).Send()

		if err != nil {
			fmt.Println("Failed to describe EC2 instances: " + err.Error())
			os.Exit(1)
		}
		instances := make(map[string]eInstance)
		for _, res := range describeInstanceResp.Reservations {
			for _, inst := range res.Instances {
				instances[*inst.InstanceId] = eInstance{
					IPAddress: *inst.PrivateIpAddress,
					ImageID:   *inst.ImageId,
					Name:      *findTag(inst.Tags, "Name").Value,
				}
			}
		}
		fmt.Printf(
			"%-20s  %-8s %5s  %8s %8s  %8s %8s  %15s %12s  %6v  %-12s  %-8s  %11s  %s\n",
			"INSTANCE ID", "STATUS", "TASKS", "CPU/used", "CPU/free",
			"MEM/used", "MEM/free", "PRIVATE IP", "INST.TYPE", "AGENT",
			"IMAGE", "DOCKER", "AGE", "NAME",
		)
		for _, cinst := range describeContainerInstancesResp.ContainerInstances {
			agentVersion := colorstring.Color("[green]" + *cinst.VersionInfo.AgentVersion)
			if *cinst.AgentConnected == false {
				agentVersion = colorstring.Color("[red]" + *cinst.VersionInfo.AgentVersion)
			}
			registeredCPU := *findResource(cinst.RegisteredResources, "CPU").IntegerValue
			remainingCPU := *findResource(cinst.RemainingResources, "CPU").IntegerValue
			registeredMemory := *findResource(cinst.RegisteredResources, "MEMORY").IntegerValue
			remainingMemory := *findResource(cinst.RemainingResources, "MEMORY").IntegerValue
			instanceType := *findAttribute(cinst.Attributes, "ecs.instance-type").Value
			dockerVersion := strings.TrimPrefix(*cinst.VersionInfo.DockerVersion, "DockerVersion: ")
			ageInDays := fmt.Sprintf("%4.1f days", time.Since(*cinst.RegisteredAt).Hours()/24)
			instance := instances[*cinst.Ec2InstanceId]
			fmt.Printf(
				"%-20s  %-8s %5d  %8d %8d  %8d %8d  %15s %12s  %-6v  %12s  %10s  %s  %s\n",
				*cinst.Ec2InstanceId, *cinst.Status, *cinst.RunningTasksCount,
				registeredCPU-remainingCPU, remainingCPU, registeredMemory-remainingMemory, remainingMemory,
				instance.IPAddress, instanceType, agentVersion, instance.ImageID,
				dockerVersion, ageInDays, instance.Name,
			)
			instanceAttributes := []string{}
			capabilities := []string{}
			var line string
			if instancesLongOutput == true {
				for _, attr := range cinst.Attributes {
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
		}
		fmt.Println()
	}
}
