package main

import (
	"github.com/conplementAG/copsctl/internal/helm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createHelm2Command() *cobra.Command {
	var command = &cobra.Command{
		Use:   "helm2",
		Short: "Command for managing the Helm 2 installation in Core Ops namespaces.",
		Long: `
Use this command to manage Helm 2 installation in your Core Ops namespaces. NOTE: for these commands to work you need the helm 2 binary in your path, available as "helm" command!
		`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("namespace", cmd.Flags().Lookup("namespace"))
		},
	}

	command.PersistentFlags().StringP("namespace", "n", "", "Namespace where you wish to manage Helm.")
	command.MarkPersistentFlagRequired("namespace")

	command.AddCommand(createHelm2InitCommand())
	command.AddCommand(createHelm2DeleteCommand())

	return command
}

func createHelm2InitCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "init",
		Short: "Initialize Helm 2 server side component (tiller)",
		Long:  "Use this command to initialize Helm 2 in a namespace of your choice.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("namespace", cmd.Flags().Lookup("namespace"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			helm.Init()
		},
	}

	return command
}

func createHelm2DeleteCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "delete",
		Short: "Remove Helm",
		Long:  "Use this command to remove Helm 2 server side component (tiller) from your namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("namespace", cmd.Flags().Lookup("namespace"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			helm.Delete()
		},
	}

	return command
}
