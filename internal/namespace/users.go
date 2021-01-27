package namespace

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"strings"

	"github.com/conplementAG/copsctl/internal/adapters/kubernetes"
	"github.com/conplementAG/copsctl/internal/common/logging"
	"github.com/spf13/viper"
)

func parseUsernames(userNames string) []string {
	var parsedUsers []string
	users := strings.Split(userNames, ",")

	for _, username := range users {
		parsedUsers = append(parsedUsers, username)
	}

	return parsedUsers
}

func AddUsers() {
	namespaceName := viper.GetString(flags.Name)
	users := viper.GetString("users")

	newUsers := parseUsernames(users)
	namespace, err := kubernetes.GetCopsNamespace(namespaceName)

	if err != nil {
		panic("Could not get the cops namespace " + err.Error())
	}

	relevantUsers := namespace.Spec.NamespaceAdminUsers
	existingServiceAccounts := namespace.Spec.NamespaceAdminServiceAccounts

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

	copsnamespace := renderTemplate(namespaceName, relevantUsers, existingServiceAccounts)

	_, err = kubernetes.ApplyString(copsnamespace)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	logging.Infof("%d user(s) have been successfully added to %s namespace", addedUserCount, namespaceName)
}

func RemoveUsers() {
	namespaceName := viper.GetString(flags.Name)
	users := viper.GetString("users")

	usersToRemove := parseUsernames(users)
	namespace, err := kubernetes.GetCopsNamespace(namespaceName)

	if err != nil {
		panic("Could not get the cops namespace " + err.Error())
	}

	existingUsers := namespace.Spec.NamespaceAdminUsers
	existingServiceAccounts := namespace.Spec.NamespaceAdminServiceAccounts
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

	copsnamespace := renderTemplate(namespaceName, relevantUsers, existingServiceAccounts)

	_, err = kubernetes.ApplyString(copsnamespace)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	logging.Infof("%d user(s) have been successfully removed from %s namespace", removedUserCount, namespaceName)
}

func ListUsers() {
	namespaceName := viper.GetString(flags.Name)
	namespace, err := kubernetes.GetCopsNamespace(namespaceName)

	if err != nil {
		panic("Could not get the cops namespace " + err.Error())
	}

	users := namespace.Spec.NamespaceAdminUsers

	logging.Info("Current users in namespace " + namespaceName + ":")

	for _, user := range users {
		logging.Info(" - " + user)
	}
}
