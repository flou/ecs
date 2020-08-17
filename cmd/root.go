package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var awsRegion string
var version = "0.0.8"
var revision string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ecs",
	Short:   "ECS Tools",
	Long:    "Command line tools to interact with your ECS clusters",
	Version: fmt.Sprintf("%s (%s)", version, revision),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&awsRegion, "region", "", "AWS region")
}
