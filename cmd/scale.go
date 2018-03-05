package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

var (
	scaleCluster      string
	scaleService      string
	scaleDesiredCount int64
)

var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale the service to a specific DesiredCount",
	Run:   runCommandScale,
}

func init() {
	rootCmd.AddCommand(scaleCmd)

	scaleCmd.Flags().StringVar(&scaleCluster, "cluster", "", "Name of the ECS cluster")
	scaleCmd.Flags().StringVar(&scaleService, "service", "", "Name of the ECS service")
	scaleCmd.Flags().Int64Var(&scaleDesiredCount, "count", 0, "New DesiredCount")
	scaleCmd.MarkFlagRequired("cluster")
	scaleCmd.MarkFlagRequired("service")
	scaleCmd.MarkFlagRequired("count")
}

func runCommandScale(cmd *cobra.Command, args []string) {
	cfg := loadAWSConfig(awsRegion)
	client := ecs.New(cfg)

	ecsService, err := findService(client, scaleCluster, scaleService)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if *ecsService.DesiredCount == scaleDesiredCount {
		colorstring.Printf("Service [yellow]%s[reset] already has a DesiredCount of %d\n",
			scaleService, scaleDesiredCount,
		)
		return
	}
	colorstring.Printf(
		"Updating [yellow]%s[reset] / DesiredCount[%d -> %d] RunningCount={%d}\n",
		scaleService, *ecsService.DesiredCount, scaleDesiredCount, *ecsService.RunningCount,
	)
	_, err = client.UpdateServiceRequest(&ecs.UpdateServiceInput{
		Cluster:      &scaleCluster,
		Service:      &scaleService,
		DesiredCount: &scaleDesiredCount,
	}).Send()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Service %s successfully updated with DesiredCount=%d", scaleService, scaleDesiredCount)
}
