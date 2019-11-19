package namespace

import (
	"time"

	"github.com/conplementAG/copsctl/internal/adapters/kubernetes"
	"github.com/conplementAG/copsctl/internal/common/logging"
	"github.com/spf13/viper"
)

func List() {
	kubernetes.PrintAllCopsNamespaces()
}

func Create() {
	namespaceName := viper.GetString("name")

	users := parseUsernames(viper.GetString("users"))
	serviceAccounts := parseServiceAccounts(viper.GetString("service-accounts"))
	copsnamespace := renderTemplate(namespaceName, users, serviceAccounts)

	_, err := kubernetes.ApplyString(copsnamespace)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	ensureNamespaceAccess(namespaceName)

	logging.Infof("Cops namespace %s successfully created / updated.", namespaceName)
}

func Delete() {
	namespaceName := viper.GetString("name")

	namespace, err := kubernetes.GetCopsNamespace(namespaceName)

	if err != nil {
		logging.Infof("Cops namespace %s does not exist", namespaceName)
		return
	}

	copsnamespace := renderTemplate(namespaceName, namespace.Spec.NamespaceAdminUsers, namespace.Spec.NamespaceAdminServiceAccounts)

	_, error := kubernetes.DeleteString(copsnamespace)

	if error != nil {
		panic("Deleting copsnamespace failed: " + err.Error())
	}

	logging.Infof("Cops namespace %s successfully deleted", namespaceName)
}

func ensureNamespaceAccess(namespace string) {
	status := false

	for i := 0; i < 20; i++ {
		status = kubernetes.CanIGetPods(namespace)
		if status == true {
			break
		}
		time.Sleep(3 * time.Second)
	}

	if status == false {
		panic("Could not verify access to pods in created namespace.")
	}
}
