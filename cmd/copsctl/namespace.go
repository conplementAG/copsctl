package main

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/namespace"
	"github.com/conplementag/cops-hq/pkg/cli"
	"github.com/conplementag/cops-hq/pkg/hq"
)

func createNamespaceCommand(hq hq.HQ) {
	namespaceCmdGroup := hq.GetCli().AddBaseCommand("namespace", "Command group for administration of k8s namespaces",
		"Use this command to administer k8s namespaces.", nil)

	namespaceOrchestrator := namespace.New(hq)

	createNamespaceListCommand(namespaceOrchestrator, namespaceCmdGroup)
	createNamespaceCreateCommand(namespaceOrchestrator, namespaceCmdGroup)
	createNamespaceDeleteCommand(namespaceOrchestrator, namespaceCmdGroup)
	createNamespaceUsersCommand(namespaceOrchestrator, namespaceCmdGroup)
	createNamespaceServiceAccountsCommand(namespaceOrchestrator, namespaceCmdGroup)
}

func createNamespaceListCommand(o *namespace.Orchestrator, cmd cli.Command) {
	cmd.AddCommand("list", "Lists all the CoreOps namespaces", "Use this list all the CoreOps namespaces inside this cluster.", func() {
		o.List()
	})
}

func createNamespaceCreateCommand(o *namespace.Orchestrator, cmd cli.Command) {
	namespaceCreateCmd := cmd.AddCommand("create", "Create a new namespace",
		"Use this command to create a new k8s namespace. This command is idempotent, "+
			"which means you can use it multiple times to ensure that the namespace is there, "+
			"as well as all other users.", func() {
			o.Create()
		})

	addUsersParam(namespaceCreateCmd)
	addServiceAccountsParam(namespaceCreateCmd)
	addNameParam(namespaceCreateCmd)
}

func createNamespaceDeleteCommand(o *namespace.Orchestrator, cmd cli.Command) {
	namespaceDeleteCmd := cmd.AddCommand("delete", "Delete an existing namespace",
		"Use this command to delete existing cops namespace.", func() {
			o.Delete()
		})

	addNameParam(namespaceDeleteCmd)
}

func createNamespaceUsersCommand(o *namespace.Orchestrator, cmd cli.Command) {
	namespaceUserCmd := cmd.AddCommand("users", "Manage users of a namespace",
		"Use this command to manage the users in an existing k8s namespace.",
		nil)

	addNameParam(namespaceUserCmd)

	createNamespaceUsersAddCommand(o, namespaceUserCmd)
	createNamespaceUsersRemoveCommand(o, namespaceUserCmd)
	createNamespaceUsersListCommand(o, namespaceUserCmd)
}

func createNamespaceUsersAddCommand(o *namespace.Orchestrator, cmd cli.Command) {
	namespaceUsersAddCmd := cmd.AddCommand("add", "Adds users to the namespace",
		"Use this command to add users to an existing k8s namespace.", func() {
			o.AddUsers()
		})

	addUsersParam(namespaceUsersAddCmd)
}

func createNamespaceUsersRemoveCommand(o *namespace.Orchestrator, cmd cli.Command) {
	namespaceUsersRemoveCmd := cmd.AddCommand("remove", "Removes users from a namespace",
		"Use this command to remove users from an existing k8s namespace.", func() {
			o.RemoveUsers()
		})

	addUsersParam(namespaceUsersRemoveCmd)
}

func createNamespaceUsersListCommand(o *namespace.Orchestrator, cmd cli.Command) {
	cmd.AddCommand("list", "List users of a namespace",
		"Use this command to list users of an existing k8s namespace.", func() {
			o.ListUsers()
		})
}

func createNamespaceServiceAccountsCommand(o *namespace.Orchestrator, cmd cli.Command) {
	namespaceServiceAccountsCmd := cmd.AddCommand("service-accounts",
		"Manage service accounts of a namespace",
		"Use this command to manage the service accounts in an existing k8s namespace.",
		nil)

	addNameParam(namespaceServiceAccountsCmd)

	createNamespaceServiceAccountsAddCommand(o, namespaceServiceAccountsCmd)
	createNamespaceServiceAccountsRemoveCommand(o, namespaceServiceAccountsCmd)
	createNamespaceServiceAccountsListCommand(o, namespaceServiceAccountsCmd)
}

func createNamespaceServiceAccountsAddCommand(o *namespace.Orchestrator, cmd cli.Command) {
	namespaceServiceAccountAddCmd := cmd.AddCommand("add", "Adds service accounts to the namespace",
		"Use this command to add service accounts to an existing k8s namespace.", func() {
			o.AddServiceAccounts()
		})

	addServiceAccountsParam(namespaceServiceAccountAddCmd)
}

func createNamespaceServiceAccountsRemoveCommand(o *namespace.Orchestrator, cmd cli.Command) {
	namespaceServiceAccountRemoveCmd := cmd.AddCommand("remove", "Removes users from a namespace",
		"Use this command to remove service-accounts from an existing k8s namespace.", func() {
			o.RemoveServiceAccounts()
		})

	addServiceAccountsParam(namespaceServiceAccountRemoveCmd)
}

func createNamespaceServiceAccountsListCommand(o *namespace.Orchestrator, cmd cli.Command) {
	cmd.AddCommand("list", "List service-accounts of a namespace",
		"Use this command to list service-accounts of an existing k8s namespace.", func() {
			o.ListServiceAccounts()
		})
}

func addUsersParam(cmd cli.Command) {
	cmd.AddPersistentParameterString(flags.Users, "", true, "u", "The email-addresses of the admin "+
		"users of the namespace. Must be identical to Azure AD (case-sensitive). "+
		"You can specify multiple users separated by commas.")
}

func addServiceAccountsParam(cmd cli.Command) {
	cmd.AddPersistentParameterString(flags.ServiceAccounts, "", false, "s", "Optionally, you can specify service accounts "+
		"which will be granted identical access level like the users. Each service accounts has to be in "+
		"the format accountname.namespace, and multiple accounts can be specified separated by commas.")
}

func addNameParam(cmd cli.Command) {
	cmd.AddPersistentParameterString(flags.Name, "", true, "n", "Name of the namespace")
}
