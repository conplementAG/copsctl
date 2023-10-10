package connection

import (
	"strings"
)

const supportedConfigVersion string = "1"
const technicalKubeConfigType string = "technicalAccount"
const developerKubeConfigType string = "developerAccount"

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

func validateKubeConfigsContainer(kubeConfigsContainer KubeConfigsContainerV1, isTechnicalAccount bool) {
	callDaCopsMessage := "Update your copsctl to the latest version, or contact the CoreOps team " +
		"if the latest version does not work."

	if kubeConfigsContainer.Version != supportedConfigVersion {
		panic("This version of copsctl can only work with supported version " + supportedConfigVersion +
			" of CoreOps cluster configurations. " + callDaCopsMessage)
	}

	if (isTechnicalAccount && kubeConfigsContainer.Type != technicalKubeConfigType) ||
		(!isTechnicalAccount && kubeConfigsContainer.Type != developerKubeConfigType) {
		panic("Unsupported type found in the downloaded credentials. " + callDaCopsMessage)
	}
}
