package helm

import (
	"strings"

	"github.com/conplementAG/copsctl/pkg/adapters/helm"

	"github.com/conplementAG/copsctl/pkg/adapters/kubernetes"
	"github.com/conplementAG/copsctl/pkg/common/logging"
	"github.com/spf13/viper"
)

func Init() {
	namespace := viper.GetString("namespace")

	roleTemplate := renderTemplate(namespace)
	kubernetes.ApplyString(roleTemplate)

	helm.InitHelm("tiller-account", namespace)

	logging.LogSuccess("helm successfully initialized")
}

func Delete() {
	namespace := viper.GetString("namespace")

	kubernetes.DeleteDeployment("tiller-deploy", namespace)

	roleTemplate := renderTemplate(namespace)
	kubernetes.DeleteString(roleTemplate)

	logging.LogSuccess("helm successfully removed")
}

func renderTemplate(namespace string) string {
	template := strings.Replace(helmTemplate, "{{ namespace }}", namespace, -1)
	return template
}
