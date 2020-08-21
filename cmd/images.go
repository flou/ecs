package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/flou/ecs/pkg/aws"
	"github.com/spf13/cobra"
)

type imagesOpts struct {
	region        string
	clusterFilter string
	serviceFilter string
	serviceType   string
}

func buildImagesCmd() *cobra.Command {
	var opts = imagesOpts{}
	var cmd = &cobra.Command{
		Use:   "images",
		Short: "List the Docker images of a service running in ECS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommandImage(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.region, "region", "r", "", "AWS region name")
	cmd.Flags().StringVarP(&opts.clusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	cmd.Flags().StringVarP(&opts.serviceFilter, "service", "s", "", "Filter by the name of the ECS service")
	cmd.Flags().StringVarP(&opts.serviceType, "type", "t", "", "Filter by service launch type")

	return cmd
}

func runCommandImage(options imagesOpts) error {
	client := ecs.New(aws.LoadAWSConfig(options.region))
	clusterNames := aws.ListClusters(client, options.clusterFilter)

	for _, cluster := range aws.DescribeClusters(client, clusterNames) {
		services := aws.ListServices(client, *cluster.ClusterName, options.serviceFilter, options.serviceType)
		fmt.Printf("--- CLUSTER: %s (%d services)\n", *cluster.ClusterName, len(services))
		for _, svc := range services {
			taskDefinition := aws.ServiceTaskDefinition(client, *svc.TaskDefinition)
			for _, container := range taskDefinition.ContainerDefinitions {
				fmt.Printf("%s: %s\n", *svc.ServiceName, *container.Image)
			}
		}
	}
	return nil
}
