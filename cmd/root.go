package cmd

import (
	"fmt"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

var (
	version string = "dev"
	commit  string
)

type rootOpts struct {
	debug bool
}

func Execute() error {
	log.SetHandler(cli.Default)
	opts := rootOpts{}

	rootCmd := &cobra.Command{
		Use:           "ecs",
		Short:         "CLI tool to interact with your ECS clusters",
		Version:       buildVersion(version, commit),
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if opts.debug {
				log.SetLevel(log.DebugLevel)
				log.Debug("debug logs enabled")
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&opts.debug, "debug", false, "Enable debug mode")

	rootCmd.AddCommand(
		buildEventsCmd(),
		buildImagesCmd(),
		buildInstancesCmd(),
		buildServicesCmd(),
		buildTasksCmd(),
		buildUpdateCmd(),
	)
	return rootCmd.Execute()
}

func buildVersion(version, commit string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	return result
}
