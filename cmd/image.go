package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

var cluster string
var service string
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Print the Docker image of a service running in ECS",
	Run:   runCommandImage,
}

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().StringVar(&cluster, "cluster", "", "Name of the ECS cluster")
	imageCmd.Flags().StringVar(&service, "service", "", "Name of the ECS service")
	imageCmd.MarkFlagRequired("cluster")
	imageCmd.MarkFlagRequired("service")
}

func runCommandImage(cmd *cobra.Command, args []string) {
	cfg := loadAWSConfig(awsRegion)
	client := ecs.New(cfg)

	ecsService, err := findService(client, cluster, service)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	taskDefinition := serviceTaskDefinition(client, &ecsService)
	for _, container := range taskDefinition.ContainerDefinitions {
		fmt.Println(*container.Image)
	}
}
