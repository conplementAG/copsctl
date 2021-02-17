package connection

import (
	"fmt"
	"github.com/conplementAG/copsctl/internal/common/logging"
	"gopkg.in/yaml.v2"
)

func marshalToYaml(input interface{}) string {
	output, err := yaml.Marshal(&input)

	if err != nil {
		logging.Error("Error occurred while generating YAML for the interface. " + err.Error())
		panic(err)
	}

	return string(output)
}

func confirmOperation(message string) {
	logging.Info(message)
	var input string
	fmt.Scanln(&input)

	if input != "yes" {
		panic("operation aborted")
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
