package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/flou/ecs/pkg/aws"
	"github.com/spf13/cobra"
)

type tasksOpts struct {
	region        string
	clusterFilter string
	serviceFilter string
	longOutput    bool
}

func buildTasksCmd() *cobra.Command {
	var opts = tasksOpts{}
	var cmd = &cobra.Command{
		Use:   "tasks",
		Short: "List tasks running in your ECS clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommandTasks(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.region, "region", "r", "", "AWS region name")
	cmd.Flags().StringVarP(&opts.clusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	cmd.Flags().StringVarP(&opts.serviceFilter, "service", "s", "", "Filter by the name of the ECS service")
	cmd.Flags().BoolVarP(&opts.longOutput, "long", "l", false, "Enable detailed output of containers parameters")

	return cmd
}

func runCommandTasks(options tasksOpts) error {
	cfg := aws.LoadAWSConfig(options.region)
	client := ecs.New(cfg)

	clusterNames := aws.ListClusters(client, options.clusterFilter)
	for _, cluster := range aws.DescribeClusters(client, clusterNames) {
		tasks := aws.ListTasks(client, *cluster.ClusterName, options.serviceFilter)
		headerLine := fmt.Sprintf(
			"--- CLUSTER: %s (%d tasks)\n", *cluster.ClusterName, len(tasks),
		)

		if len(tasks) != 0 {
			fmt.Printf(headerLine)
			for _, task := range tasks {
				aws.PrintTaskDetails(client, &task, options.longOutput)
			}
		}
	}
	return nil
}
