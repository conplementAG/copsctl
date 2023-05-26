package cluster_info

import (
	"encoding/json"
	"fmt"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type Orchestrator struct {
	hq       hq.HQ
	executor commands.Executor
}

func New(hq hq.HQ) *Orchestrator {
	return &Orchestrator{
		hq:       hq,
		executor: hq.GetExecutor(),
	}
}

func (o *Orchestrator) ShowClusterInfo() {
	printConfigSilenceEverythingElse := viper.GetBool(flags.PrintToStdoutSilenceEverythingElse)

	if !printConfigSilenceEverythingElse {
		logrus.Info("Reading the cluster info ...")
		logrus.Info("NOTE: you can use the " + flags.PrintToStdoutSilenceEverythingElse + " flag to silence these outputs (useful for automation)")

		logrus.Info("===========================================================")
		logrus.Info("==================== Cluster Info:  =======================")
		logrus.Info("===========================================================")
	}

	result, err := o.executor.Execute("kubectl get configmap -n coreops-public -o jsonpath=\"{.data['info\\.json']}\" coreops-cluster-info")

	if err != nil {
		logrus.Errorf(err.Error())
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
			logrus.Errorf(err.Error())
			panic(err)
		}

		fmt.Println(string(indented))
	}
}
