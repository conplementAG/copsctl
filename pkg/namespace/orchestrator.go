package namespace

import (
	"strings"
	"time"

	"github.com/conplementAG/copsctl/pkg/adapters/kubernetes"
	"github.com/conplementAG/copsctl/pkg/common/fileprocessing"
	"github.com/conplementAG/copsctl/pkg/common/logging"
	"github.com/spf13/viper"
)

// Create creates a CopsNamespace Custom-Resource-Definition with the given name and user
func Create() {
	namespaceName := viper.GetString("name")
	userNames := viper.GetString("users")

	users := parseUsernames(userNames)
	copsnamespace := renderTemplate(namespaceName, users)

	temporaryDirectory, temporaryFile := fileprocessing.WriteStringToTemporaryFile(copsnamespace, "copsnamespace.yaml")
	kubernetes.Apply(temporaryFile)
	fileprocessing.DeletePath(temporaryDirectory)

	ensureNamespaceAccess(namespaceName)

	logging.LogSuccessf("Cops namespace %s successfully created\n", namespaceName)
}

func renderTemplate(namespaceName string, userNames []string) string {
	copsnamespace := strings.Replace(copsNamespaceTemplate, "{{ namespaceName }}", namespaceName, -1)
	copsnamespace = strings.Replace(copsnamespace, "{{ usernames }}", renderUsernames(userNames), -1)
	return copsnamespace
}

func parseUsernames(userNames string) []string {
	var parsedUsers []string
	users := strings.Split(userNames, ",")
	for _, username := range users {
		parsedUsers = append(parsedUsers, username)
	}
	return parsedUsers
}

func renderUsernames(userNames []string) string {
	userlist := ""
	length := len(userNames)
	for index, userName := range userNames {
		userlist += "  - " + userName
		if index != (length - 1) {
			userlist += "\n"
		}
	}
	return userlist
}

func ensureNamespaceAccess(namespace string) {
	status := false
	for i := 0; i < 10; i++ {
		status = kubernetes.CanIGetPods(namespace)
		if status == true {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if status == false {
		panic("Could not verify access to pods in created namespace.")
	}
}
