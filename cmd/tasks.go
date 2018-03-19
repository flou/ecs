package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

var (
	tasksClusterFilter string
	tasksServiceFilter string
	tasksLongOutput    bool
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks running in your ECS clusters",
	Run:   runCommandTasks,
}

func init() {
	rootCmd.AddCommand(tasksCmd)

	tasksCmd.Flags().StringVarP(&tasksClusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	tasksCmd.Flags().StringVarP(&tasksServiceFilter, "service", "s", "", "Filter by the name of the ECS service")
	tasksCmd.Flags().BoolVarP(&tasksLongOutput, "long", "l", false, "Enable detailed output of containers parameters")
}

func runCommandTasks(cmd *cobra.Command, args []string) {
	cfg := loadAWSConfig(awsRegion)
	client := ecs.New(cfg)

	clusterNames := listClusters(client, tasksClusterFilter)
	for _, cluster := range describeClusters(client, clusterNames) {
		tasks := listTasks(client, *cluster.ClusterName, tasksServiceFilter)
		headerLine := fmt.Sprintf(
			"--- CLUSTER: %s (%d tasks)\n", *cluster.ClusterName, len(tasks),
		)

		if len(tasks) != 0 {
			fmt.Printf(headerLine)
			for _, task := range tasks {
				printTaskDetails(client, &task)
			}
			fmt.Println()
		}
	}
}
