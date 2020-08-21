package aws

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
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

// DescribeClusters describes ECS clusters to fetch detailed information
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
