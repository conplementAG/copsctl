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
	_, err := kubernetes.ApplyString(roleTemplate)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	_, err = helm.InitHelm("tiller-account", namespace)

	if err != nil {
		panic("Helm init failed: " + err.Error())
	}

	logging.Info("helm successfully initialized")
}

func Delete() {
	namespace := viper.GetString("namespace")

	err := kubernetes.DeleteDeployment("tiller-deploy", namespace)

	if err != nil {
		panic("Deleting tiller deployment failed: " + err.Error())
	}

	roleTemplate := renderTemplate(namespace)
	_, err = kubernetes.DeleteString(roleTemplate)

	if err != nil {
		panic("Deleting role template failed: " + err.Error())
	}

	logging.Info("helm successfully removed")
}

func renderTemplate(namespace string) string {
	template := strings.Replace(helmTemplate, "{{ namespace }}", namespace, -1)
	return template
}
