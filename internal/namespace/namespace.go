package namespace

import (
	"fmt"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/hq"
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

func (o *Orchestrator) Exists() {
	namespaceName := viper.GetString(flags.Name)
	result, err := kubernetes.ExistsCopsNamespace(o.executor, namespaceName)

	if err != nil {
		panic("Get Cops namespace failed: " + err.Error())
	}

	if result {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}

func (o *Orchestrator) Create() {
	namespaceName := viper.GetString(flags.Name)

	users := parseUsernames(viper.GetString("users"))
	serviceAccounts := parseServiceAccounts(viper.GetString("service-accounts"))
	copsNamespace := renderTemplate(namespaceName, users, serviceAccounts, viper.GetString(flags.ProjectName), viper.GetString(flags.ProjectCostCenter))

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
		panic("Get Cops namespace failed: " + err.Error())
	}

	if namespace == nil {
		logrus.Infof("Cops namespace '%s' does not exist", namespaceName)
		return
	}

	copsNamespace := renderTemplate(namespaceName, namespace.Spec.NamespaceAdminUsers, namespace.Spec.NamespaceAdminServiceAccounts, namespace.Spec.Project.Name, namespace.Spec.Project.CostCenter)

	_, err = kubernetes.DeleteString(o.executor, copsNamespace)

	if err != nil {
		panic("Deleting Cops namespace failed: " + err.Error())
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
