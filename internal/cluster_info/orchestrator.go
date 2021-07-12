package cluster_info

import (
	"fmt"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/common/commands"
	"github.com/conplementAG/copsctl/internal/common/logging"
	"github.com/spf13/viper"
	"strings"
)

func ShowClusterInfo() {
	printConfigSilenceEverythingElse := viper.GetBool(flags.PrintToStdoutSilenceEverythingElse)

	if !printConfigSilenceEverythingElse {
		logging.Info("Reading the cluster info ..." )
		logging.Info("NOTE: you can use the " + flags.PrintToStdoutSilenceEverythingElse + " flag to silence these outputs (useful for automation)" )

		logging.Info("===========================================================")
		logging.Info("==================== Cluster Info:  =======================")
		logging.Info("===========================================================")
	}

	command := "kubectl get configmap -n coreops-public -o jsonpath='{.data}' coreops-cluster-info"
	result, err := commands.ExecuteCommand(commands.Create(command))

	if err != nil {
		logging.Errorf(err.Error())
		panic(err)
	}

	result = strings.TrimPrefix(result, "'")
	result = strings.TrimSuffix(result, "'")
	fmt.Println(result)
}