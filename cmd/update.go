package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/fatih/color"
	"github.com/flou/ecs/pkg/aws"
	"github.com/spf13/cobra"
)

type updateOpts struct {
	region       string
	cluster      string
	service      string
	desiredCount int64
	force        bool
}

func buildUpdateCmd() *cobra.Command {
	var opts = updateOpts{}
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "Update the service to a specific DesiredCount",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommandUpdate(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.region, "region", "r", "", "AWS region name")
	cmd.Flags().StringVarP(&opts.cluster, "cluster", "c", "", "Name of the ECS cluster")
	cmd.MarkFlagRequired("cluster")
	cmd.Flags().StringVarP(&opts.service, "service", "s", "", "Name of the ECS service")
	cmd.MarkFlagRequired("service")

	cmd.Flags().Int64Var(&opts.desiredCount, "count", -1, "New DesiredCount")
	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "Force a new deployment of the service")

	return cmd
}

func runCommandUpdate(options updateOpts) error {
	cfg := aws.LoadAWSConfig(options.region)
	client := ecs.New(cfg)

	ecsService, err := aws.FindService(client, options.cluster, options.service)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	params := ecs.UpdateServiceInput{
		Cluster:            &options.cluster,
		Service:            &options.service,
		ForceNewDeployment: &options.force,
	}

	if options.desiredCount >= 0 {
		if *ecsService.DesiredCount == options.desiredCount {
			fmt.Printf("Service %s already has a DesiredCount of %d\n",
				color.YellowString(options.service), options.desiredCount,
			)
			return nil
		}
		fmt.Printf(
			"Updating %s / DesiredCount[%d -> %d] RunningCount={%d}\n",
			color.YellowString(options.service), *ecsService.DesiredCount, options.desiredCount, *ecsService.RunningCount,
		)
		params.DesiredCount = &options.desiredCount
	}

	_, err = client.UpdateServiceRequest(&params).Send(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Service %s successfully updated: DesiredCount=%d\n", color.YellowString(options.service), options.desiredCount)
	return nil
}
