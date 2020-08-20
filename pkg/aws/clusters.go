package aws

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mitchellh/colorstring"
)

// ListClusters returns the list of clusters sorted by name
func ListClusters(client *ecs.Client, filter string) []string {
	clusterNames := []string{}
	listClusterOutput, err := client.ListClustersRequest(&ecs.ListClustersInput{}).Send(context.Background())
	if err != nil {
		fmt.Println("Failed to list clusters: " + err.Error())
		os.Exit(1)
	}

	if filter == "" {
		clusterNames = listClusterOutput.ClusterArns
	} else {
		for _, clusterArn := range listClusterOutput.ClusterArns {
			cluster := clusterNameFromArn(clusterArn)
			if strings.Contains(strings.ToLower(cluster), strings.ToLower(filter)) {
				clusterNames = append(clusterNames, cluster)
			}
		}
	}
	sort.Strings(clusterNames)
	return clusterNames
}

func DescribeClusters(client *ecs.Client, clusters []string) []ecs.Cluster {
	descClusterReq := client.DescribeClustersRequest(&ecs.DescribeClustersInput{Clusters: clusters})
	descClusterOutput, err := descClusterReq.Send(context.Background())
	if err != nil {
		fmt.Println("Failed to describe clusters: " + err.Error())
		os.Exit(1)
	}
	sort.Slice(descClusterOutput.Clusters, func(i, j int) bool {
		return *descClusterOutput.Clusters[i].ClusterName < *descClusterOutput.Clusters[j].ClusterName
	})
	return descClusterOutput.Clusters
}

func DetailedInstanceOutput(containerInstance *ecs.ContainerInstance) {
	var line string
	instanceAttributes := make([]string, 0)
	capabilities := make([]string, 0)
	for _, attr := range containerInstance.Attributes {
		if strings.Contains(*attr.Name, "ecs.capability.") {
			capability := strings.SplitAfter(*attr.Name, "ecs.capability.")[1]
			if strings.HasPrefix(capability, "docker-remote-api.") {
				continue
			}
			if attr.Value == nil {
				line = fmt.Sprintf(" - %s", capability)
			} else {
				line = fmt.Sprintf(" - %-22s %s", capability, colorstring.Color("[yellow]"+*attr.Value))
			}
			capabilities = append(capabilities, line)
		} else {
			if attr.Value == nil {
				line = fmt.Sprintf(" - %s", *attr.Name)
			} else {
				line = fmt.Sprintf(" - %-22s %s", *attr.Name, colorstring.Color("[yellow]"+*attr.Value))
			}
			instanceAttributes = append(instanceAttributes, line)
		}
	}
	fmt.Println("Attributes:")
	sort.Strings(instanceAttributes)
	for _, attr := range instanceAttributes {
		fmt.Println(attr)
	}
	fmt.Println("Capabilities:")
	sort.Strings(capabilities)
	for _, attr := range capabilities {
		fmt.Println(attr)
	}
	fmt.Println()
}
