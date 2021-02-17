package connection

import (
	"gopkg.in/yaml.v2"
	"strings"
)

func validateConnectionString(connectionString string, environmentTag string, isTechnicalAccount bool) {
	expectedToFind := environmentTag + "-developer.yaml"

	if isTechnicalAccount {
		expectedToFind = environmentTag + "-technical-account.yaml"
	}

	if !strings.Contains(connectionString, expectedToFind) {
		panic("Provided connection string is not valid for the specified environment " +
			"or the desired authentication mode.")
	}
}

func validateDownloadedBlob(blob string, isTechnicalAccount bool) {
	bareContainer := getContainerBare(blob)
	supportedVersion := "1"

	callDaCopsMessage := "Update your copsctl to the latest version, or contact the CoreOps team " +
		"if the latest version does not work."

	if bareContainer.Version != supportedVersion {
		panic("This version of copsctl can only work with supported version " + supportedVersion +
			" of CoreOps cluster configurations. " + callDaCopsMessage)
	}

	if (isTechnicalAccount && bareContainer.Type != "technicalAccount") ||
		(!isTechnicalAccount && bareContainer.Type != "developerAccount") {
		panic("Unsupported type found in the downloaded credentials. " + callDaCopsMessage)
	}
}

type KubeConfigsContainerBare struct {
	Version string `yaml:"version"`
	Type    string `yaml:"type"`
}

func getContainerBare(yamlString string) KubeConfigsContainerBare {
	var container KubeConfigsContainerBare

	err := yaml.Unmarshal([]byte(yamlString), &container)

	if err != nil {
		panic(err)
	}

	return container
}
