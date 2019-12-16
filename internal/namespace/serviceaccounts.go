package namespace

import (
	"strings"

	"github.com/conplementAG/copsctl/internal/adapters/kubernetes"
	"github.com/conplementAG/copsctl/internal/common/logging"
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

func AddServiceAccounts() {
	namespaceName := viper.GetString("name")
	accounts := viper.GetString("service-accounts")

	newAccounts := parseServiceAccounts(accounts)
	namespace, err := kubernetes.GetCopsNamespace(namespaceName)

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

	copsnamespace := renderTemplate(namespaceName, existingUsers, relevantAccounts)

	_, err = kubernetes.ApplyString(copsnamespace)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	logging.Infof("%d service account(s) have been successfully added to %s namespace", addedAccountsCount, namespaceName)
}

func RemoveServiceAccounts() {
	namespaceName := viper.GetString("name")
	accounts := viper.GetString("service-accounts")

	accountsToRemove := parseServiceAccounts(accounts)
	namespace, err := kubernetes.GetCopsNamespace(namespaceName)

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

	copsnamespace := renderTemplate(namespaceName, existingUsers, resultingAccounts)

	_, err = kubernetes.ApplyString(copsnamespace)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	logging.Infof("%d service account(s) have been removed from the %s namespace", removedAccountsCount, namespaceName)
}

func ListServiceAccounts() {
	namespaceName := viper.GetString("name")
	namespace, err := kubernetes.GetCopsNamespace(namespaceName)

	if err != nil {
		panic("Could not get the cops namespace " + err.Error())
	}

	serviceAccounts := namespace.Spec.NamespaceAdminServiceAccounts

	logging.Info("Current service accounts in the namespace " + namespaceName + ":")

	for _, sa := range serviceAccounts {
		logging.Info(" - " + sa.ServiceAccount + "." + sa.Namespace)
	}
}
