package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

var instancesFilter string
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List container instances in your ECS clusters",
	Run:   runCommandInstances,
}

func init() {
	rootCmd.AddCommand(instancesCmd)

	instancesCmd.Flags().StringVarP(&instancesFilter, "filter", "f", "", "Filter by the name of the ECS cluster")
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
		instancesAttrs := make(map[string]map[string]string)
		for _, res := range describeInstanceResp.Reservations {
			for _, inst := range res.Instances {
				instance := map[string]string{
					"IpAddress": *inst.PrivateIpAddress,
					"ImageId":   *inst.ImageId,
					"NameTag":   *findTag(inst.Tags, "Name").Value,
				}
				instancesAttrs[*inst.InstanceId] = instance
			}
		}
		fmt.Printf(
			"%-20s  %-8s %5s  %8s %8s  %8s %8s  %15s %12s  %5v  %-12s  %s\n",
			"INSTANCE ID", "STATUS", "TASKS", "CPU/used", "CPU/free",
			"MEM/used", "MEM/free", "PRIVATE IP", "INST.TYPE", "AGENT",
			"IMAGE", "NAME",
		)
		for _, cinst := range describeContainerInstancesResp.ContainerInstances {
			registeredCPU := *findResource(cinst.RegisteredResources, "CPU").IntegerValue
			remainingCPU := *findResource(cinst.RemainingResources, "CPU").IntegerValue
			registeredMemory := *findResource(cinst.RegisteredResources, "MEMORY").IntegerValue
			remainingMemory := *findResource(cinst.RemainingResources, "MEMORY").IntegerValue
			instanceType := *findAttribute(cinst.Attributes, "ecs.instance-type").Value
			instanceAttrs := instancesAttrs[*cinst.Ec2InstanceId]
			fmt.Printf(
				"%-20s  %-8s %5d  %8d %8d  %8d %8d  %15s %12s  %-5v  %12s  %s\n",
				*cinst.Ec2InstanceId, *cinst.Status, *cinst.RunningTasksCount,
				registeredCPU-remainingCPU, remainingCPU, registeredMemory-remainingMemory, remainingMemory,
				instanceAttrs["IpAddress"], instanceType, *cinst.AgentConnected, instanceAttrs["ImageId"],
				instanceAttrs["NameTag"],
			)
		}
		fmt.Println()
	}
}
