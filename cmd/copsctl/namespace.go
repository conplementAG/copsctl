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
		Long: `
Use this command to administer k8s namespaces.
        `,
		Run: func(cmd *cobra.Command, args []string) {
			// since "namespace" is not really a command, but rather a group of commands, we
			// show the help for the command group instead
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	command.AddCommand(createNamespaceCreateCommand())

	command.PersistentFlags().StringP("name", "n", "", "Name of the namespace")
	command.MarkPersistentFlagRequired("name")

	return command
}

func createNamespaceCreateCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace",
		Long: `
Use this command to create a new k8s namespace.
        `,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("users", cmd.Flags().Lookup("users"))
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace.Create()
		},
	}

	command.PersistentFlags().StringP("users", "u", "", "The email-addresses of the admin users of the namespace. Must be identical to Azure AD (case-sensitive). You can specify multiple users separated by comma.")
	command.MarkPersistentFlagRequired("users")
	return command
}
