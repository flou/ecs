package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func shortTaskDefinitionName(taskDefinition string) string {
	splitTaskDefinition := strings.Split(taskDefinition, "/")
	return strings.Split(taskDefinition, "/")[len(splitTaskDefinition)-1]
}

func clusterNameFromArn(clusterArn string) string {
	splitClusterArn := strings.Split(clusterArn, "/")
	return strings.Split(clusterArn, "/")[len(splitClusterArn)-1]
}

// Split a list of strings into a list of smaller lists containing up to `count` items
func chunk(list []string, count int) [][]string {
	newList := make([][]string, len(list)/count+1)
	for index := 0; index < len(list); index += count {
		upperBound := index + count
		if index+count > len(list) {
			upperBound = len(list)
		}
		newIdx := index / count
		newList[newIdx] = list[index:upperBound]
	}
	return newList
}

func linkToIAM(roleArn string) string {
	return fmt.Sprintf("https://console.aws.amazon.com/iam/home#/roles/%s", roleArn)
}

func linkToConsole(service *ecs.Service, cluster string) string {
	awsRegion := strings.Split(*service.ServiceArn, ":")[3]
	return fmt.Sprintf(
		"https://%s.console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services/%s/events",
		awsRegion, awsRegion, cluster, *service.ServiceName,
	)
}
