package cmd

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

var (
	eventsClusterFilter string
	eventsServiceFilter string
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List events for services running in your ECS clusters",
	Run:   runCommandEvents,
}

func init() {
	rootCmd.AddCommand(eventsCmd)

	eventsCmd.Flags().StringVarP(&eventsClusterFilter, "cluster", "c", "", "Filter by the name of the ECS cluster")
	eventsCmd.Flags().StringVarP(&eventsServiceFilter, "service", "s", "", "Filter by the name of the ECS service")
}

func runCommandEvents(cmd *cobra.Command, args []string) {
	var (
		services []ecs.Service
		events   []ecs.ServiceEvent
	)

	client := ecs.New(loadAWSConfig(awsRegion))
	clusterNames := listClusters(client, eventsClusterFilter)

	for _, cluster := range describeClusters(client, clusterNames) {
		services = append(services, listServices(client, *cluster.ClusterName, eventsServiceFilter)...)
	}
	for _, svc := range services {
		events = append(events, svc.Events...)
	}
	sort.Sort(byCreatedAt(events))
	for _, event := range events {
		fmt.Printf("%s: %s\n", *event.CreatedAt, *event.Message)
	}
}

type byCreatedAt []ecs.ServiceEvent

func (c byCreatedAt) Len() int           { return len(c) }
func (c byCreatedAt) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byCreatedAt) Less(i, j int) bool { return c[i].CreatedAt.Before(*c[j].CreatedAt) }
