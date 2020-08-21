package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/flou/ecs/pkg/aws"
	"github.com/spf13/cobra"
)

type servicesOpts struct {
	region        string
	clusterFilter string
	serviceFilter string
	serviceType   string
	printAll      bool
	longOutput    bool
}

func buildServicesCmd() *cobra.Command {
	var opts = servicesOpts{}
	var cmd = &cobra.Command{
		Use:   "services",
		Short: "List services in your ECS clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommandServices(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.region, "region", "r", "", "AWS region name")
	cmd.Flags().StringVarP(&opts.clusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	cmd.Flags().StringVarP(&opts.serviceFilter, "service", "s", "", "Filter by the name of the ECS service")
	cmd.Flags().StringVarP(&opts.serviceType, "type", "t", "", "Filter by service launch type")
	cmd.Flags().BoolVarP(&opts.printAll, "all", "a", false, "Print all services, ignoring their status")
	cmd.Flags().BoolVarP(&opts.longOutput, "long", "l", false, "Enable detailed output of containers parameters")

	return cmd
}

func runCommandServices(options servicesOpts) error {
	cfg := aws.LoadAWSConfig(options.region)
	client := ecs.New(cfg)

	clusterNames := aws.ListClusters(client, options.clusterFilter)
	for _, cluster := range aws.DescribeClusters(client, clusterNames) {
		services := aws.ListServices(client, *cluster.ClusterName, options.serviceFilter, options.serviceType)
		headerLine := fmt.Sprintf(
			"--- CLUSTER: %s (%d services)\n", *cluster.ClusterName, len(services),
		)

		if options.printAll == false {
			var displayedServices []ecs.Service
			for _, svc := range services {
				if !aws.ServiceOk(&svc) {
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
				aws.PrintServiceDetails(client, &svc, options.longOutput)
			}
		}
		fmt.Println()
	}
	return nil
}
