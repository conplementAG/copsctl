package connection

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func marshalToYaml(input interface{}) string {
	output, err := yaml.Marshal(&input)

	if err != nil {
		logrus.Error("Error occurred while generating YAML for the interface. " + err.Error())
		panic(err)
	}

	return string(output)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
