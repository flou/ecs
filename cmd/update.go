package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

var (
	updateCluster      string
	updateService      string
	updateDesiredCount int64
	updateForce        bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the service to a specific DesiredCount",
	Run:   runCommandUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateCluster, "cluster", "c", "", "Name of the ECS cluster")
	updateCmd.MarkFlagRequired("cluster")
	updateCmd.Flags().StringVarP(&updateService, "service", "s", "", "Name of the ECS service")
	updateCmd.MarkFlagRequired("service")

	updateCmd.Flags().Int64Var(&updateDesiredCount, "count", -1, "New DesiredCount")
	updateCmd.Flags().BoolVarP(&updateForce, "force", "f", false, "Force a new deployment of the service")
}

func runCommandUpdate(cmd *cobra.Command, args []string) {
	cfg := loadAWSConfig(awsRegion)
	client := ecs.New(cfg)

	ecsService, err := findService(client, updateCluster, updateService)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	params := ecs.UpdateServiceInput{
		Cluster:            &updateCluster,
		Service:            &updateService,
		ForceNewDeployment: &updateForce,
	}

	if updateDesiredCount >= 0 {
		if *ecsService.DesiredCount == updateDesiredCount {
			colorstring.Printf("Service [yellow]%s[reset] already has a DesiredCount of %d\n",
				updateService, updateDesiredCount,
			)
			return
		}
		colorstring.Printf(
			"Updating [yellow]%s[reset] / DesiredCount[%d -> %d] RunningCount={%d}\n",
			updateService, *ecsService.DesiredCount, updateDesiredCount, *ecsService.RunningCount,
		)
		params.DesiredCount = &updateDesiredCount
	}

	_, err = client.UpdateServiceRequest(&params).Send(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	colorstring.Printf("Service [yellow]%s[reset] successfully updated: DesiredCount=%d\n", updateService, updateDesiredCount)
}
