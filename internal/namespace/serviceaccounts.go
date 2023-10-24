package namespace

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/conplementAG/copsctl/internal/adapters/kubernetes"
	"github.com/spf13/viper"
)

func parseServiceAccounts(rawAccounts string) []kubernetes.CopsServiceAccount {
	parsedAccounts := make([]kubernetes.CopsServiceAccount, 0)

	if rawAccounts != "" {
		accounts := strings.Split(rawAccounts, ",")

		for _, account := range accounts {
			if !strings.Contains(account, ".") {
				panic("Service accounts should contain dot (.) separated service account(s) (e.g. accountname.namespace). " +
					"Input value: " + account)
			}

			splitAccount := strings.Split(account, ".")

			if len(splitAccount) != 2 || splitAccount[0] == "" || splitAccount[1] == "" {
				panic("Service account not in the expected format: " + account)
			}

			parsedAccounts = append(parsedAccounts, kubernetes.CopsServiceAccount{
				ServiceAccount: splitAccount[0],
				Namespace:      splitAccount[1],
			})
		}
	}

	return parsedAccounts
}

func (o *Orchestrator) AddServiceAccounts() {
	namespaceName := viper.GetString(flags.Name)
	accounts := viper.GetString("service-accounts")

	newAccounts := parseServiceAccounts(accounts)
	namespace, err := kubernetes.GetCopsNamespace(o.executor, namespaceName)

	if err != nil {
		panic("Could not get the cops namespace " + err.Error())
	}

	existingUsers := namespace.Spec.NamespaceAdminUsers
	relevantAccounts := namespace.Spec.NamespaceAdminServiceAccounts

	addedAccountsCount := 0

	for _, newAccount := range newAccounts {
		alreadyExists := false

		for _, existingAccount := range relevantAccounts {
			if existingAccount.ServiceAccount == newAccount.ServiceAccount && existingAccount.Namespace == newAccount.Namespace {
				alreadyExists = true
				break
			}
		}

		if !alreadyExists {
			relevantAccounts = append(relevantAccounts, newAccount)
			addedAccountsCount++
		}
	}

	copsNamespace := renderTemplate(namespaceName, existingUsers, relevantAccounts, namespace.Spec.Project.Name, namespace.Spec.Project.CostCenter)

	_, err = kubernetes.ApplyString(o.executor, copsNamespace)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	logrus.Infof("%d service account(s) have been successfully added to %s namespace", addedAccountsCount, namespaceName)
}

func (o *Orchestrator) RemoveServiceAccounts() {
	namespaceName := viper.GetString(flags.Name)
	accounts := viper.GetString("service-accounts")

	accountsToRemove := parseServiceAccounts(accounts)
	namespace, err := kubernetes.GetCopsNamespace(o.executor, namespaceName)

	if err != nil {
		panic("Could not get the cops namespace " + err.Error())
	}

	existingUsers := namespace.Spec.NamespaceAdminUsers
	existingAccounts := namespace.Spec.NamespaceAdminServiceAccounts
	var resultingAccounts []kubernetes.CopsServiceAccount

	removedAccountsCount := 0

	for _, existingAccount := range existingAccounts {
		shouldAccountBeRemoved := false

		for _, account := range accountsToRemove {
			if account.ServiceAccount == existingAccount.ServiceAccount && account.Namespace == existingAccount.Namespace {
				shouldAccountBeRemoved = true
				removedAccountsCount++
				break
			}
		}

		if !shouldAccountBeRemoved {
			resultingAccounts = append(resultingAccounts, existingAccount)
		}
	}

	copsnamespace := renderTemplate(namespaceName, existingUsers, resultingAccounts, namespace.Spec.Project.Name, namespace.Spec.Project.CostCenter)

	_, err = kubernetes.ApplyString(o.executor, copsnamespace)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	logrus.Infof("%d service account(s) have been removed from the %s namespace", removedAccountsCount, namespaceName)
}

func (o *Orchestrator) ListServiceAccounts() {
	namespaceName := viper.GetString(flags.Name)
	namespace, err := kubernetes.GetCopsNamespace(o.executor, namespaceName)

	if err != nil {
		panic("Could not get the cops namespace " + err.Error())
	}

	serviceAccounts := namespace.Spec.NamespaceAdminServiceAccounts

	logrus.Info("Current service accounts in the namespace " + namespaceName + ":")

	for _, sa := range serviceAccounts {
		logrus.Info(" - " + sa.ServiceAccount + "." + sa.Namespace)
	}
}
