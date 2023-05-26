package main

import (
	"fmt"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/info"
	"github.com/conplementag/cops-hq/v2/pkg/cli"
	"github.com/conplementag/cops-hq/v2/pkg/hq"
)

func createInfoCommands(hq hq.HQ) {
	infoCmdGroup := hq.GetCli().AddBaseCommand("info", "Command group for informations of k8s",
		"Use this command to get informations on cluster / environment.", nil)

	orchestrator := info.New(hq)

	createClusterInfoCommand(hq, orchestrator)
	createInfoClusterCommand(infoCmdGroup, orchestrator)
	createInfoEnvironementCommand(infoCmdGroup, orchestrator)
}

func createClusterInfoCommand(hq hq.HQ, o *info.Orchestrator) {
	clusterInfoCmdGroup := hq.GetCli().AddBaseCommand("cluster-info", "[DEPRECATED] Command for showing the CoreOps cluster information",
		fmt.Sprintf("[DEPRECATED] Use this command to get the cluster info which might be useful for your. For example, if the static outbound IPs are enabled for the cluster, then you can use this command to get these IPs. Make sure you are connected to the cluster first. "+
			"Use the %s flag for automation.", flags.PrintToStdoutSilenceEverythingElse), func() {
			o.ShowEnvironmentInfo()
		})

	addSilenceParam(clusterInfoCmdGroup)
}

func createInfoClusterCommand(cmd cli.Command, o *info.Orchestrator) {
	infoClusterCmd := cmd.AddCommand("cluster", "Get cluster infos",
		"Use this command to get informations around cluster which have the same lifecycle as the cluster its self. "+
			"Check 'info environment' for information with common lifecycle. ", func() {
			o.ShowClusterInfo()
		})
	addSilenceParam(infoClusterCmd)
}

func createInfoEnvironementCommand(cmd cli.Command, o *info.Orchestrator) {
	infoClusterCmd := cmd.AddCommand("environment", "Get environment infos",
		"Use this command to get informations around environment with common lifecycle. "+
			"Check 'info cluster' for cluster specific information with cluster lifecycle. ", func() {
			o.ShowEnvironmentInfo()
		})
	addSilenceParam(infoClusterCmd)
}

func addSilenceParam(cmd cli.Command) {
	cmd.AddPersistentParameterBool(flags.PrintToStdoutSilenceEverythingElse, false, false, "q",
		"Similar to print-to-stdout, but silences all other logging outputs. Useful for automation.")

}
