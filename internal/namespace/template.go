package namespace

import (
	"strings"

	"github.com/conplementAG/copsctl/internal/adapters/kubernetes"
)

func renderUsernames(userNames []string) string {
	userlist := ""
	length := len(userNames)

	if len(userNames) > 0 {
		userlist += "  namespaceAdminUsers:\n"

		for index, userName := range userNames {
			userlist += "  - " + userName

			if index != (length - 1) {
				userlist += "\n"
			}
		}
	}

	return userlist
}

func renderServiceAccounts(serviceAccounts []kubernetes.CopsServiceAccount) string {
	accountsList := ""

	if len(serviceAccounts) > 0 {
		accountsList += "  namespaceAdminServiceAccounts:\n"

		for index, account := range serviceAccounts {
			accountsList += "  - serviceAccount: " + account.ServiceAccount + "\n"
			accountsList += "    namespace: " + account.Namespace

			if index != (len(serviceAccounts) - 1) {
				accountsList += "\n"
			}
		}
	}

	return accountsList
}

func renderTemplate(namespaceName string, userNames []string, serviceAccounts []kubernetes.CopsServiceAccount) string {
	copsnamespace := strings.Replace(copsNamespaceTemplate, "{{ namespaceName }}", namespaceName, -1)
	copsnamespace = strings.Replace(copsnamespace, "{{ usernames }}", renderUsernames(userNames), -1)
	copsnamespace = strings.Replace(copsnamespace, "{{ serviceAccounts }}", renderServiceAccounts(serviceAccounts), -1)
	return copsnamespace
}

const copsNamespaceTemplate string = `
apiVersion: coreops.conplement.cloud/v1
kind: CopsNamespace
metadata:
  name: {{ namespaceName }}
spec:
{{ usernames }}
{{ serviceAccounts }}`
