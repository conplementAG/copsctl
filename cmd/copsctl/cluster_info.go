package main

import (
	"fmt"
	"github.com/conplementAG/copsctl/internal/cluster_info"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementag/cops-hq/pkg/hq"
)

func createClusterInfoCommand(hq hq.HQ) {
	clusterInfoCmdGroup := hq.GetCli().AddBaseCommand("cluster-info", "Command for showing the CoreOps cluster information",
		fmt.Sprintf("Use this command to get the cluster info which might be useful for your. For example, if the static outbound IPs are enabled for the cluster, then you can use this command to get these IPs. Make sure you are connected to the cluster first. "+
			"Use the %s flag for automation.", flags.PrintToStdoutSilenceEverythingElse), func() {
			cluster_info.New(hq).ShowClusterInfo()
		})

	clusterInfoCmdGroup.AddPersistentParameterBool(flags.PrintToStdoutSilenceEverythingElse, false, false, "q",
		"Similar to print-to-stdout, but silences all other logging outputs. Useful for automation.")
}
