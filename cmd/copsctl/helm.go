package main

import (
	"github.com/conplementAG/copsctl/pkg/helm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createHelmCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "helm",
		Short: "Command for managing the helm installation in Core Ops namespaces.",
		Long: `
Use this command to manage helm installation in your Core Ops namespaces.
		`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("namespace", cmd.Flags().Lookup("namespace"))
		},
	}

	command.PersistentFlags().StringP("namespace", "n", "", "Namespace where you wish to manage Helm.")
	command.MarkPersistentFlagRequired("namespace")

	command.AddCommand(createHelmInitCommand())
	command.AddCommand(createHelmDeleteCommand())

	return command
}

func createHelmInitCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "init",
		Short: "Initialize Helm",
		Long:  "Use this command to initialize Helm in a namespace of your choice.",
		Run: func(cmd *cobra.Command, args []string) {
			helm.Init()
		},
	}

	return command
}

func createHelmDeleteCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "delete",
		Short: "Remove Helm",
		Long:  "Use this command to remove Helm from your namespace.",
		Run: func(cmd *cobra.Command, args []string) {
			helm.Delete()
		},
	}

	return command
}
