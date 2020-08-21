package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/flou/ecs/pkg/aws"
	"github.com/spf13/cobra"
)

type eventsOpts struct {
	region        string
	clusterFilter string
	serviceFilter string
	serviceType   string
	skipSteady    bool
}

func buildEventsCmd() *cobra.Command {
	var opts = eventsOpts{}
	var cmd = &cobra.Command{
		Use:   "events",
		Short: "List events for services running in your ECS clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommandEvents(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.region, "region", "r", "", "AWS region name")
	cmd.Flags().StringVarP(&opts.clusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	cmd.Flags().StringVarP(&opts.serviceFilter, "service", "s", "", "Filter by the name of the ECS service")
	cmd.Flags().StringVarP(&opts.serviceType, "type", "t", "", "Filter by service launch type")
	cmd.Flags().BoolVar(&opts.skipSteady, "skip-steady", false, "Don't display events that say the service is steady")

	return cmd
}

func runCommandEvents(options eventsOpts) error {
	var services []ecs.Service
	var events []ecs.ServiceEvent

	client := ecs.New(aws.LoadAWSConfig(options.region))
	clusterNames := aws.ListClusters(client, options.clusterFilter)

	for _, cluster := range aws.DescribeClusters(client, clusterNames) {
		services = append(
			services,
			aws.ListServices(client, *cluster.ClusterName, options.serviceFilter, options.serviceType)...,
		)
	}
	for _, svc := range services {
		events = append(events, svc.Events...)
	}
	sort.Sort(aws.EventsByCreatedAt(events))
	for _, event := range events {
		if !strings.Contains(*event.Message, "has reached a steady state.") || options.skipSteady == false {
			fmt.Printf("%s: %s\n", *event.CreatedAt, *event.Message)
		}
	}
	return nil
}
