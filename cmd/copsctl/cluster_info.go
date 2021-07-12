package main

import (
	"github.com/conplementAG/copsctl/internal/cluster_info"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createClusterInfoCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "cluster-info",
		Short: "Command for showing the CoreOps cluster information",
		Long: `
Use this command to get the cluster info which might be useful for your. For example, if the static outbound IPs are enabled for the cluster,
then you can use this command to get these IPs. Make sure you are connected to the cluster first.
` + "Use the " + flags.PrintToStdoutSilenceEverythingElse + " flag for automation.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.PrintToStdoutSilenceEverythingElse,
				cmd.Flags().Lookup(flags.PrintToStdoutSilenceEverythingElse))
		},
		Run: func(cmd *cobra.Command, args []string) {
			cluster_info.ShowClusterInfo()
		},
	}

	command.PersistentFlags().BoolP(flags.PrintToStdoutSilenceEverythingElse, "q", false,
		"Similar to print-to-stdout, but silences all other logging outputs. Useful for automation.")

	return command
}
