package helm

import (
	"github.com/conplementAG/copsctl/pkg/adapters/kubernetes"
	"github.com/conplementAG/copsctl/pkg/common/logging"
	"github.com/spf13/viper"
)

func Init() {
	namespace := viper.GetString("namespace")
	kubernetes.CreateServiceAccount(namespace, "tiller-account")

	logging.LogSuccess("helm successfully initialized")
}

func Delete() {
	namespace := viper.GetString("namespace")
	kubernetes.RemoveServiceAccount(namespace, "tiller-account")

	logging.LogSuccess("helm successfully removed")
}
