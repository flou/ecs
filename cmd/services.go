package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

var (
	servicesClusterName   string
	servicesClusterFilter string
	servicesServiceFilter string
	printAll              bool
	longOutput            bool
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List services in your ECS clusters",
	Run:   runCommandServices,
}

func init() {
	rootCmd.AddCommand(servicesCmd)

	servicesCmd.Flags().StringVarP(&servicesClusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	servicesCmd.Flags().StringVarP(&servicesServiceFilter, "service", "s", "", "Filter by the name of the ECS service")
	servicesCmd.Flags().BoolVarP(&printAll, "all", "a", false, "Print all services, ignoring their status")
	servicesCmd.Flags().BoolVarP(&longOutput, "long", "l", false, "Enable detailed output of containers parameters")
}

func runCommandServices(cmd *cobra.Command, args []string) {
	cfg := loadAWSConfig(awsRegion)
	client := ecs.New(cfg)

	clusterNames := []string{servicesClusterName}
	if servicesClusterName == "" {
		clusterNames = listClusters(client, servicesClusterFilter)
	}
	for _, cluster := range describeClusters(client, clusterNames) {
		services := listServices(client, *cluster.ClusterName, servicesServiceFilter)
		headerLine := fmt.Sprintf(
			"--- CLUSTER: %s (%d services)\n", *cluster.ClusterName, len(services),
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
		if len(services) != 0 {
			fmt.Printf(headerLine)
			for _, svc := range services {
				printServiceDetails(client, &svc, longOutput)
			}
			fmt.Println()
		}
	}
}
