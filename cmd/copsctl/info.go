package main

import (
	"fmt"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/info"
	"github.com/conplementag/cops-hq/v2/pkg/cli"
	"github.com/conplementag/cops-hq/v2/pkg/hq"
)

func createInfoCommands(hq hq.HQ) {
	infoCmdGroup := hq.GetCli().AddBaseCommand("info",
		"Command group for information around the CoreOps environments and clusters.",
		"Use this command to get information on environments and clusters. Environments are defined as "+
			"sort-of-containers for both primary and secondary clusters, including all of the common resources. Clusters "+
			"and their information is specific to a cluster, for example calling 'info cluster' on a secondary cluster does "+
			"not show the same results like on the primary.", nil)

	orchestrator := info.New(hq)

	createClusterInfoCommand(hq, orchestrator)
	createInfoClusterCommand(infoCmdGroup, orchestrator)
	createInfoEnvironmentCommand(infoCmdGroup, orchestrator)
}

func createClusterInfoCommand(hq hq.HQ, o *info.Orchestrator) {
	clusterInfoCmdGroup := hq.GetCli().AddBaseCommand("cluster-info",
		"[DEPRECATED] Command for showing the CoreOps cluster information",
		fmt.Sprintf("[DEPRECATED] Use this command to get the cluster info which might be useful for your. "+
			"For example, if the static outbound IPs are enabled for the cluster, then you can use this command to get "+
			"these IPs. Make sure you are connected to the cluster first. Use the %s flag for automation.",
			flags.PrintToStdoutSilenceEverythingElse),
		func() {
			o.ShowEnvironmentInfo()
		})

	addSilenceParam(clusterInfoCmdGroup)
}

func createInfoClusterCommand(cmd cli.Command, o *info.Orchestrator) {
	infoClusterCmd := cmd.AddCommand("cluster", "Get cluster infos",
		"Use this command to get information around cluster which have the same lifecycle as the cluster itself "+
			"(e.g. primary only or secondary only information). Check 'info environment' for information with common lifecycle. ",
		func() {
			o.ShowClusterInfo()
		})
	addSilenceParam(infoClusterCmd)
}

func createInfoEnvironmentCommand(cmd cli.Command, o *info.Orchestrator) {
	infoClusterCmd := cmd.AddCommand("environment", "Get environment infos",
		"Use this command to get information around the entire environment, applicable to both primary and "+
			"secondary clusters in it. Check 'info cluster' for cluster specific information with cluster lifecycle.",
		func() {
			o.ShowEnvironmentInfo()
		})
	addSilenceParam(infoClusterCmd)
}

func addSilenceParam(cmd cli.Command) {
	cmd.AddPersistentParameterBool(flags.PrintToStdoutSilenceEverythingElse, false, false, "q",
		"Similar to print-to-stdout, but silences all other logging outputs. Useful for automation.")

}
