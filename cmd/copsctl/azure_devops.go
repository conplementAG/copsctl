package main

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"os"

	"github.com/conplementAG/copsctl/internal/azure_devops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createAzureDevopsCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "azure-devops",
		Short: "[DEPRECATED] Command group for administration of Azure DevOps accounts",
		Long: `
[DEPRECATED] Instead of setting up the Azure DevOps service connection, consider using copsctl connect just before deploying via Helm. 
This way you will not have a dependency on Azure DevOps and a prerequisite to setup the service connection before deploying.
		`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// bind these during pre-run to avoid viper overwriting itself when same PFlag used with multiple commands
			viper.BindPFlag(flags.EnvironmentTag, cmd.Flags().Lookup(flags.EnvironmentTag))
			viper.BindPFlag(flags.Organization, cmd.Flags().Lookup(flags.Organization))
			viper.BindPFlag(flags.Project, cmd.Flags().Lookup(flags.Project))
			viper.BindPFlag(flags.Namespace, cmd.Flags().Lookup(flags.Namespace))
			viper.BindPFlag(flags.Username, cmd.Flags().Lookup(flags.Username))
			viper.BindPFlag(flags.AccessToken, cmd.Flags().Lookup(flags.AccessToken))
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	command.PersistentFlags().StringP(flags.EnvironmentTag, "e", "", "Tag to use to make all environment names unique, e.g. your env identifier")
	command.MarkPersistentFlagRequired(flags.EnvironmentTag)

	command.PersistentFlags().StringP(flags.Organization, "o", "", "Name of the Azure DevOps organization")
	command.MarkPersistentFlagRequired(flags.Organization)

	command.PersistentFlags().StringP(flags.Project, "p", "", "The name of Azure DevOps project")
	command.MarkPersistentFlagRequired(flags.Project)

	command.PersistentFlags().StringP(flags.AccessToken, "t", "", "Your access token to access Azure DevOps (see. Azure DevOps->Profile->Security)")
	command.MarkPersistentFlagRequired(flags.AccessToken)

	command.PersistentFlags().StringP(flags.Username, "u", "", "Your username to access Azure DevOps (see. Azure DevOps->Profile->Security)")
	command.MarkPersistentFlagRequired(flags.Username)

	command.PersistentFlags().StringP(flags.Namespace, "n", "", "Namespace to which the endpoint can be scoped (RBAC rights)")

	command.AddCommand(createAzureDevopsCreateEndpointCommand())

	return command
}

func createAzureDevopsCreateEndpointCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create-endpoint",
		Short: "[DEPRECATED] Connect the Azure DevOps account to the current Kubernetes cluster.",
		Long: `
[DEPRECATED] Use this command to provision a service account which will be used to create a service endpoint in the Azure DevOps project.
This command is idempotent. Check azure-devops command info for deprecation explanation.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			azure_devops.NewOrchestrator().ConfigureEndpoint()
		},
	}
}
