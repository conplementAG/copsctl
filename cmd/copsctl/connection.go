package main

import (
	"github.com/conplementAG/copsctl/pkg/connection"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createConnectCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "connect",
		Short: "Command for managing the connection to k8s clusters",
		Long: `
Use this command to manage the connection to a k8s cluster.
		`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("environment-tag", cmd.Flags().Lookup("environment-tag"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			connection.Connect()
		},
	}

	command.PersistentFlags().StringP("environment-tag", "e", "", "The environment tag of the cluster you want to connect to")
	command.MarkPersistentFlagRequired("environment-tag")

	return command
}
