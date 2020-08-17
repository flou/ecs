package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/flou/ecs/pkg/aws"
	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

type updateCmd struct {
	cmd  *cobra.Command
	opts updateOpts
}

type updateOpts struct {
	region       string
	cluster      string
	service      string
	desiredCount int64
	force        bool
}

func buildUpdateCmd() *updateCmd {
	var root = &updateCmd{}
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "Update the service to a specific DesiredCount",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommandUpdate(root.opts)
		},
	}

	cmd.Flags().StringVarP(&root.opts.region, "region", "r", "", "AWS region name")
	cmd.Flags().StringVarP(&root.opts.cluster, "cluster", "c", "", "Name of the ECS cluster")
	cmd.MarkFlagRequired("cluster")
	cmd.Flags().StringVarP(&root.opts.service, "service", "s", "", "Name of the ECS service")
	cmd.MarkFlagRequired("service")

	cmd.Flags().Int64Var(&root.opts.desiredCount, "count", -1, "New DesiredCount")
	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Force a new deployment of the service")

	root.cmd = cmd
	return root
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
			colorstring.Printf("Service [yellow]%s[reset] already has a DesiredCount of %d\n",
				options.service, options.desiredCount,
			)
			return nil
		}
		colorstring.Printf(
			"Updating [yellow]%s[reset] / DesiredCount[%d -> %d] RunningCount={%d}\n",
			options.service, *ecsService.DesiredCount, options.desiredCount, *ecsService.RunningCount,
		)
		params.DesiredCount = &options.desiredCount
	}

	_, err = client.UpdateServiceRequest(&params).Send(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	colorstring.Printf("Service [yellow]%s[reset] successfully updated: DesiredCount=%d\n", options.service, options.desiredCount)
	return nil
}
