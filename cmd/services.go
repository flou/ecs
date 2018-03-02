package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

var servicesCluster string
var servicesFilter string
var printAll bool
var longOutput bool
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List unhealthy services in your ECS clusters",
	Run:   runCommandServices,
}

func init() {
	rootCmd.AddCommand(servicesCmd)

	servicesCmd.Flags().StringVarP(&servicesFilter, "filter", "f", "", "Filter by the name of the ECS cluster")
	servicesCmd.Flags().StringVar(&servicesCluster, "cluster", "", "Filter by the name of the ECS cluster")
	servicesCmd.Flags().BoolVarP(&printAll, "all", "a", false, "Filter by the name of the ECS cluster")
	servicesCmd.Flags().BoolVarP(&longOutput, "long", "l", false, "Filter by the name of the ECS cluster")
}

func runCommandServices(cmd *cobra.Command, args []string) {
	cfg := loadAWSConfig(awsRegion)
	client := ecs.New(cfg)

	clusterNames := []string{servicesCluster}
	if servicesCluster == "" {
		clusterNames = listClusters(client, servicesFilter)
	}
	for _, cluster := range describeClusters(client, clusterNames) {
		services := listServices(client, *cluster.ClusterName)
		headerLine := fmt.Sprintf(
			"--- CLUSTER: %s (%d services)\n",
			*cluster.ClusterName, len(services),
		)

		if printAll == false {
			var displayedServices []ecs.Service
			for _, svc := range services {
				if !serviceOk(&svc) {
					displayedServices = append(displayedServices, svc)
				}
			}
			if len(displayedServices) > 0 {
				headerLine = fmt.Sprintf(
					"--- CLUSTER: %s (listing %d/%d services)\n",
					*cluster.ClusterName, len(displayedServices), len(services),
				)
				services = displayedServices
			}
		}
		fmt.Printf(headerLine)
		for _, svc := range services {
			printServiceDetails(client, &svc, longOutput)
		}
		fmt.Println()
	}
}
