package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createHelmCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "helm",
		Short: "DEPRECATED. Use helm2 command instead.",
		Long: `
		DEPRECATED. Use helm2 command instead.
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
