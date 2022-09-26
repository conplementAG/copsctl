package namespace

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/conplementag/cops-hq/pkg/hq"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/conplementAG/copsctl/internal/adapters/kubernetes"
	"github.com/spf13/viper"
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

func (o *Orchestrator) List() {
	kubernetes.PrintAllCopsNamespaces(o.executor)
}

func (o *Orchestrator) Create() {
	namespaceName := viper.GetString(flags.Name)

	users := parseUsernames(viper.GetString("users"))
	serviceAccounts := parseServiceAccounts(viper.GetString("service-accounts"))
	copsNamespace := renderTemplate(namespaceName, users, serviceAccounts)

	_, err := kubernetes.ApplyString(o.executor, copsNamespace)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	ensureNamespaceAccess(o.executor, namespaceName)

	logrus.Infof("Cops namespace %s successfully created / updated.", namespaceName)
}

func (o *Orchestrator) Delete() {
	namespaceName := viper.GetString(flags.Name)

	namespace, err := kubernetes.GetCopsNamespace(o.executor, namespaceName)

	if err != nil {
		logrus.Infof("Cops namespace %s does not exist", namespaceName)
		return
	}

	copsNamespace := renderTemplate(namespaceName, namespace.Spec.NamespaceAdminUsers, namespace.Spec.NamespaceAdminServiceAccounts)

	_, error := kubernetes.DeleteString(o.executor, copsNamespace)

	if error != nil {
		panic("Deleting copsnamespace failed: " + err.Error())
	}

	logrus.Infof("Cops namespace %s successfully deleted", namespaceName)
}

func ensureNamespaceAccess(executor commands.Executor, namespace string) {
	status := false

	for i := 0; i < 20; i++ {
		status = kubernetes.CanIGetPods(executor, namespace)
		if status == true {
			break
		}
		time.Sleep(3 * time.Second)
	}

	if status == false {
		panic("Could not verify access to pods in created namespace.")
	}
}
