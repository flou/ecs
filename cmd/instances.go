package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/flou/ecs/pkg/aws"
	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

type instanceOpts struct {
	region        string
	clusterFilter string
	longOutput    bool
}

func buildInstancesCmd() *cobra.Command {
	var opts = instanceOpts{}
	var cmd = &cobra.Command{
		Use:   "instances",
		Short: "List container instances in your ECS clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommandInstances(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.region, "region", "r", "", "AWS region name")
	cmd.Flags().StringVarP(&opts.clusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	cmd.Flags().BoolVarP(&opts.longOutput, "long", "l", false, "Enable detailed output of containers instances")

	return cmd
}

type eInstance struct {
	IPAddress string
	ImageID   string
	Name      string
}

func runCommandInstances(options instanceOpts) error {
	cfg := aws.LoadAWSConfig(options.region)
	client := ecs.New(cfg)
	ec2Client := ec2.New(cfg)

	clusterNames := aws.ListClusters(client, options.clusterFilter)
	for _, cluster := range aws.DescribeClusters(client, clusterNames) {
		listContainerResp, err := client.ListContainerInstancesRequest(
			&ecs.ListContainerInstancesInput{Cluster: cluster.ClusterName}).Send(context.Background())
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
			}).Send(context.Background())
		if err != nil {
			fmt.Println("Failed to describe container instances: " + err.Error())
			os.Exit(1)
		}

		containerInstanceIds := make([]string, 0)
		for _, cinst := range describeContainerInstancesResp.ContainerInstances {
			containerInstanceIds = append(containerInstanceIds, *cinst.Ec2InstanceId)
		}
		describeInstanceResp, err := ec2Client.DescribeInstancesRequest(
			&ec2.DescribeInstancesInput{InstanceIds: containerInstanceIds}).Send(context.Background())

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
					Name:      *aws.FindTag(inst.Tags, "Name").Value,
				}
			}
		}
		fmt.Printf(
			"%-20s  %-8s %5s  %10s  %10s  %15s %10s  %6v  %-21s  %-6s  %10s\n",
			"INSTANCE ID", "STATUS", "TASKS", "CPU:used/free", "MEM:used/free",
			"PRIVATE IP", "INST.TYPE", "AGENT", "AMI", "DOCKER", "AGE",
		)
		for _, cinst := range describeContainerInstancesResp.ContainerInstances {
			agentVersion := colorstring.Color("[green]" + *cinst.VersionInfo.AgentVersion)
			if *cinst.AgentConnected == false {
				agentVersion = colorstring.Color("[red]" + *cinst.VersionInfo.AgentVersion)
			}
			registeredCPU := aws.FindResource(cinst.RegisteredResources, "CPU").IntegerValue
			remainingCPU := aws.FindResource(cinst.RemainingResources, "CPU").IntegerValue
			registeredMemory := aws.FindResource(cinst.RegisteredResources, "MEMORY").IntegerValue
			remainingMemory := aws.FindResource(cinst.RemainingResources, "MEMORY").IntegerValue
			instanceType := aws.FindAttribute(cinst.Attributes, "ecs.instance-type").Value
			dockerVersion := strings.TrimPrefix(*cinst.VersionInfo.DockerVersion, "DockerVersion: ")
			ageInDays := fmt.Sprintf("%4.1f days", time.Since(*cinst.RegisteredAt).Hours()/24)
			instance := instances[*cinst.Ec2InstanceId]
			fmt.Printf(
				"%-20s  %-8s %5d  %13s  %13s  %15s %10s  %-6v  %12s  %7s  %s\n",
				*cinst.Ec2InstanceId, *cinst.Status, *cinst.RunningTasksCount,
				fmt.Sprintf("%d/%d", *registeredCPU-*remainingCPU, *remainingCPU),
				fmt.Sprintf("%d/%d", *registeredMemory-*remainingMemory, *remainingMemory),
				instance.IPAddress, *instanceType, agentVersion, instance.ImageID,
				dockerVersion, ageInDays,
			)
			if options.longOutput == true {
				aws.DetailedInstanceOutput(&cinst)
			}
		}
		fmt.Println()
	}
	return nil
}
