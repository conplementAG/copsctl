package main

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"os"

	"github.com/conplementAG/copsctl/internal/namespace"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createNamespaceCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "namespace",
		Short: "Command group for administration of k8s namespaces",
		Long:  "Use this command to administer k8s namespaces.",
		Run: func(cmd *cobra.Command, args []string) {
			// since "namespace" is not really a command, but rather a group of commands, we
			// show the help for the command group instead
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	command.AddCommand(createNamespaceListCommand())
	command.AddCommand(createNamespaceCreateCommand())
	command.AddCommand(createNamespaceDeleteCommand())
	command.AddCommand(createNamespaceUsersCommand())
	command.AddCommand(createNamespaceServiceAccountsCommand())

	return command
}

func createNamespaceListCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "Lists all the CoreOps namespaces",
		Long:  "Use this list all the CoreOps namespaces inside this cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			namespace.List()
		},
	}

	return command
}

func createNamespaceCreateCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace",
		Long: "Use this command to create a new k8s namespace. This command is idempotent, " +
			"which means you can use it multiple times to ensure that the namespace is there, " +
			"as well as all other users.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.Users, cmd.Flags().Lookup(flags.Users))
			viper.BindPFlag(flags.ServiceAccounts, cmd.Flags().Lookup(flags.ServiceAccounts))
			viper.BindPFlag(flags.Name, cmd.Flags().Lookup(flags.Name))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.Create()
		},
	}

	command.PersistentFlags().StringP(flags.Name, "n", "", "Name of the namespace")
	command.MarkPersistentFlagRequired(flags.Name)

	addUsersPersistentFlag(command)
	command.MarkPersistentFlagRequired(flags.Users)

	addServiceAccountsPersistentFlag(command)

	return command
}

func createNamespaceDeleteCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing namespace",
		Long:  "Use this command to delete existing cops namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.Name, cmd.Flags().Lookup(flags.Name))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.Delete()
		},
	}

	command.PersistentFlags().StringP(flags.Name, "n", "", "Name of the namespace")
	command.MarkPersistentFlagRequired(flags.Name)

	return command
}

func createNamespaceUsersCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "users",
		Short: "Manage users of a namespace",
		Long:  "Use this command to manage the users in an existing k8s namespace.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	command.PersistentFlags().StringP(flags.Name, "n", "", "Name of the namespace")
	command.MarkPersistentFlagRequired(flags.Name)

	command.AddCommand(createNamespaceUsersAddCommand())
	command.AddCommand(createNamespaceUsersRemoveCommand())
	command.AddCommand(createNamespaceUsersListCommand())

	return command
}

func createNamespaceUsersAddCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "add",
		Short: "Adds users to the namespace",
		Long:  "Use this command to add users to an existing k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.Users, cmd.Flags().Lookup(flags.Users))
			viper.BindPFlag(flags.Name, cmd.Flags().Lookup(flags.Name))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.AddUsers()
		},
	}

	addUsersPersistentFlag(command)
	command.MarkPersistentFlagRequired(flags.Users)
	return command
}

func createNamespaceUsersRemoveCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "remove",
		Short: "Removes users from a namespace",
		Long:  "Use this command to remove users from an existing k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.Users, cmd.Flags().Lookup(flags.Users))
			viper.BindPFlag(flags.Name, cmd.Flags().Lookup(flags.Name))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.RemoveUsers()
		},
	}

	addUsersPersistentFlag(command)
	command.MarkPersistentFlagRequired(flags.Users)
	return command
}

func createNamespaceUsersListCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "List users of a namespace",
		Long:  "Use this command to list users of an existing k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.Name, cmd.Flags().Lookup(flags.Name))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.ListUsers()
		},
	}
	return command
}

func createNamespaceServiceAccountsCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "service-accounts",
		Short: "Manage service accounts of a namespace",
		Long:  "Use this command to manage the service accounts in an existing k8s namespace.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	command.PersistentFlags().StringP(flags.Name, "n", "", "Name of the namespace")
	command.MarkPersistentFlagRequired(flags.Name)

	command.AddCommand(createNamespaceServiceAccountsAddCommand())
	command.AddCommand(createNamespaceServiceAccountsRemoveCommand())
	command.AddCommand(createNamespaceServiceAccountsListCommand())

	return command
}

func createNamespaceServiceAccountsAddCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "add",
		Short: "Adds service accounts to the namespace",
		Long:  "Use this command to add service accounts to an existing k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.ServiceAccounts, cmd.Flags().Lookup(flags.ServiceAccounts))
			viper.BindPFlag(flags.Name, cmd.Flags().Lookup(flags.Name))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.AddServiceAccounts()
		},
	}

	addServiceAccountsPersistentFlag(command)
	command.MarkPersistentFlagRequired(flags.ServiceAccounts)
	return command
}

func createNamespaceServiceAccountsRemoveCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "remove",
		Short: "Removes users from a namespace",
		Long:  "Use this command to remove service-accounts from an existing k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.ServiceAccounts, cmd.Flags().Lookup(flags.ServiceAccounts))
			viper.BindPFlag(flags.Name, cmd.Flags().Lookup(flags.Name))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.RemoveServiceAccounts()
		},
	}

	addServiceAccountsPersistentFlag(command)
	command.MarkPersistentFlagRequired(flags.ServiceAccounts)
	return command
}

func createNamespaceServiceAccountsListCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "List service-accounts of a namespace",
		Long:  "Use this command to list service-accounts of an existing k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag(flags.Name, cmd.Flags().Lookup(flags.Name))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.ListServiceAccounts()
		},
	}
	return command
}

func addUsersPersistentFlag(command *cobra.Command) {
	command.PersistentFlags().StringP(flags.Users, "u", "", "The email-addresses of the admin "+
		"users of the namespace. Must be identical to Azure AD (case-sensitive). "+
		"You can specify multiple users separated by commas.")
}

func addServiceAccountsPersistentFlag(command *cobra.Command) {
	command.PersistentFlags().StringP(flags.ServiceAccounts, "s", "", "Optionally, you can specify service accounts "+
		"which will be granted idential access level like the users. Each service accounts has to be in "+
		"the format accountname.namespace, and multiple accounts can be specified separated by commas.")
}
