package main

import (
	"os"

	"github.com/conplementAG/copsctl/internal/azuredevops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createAzureDevopsCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "azure-devops",
		Short: "Command group for administration of Azure DevOps accounts",
		Long: `
Use this command to administer Azure DevOps accounts (aka VSTS / hosted TFS).
		`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// bind these during pre-run to avoid viper overwriting itself when same PFlag used with multiple commands
			viper.BindPFlag("environment-tag", cmd.Flags().Lookup("environment-tag"))
			viper.BindPFlag("region", cmd.Flags().Lookup("region"))
			viper.BindPFlag("organization", cmd.Flags().Lookup("organization"))
			viper.BindPFlag("project", cmd.Flags().Lookup("project"))
			viper.BindPFlag("namespace", cmd.Flags().Lookup("namespace"))
			viper.BindPFlag("username", cmd.Flags().Lookup("username"))
			viper.BindPFlag("accesstoken", cmd.Flags().Lookup("accesstoken"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	command.PersistentFlags().StringP("environment-tag", "e", "", "Tag to use to make all environment names unique, e.g. your env identifier")
	command.MarkPersistentFlagRequired("environment-tag")

	command.PersistentFlags().StringP("region", "r", "northeurope", "Region where the resources should be deployed, e.g. northeurope")

	command.PersistentFlags().StringP("organization", "o", "", "Name of the VSTS organization")
	command.MarkPersistentFlagRequired("organization")

	command.PersistentFlags().StringP("project", "p", "", "The name of VSTS project")
	command.MarkPersistentFlagRequired("project")

	command.PersistentFlags().StringP("accesstoken", "t", "", "Your access token to access Azure DevOps (see. Azure DevOps->Profile->Security)")
	command.MarkPersistentFlagRequired("accesstoken")

	command.PersistentFlags().StringP("username", "u", "", "Your username to access Azure DevOps (see. Azure DevOps->Profile->Security)")
	command.MarkPersistentFlagRequired("username")

	command.PersistentFlags().StringP("namespace", "n", "", "Namespace to which the endpoint can be scoped (RBAC rights)")

	command.AddCommand(createAzureDevopsCreateEndpointCommand())

	return command
}

func createAzureDevopsCreateEndpointCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create-endpoint",
		Short: "Connect the Azure DevOps account to the current Kubernetes cluster.",
		Long: `
Use this command to provision a service account which will be used to create a service endpoint in the Azure DevOps project.
This command is idempotent.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			azuredevops.NewOrchestrator().ConfigureEndpoint()
		},
	}
}
