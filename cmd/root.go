package cmd

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

type rootOpts struct {
	debug bool
}

// Execute is the root command for the ecs CLI
func Execute(version string) error {
	log.SetHandler(cli.Default)
	var opts = rootOpts{}
	var cmd = &cobra.Command{
		Use:           "ecs",
		Short:         "CLI tool to interact with your ECS clusters",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if opts.debug {
				log.SetLevel(log.DebugLevel)
				log.Debug("debug logs enabled")
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.debug, "debug", false, "Enable debug mode")

	cmd.AddCommand(
		buildEventsCmd(),
		buildImagesCmd(),
		buildInstancesCmd(),
		buildServicesCmd(),
		buildTasksCmd(),
		buildUpdateCmd(),
		buildCompletionCmd(),
	)
	return cmd.Execute()
}
