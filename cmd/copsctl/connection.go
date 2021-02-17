package main

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/connection"
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
			viper.BindPFlag(flags.EnvironmentTag, cmd.Flags().Lookup(flags.EnvironmentTag))
			viper.BindPFlag(flags.ConnectionString, cmd.Flags().Lookup(flags.ConnectionString))
			viper.BindPFlag(flags.TechnicalAccount, cmd.Flags().Lookup(flags.TechnicalAccount))
			viper.BindPFlag(flags.Secondary, cmd.Flags().Lookup(flags.Secondary))

			viper.BindPFlag(flags.AutoApprove, cmd.Flags().Lookup(flags.AutoApprove))
			viper.BindPFlag(flags.PrintToStdout, cmd.Flags().Lookup(flags.PrintToStdout))
			viper.BindPFlag(flags.PrintToStdoutSilenceEverythingElse,
				cmd.Flags().Lookup(flags.PrintToStdoutSilenceEverythingElse))
		},
		Run: func(cmd *cobra.Command, args []string) {
			connection.Connect()
		},
	}

	command.PersistentFlags().StringP(flags.EnvironmentTag, "e", "",
		"The environment tag of the environment you want to connect to")
	command.MarkPersistentFlagRequired(flags.EnvironmentTag)

	command.PersistentFlags().StringP(flags.ConnectionString, "c", "",
		"The connection string for the environment.")
	command.MarkPersistentFlagRequired(flags.ConnectionString)

	command.PersistentFlags().BoolP(flags.TechnicalAccount, "t", false,
		"Use to switch to non-interactive technical account mode. "+
			"Make sure that the connection string you provide is valid for a technical account!")
	command.PersistentFlags().BoolP(flags.Secondary, "s", false,
		"Connect to the secondary instead of the primary cluster.")

	command.PersistentFlags().BoolP(flags.AutoApprove, "a", false,
		"Auto approve any confirmation prompts.")
	command.PersistentFlags().BoolP(flags.PrintToStdout, "p", false,
		"Print the config to stdout instead of overwriting $HOME/.kube/config. Useful for developers")
	command.PersistentFlags().BoolP(flags.PrintToStdoutSilenceEverythingElse, "q", false,
		"Similar to print-to-stdout, but silences all other logging outputs. Useful for automation.")

	return command
}
