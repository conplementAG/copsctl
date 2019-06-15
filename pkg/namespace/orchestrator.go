package namespace

import (
	"log"
	"strings"
	"time"

	"github.com/conplementAG/copsctl/pkg/adapters/kubernetes"
	"github.com/conplementAG/copsctl/pkg/common/logging"
	"github.com/spf13/viper"
)

// List simply lists all of the coreops namespaces
func List() {
	kubernetes.PrintAllCopsNamespaces()
}

// Create creates a CopsNamespace Custom-Resource-Definition with the given name and user
func Create() {
	namespaceName := viper.GetString("name")
	userNames := viper.GetString("users")

	users := parseUsernames(userNames)
	copsnamespace := renderTemplate(namespaceName, users)

	kubernetes.ApplyString(copsnamespace)

	ensureNamespaceAccess(namespaceName)

	logging.LogSuccessf("Cops namespace %s successfully created\n", namespaceName)
}

// AddUsers adds the given users to the clusters
func AddUsers() {
	namespaceName := viper.GetString("name")
	users := viper.GetString("users")

	newUsers := parseUsernames(users)
	namespace := kubernetes.GetCopsNamespace(namespaceName)
	relevantUsers := namespace.Spec.NamespaceAdminUsers

	addedUserCount := 0

	for _, newUser := range newUsers {
		userAlreadyExists := false

		for _, existingUser := range relevantUsers {
			if existingUser == newUser {
				userAlreadyExists = true
				break
			}
		}

		if !userAlreadyExists {
			relevantUsers = append(relevantUsers, newUser)
			addedUserCount++
		}
	}

	copsnamespace := renderTemplate(namespaceName, relevantUsers)
	kubernetes.ApplyString(copsnamespace)

	logging.LogSuccessf("%d user(s) have been successfully added to %s namespace\n", addedUserCount, namespaceName)
}

// RemoveUsers removes the given users from the clusters
func RemoveUsers() {
	namespaceName := viper.GetString("name")
	users := viper.GetString("users")

	usersToRemove := parseUsernames(users)
	namespace := kubernetes.GetCopsNamespace(namespaceName)
	existingUsers := namespace.Spec.NamespaceAdminUsers
	var relevantUsers []string

	removedUserCount := 0

	for _, existingUser := range existingUsers {
		shouldUserBeRemoved := false

		for _, userToRemove := range usersToRemove {
			if userToRemove == existingUser {
				shouldUserBeRemoved = true
				removedUserCount++
				break
			}
		}

		if !shouldUserBeRemoved {
			relevantUsers = append(relevantUsers, existingUser)
		}
	}

	copsnamespace := renderTemplate(namespaceName, relevantUsers)
	kubernetes.ApplyString(copsnamespace)

	logging.LogSuccessf("%d user(s) have been successfully removed from %s namespace\n", removedUserCount, namespaceName)
}

// ListUsers prints the current users of the given namespace
func ListUsers() {
	namespaceName := viper.GetString("name")
	namespace := kubernetes.GetCopsNamespace(namespaceName)
	users := namespace.Spec.NamespaceAdminUsers

	log.Println("Current users in namespace " + namespaceName + ":")

	for _, user := range users {
		log.Println(" - " + user)
	}
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
