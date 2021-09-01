package cluster_info

import (
	"encoding/json"
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

	if printConfigSilenceEverythingElse {
		// no pretty formats or anything is needed if printing for parse purposes
		fmt.Println(result)
	} else {
		var mapResult map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(result), &mapResult)
		
		indented, err := json.MarshalIndent(mapResult, "", "    ")

		if err != nil {
			logging.Errorf(err.Error())
			panic(err)
		}

		fmt.Println(string(indented))
	}
}