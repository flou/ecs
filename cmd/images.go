package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

var (
	imagesClusterFilter string
	imagesServiceFilter string
)
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List the Docker images of a service running in ECS",
	Run:   runCommandImage,
}

func init() {
	rootCmd.AddCommand(imagesCmd)

	imagesCmd.Flags().StringVarP(&imagesClusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	imagesCmd.Flags().StringVarP(&imagesServiceFilter, "service", "s", "", "Filter by the name of the ECS service")
}

func runCommandImage(cmd *cobra.Command, args []string) {
	client := ecs.New(loadAWSConfig(awsRegion))
	clusterNames := listClusters(client, imagesClusterFilter)

	for _, cluster := range describeClusters(client, clusterNames) {
		services := listServices(client, *cluster.ClusterName, imagesServiceFilter)
		fmt.Printf("--- CLUSTER: %s (%d services)\n", *cluster.ClusterName, len(services))
		for _, svc := range services {
			taskDefinition := serviceTaskDefinition(client, *svc.TaskDefinition)
			for _, container := range taskDefinition.ContainerDefinitions {
				fmt.Printf("%s: %s\n", *svc.ServiceName, *container.Image)
			}
		}
		fmt.Println()
	}
}
