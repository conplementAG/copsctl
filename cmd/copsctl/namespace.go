package main

import (
	"os"

	"github.com/conplementAG/copsctl/pkg/namespace"
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

	return command
}

func createNamespaceListCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "lists all the CoreOps namespaces",
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
		Long:  "Use this command to create a new k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("users", cmd.Flags().Lookup("users"))
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.Create()
		},
	}

	command.PersistentFlags().StringP("name", "n", "", "Name of the namespace")
	command.MarkPersistentFlagRequired("name")

	command.PersistentFlags().StringP("users", "u", "", "The email-addresses of the admin users of the namespace. Must be identical to Azure AD (case-sensitive). You can specify multiple users separated by comma.")
	command.MarkPersistentFlagRequired("users")
	return command
}

func createNamespaceDeleteCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "delete",
		Short: "Delete a new namespace",
		Long:  "Use this command to delete existing copsnamespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.Delete()
		},
	}

	command.PersistentFlags().StringP("name", "n", "", "Name of the namespace")
	command.MarkPersistentFlagRequired("name")

	return command
}

func createNamespaceUsersCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "users",
		Short: "Mange users of a namespace",
		Long:  "Use this command to mange the users in an existing k8s namespace.",
		Run: func(cmd *cobra.Command, args []string) {
			// since "namespace" is not really a command, but rather a group of commands, we
			// show the help for the command group instead
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	command.PersistentFlags().StringP("name", "n", "", "Name of the namespace")
	command.MarkPersistentFlagRequired("name")

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
			viper.BindPFlag("users", cmd.Flags().Lookup("users"))
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.AddUsers()
		},
	}

	command.PersistentFlags().StringP("users", "u", "", "The email-addresses of the users that you want to add. Must be identical to Azure AD (case-sensitive). You can specify multiple users separated by comma.")
	command.MarkPersistentFlagRequired("users")
	return command
}

func createNamespaceUsersRemoveCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "remove",
		Short: "Removes users from a namespace",
		Long:  "Use this command to remove users from an existing k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("users", cmd.Flags().Lookup("users"))
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.RemoveUsers()
		},
	}

	command.PersistentFlags().StringP("users", "u", "", "The email-addresses of the users that you want to remove from the namespace. Must be identical to Azure AD (case-sensitive). You can specify multiple users separated by comma.")
	command.MarkPersistentFlagRequired("users")
	return command
}

func createNamespaceUsersListCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "List users of a namespace",
		Long:  "Use this command to list users of an existing k8s namespace.",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.ListUsers()
		},
	}
	return command
}
