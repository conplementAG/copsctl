package main

import (
	"github.com/conplementAG/copsctl/internal/azure_devops"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementag/cops-hq/v2/pkg/cli"
	"github.com/conplementag/cops-hq/v2/pkg/hq"
)

func createAzureDevopsCommand(hq hq.HQ) {
	azureDevOpsCmdGroup := hq.GetCli().AddBaseCommand("azure-devops",
		"[DEPRECATED] Command group for administration of Azure DevOps accounts",
		"[DEPRECATED] Instead of setting up the Azure DevOps service connection, consider using copsctl connect just before deploying via Helm. "+
			"This way you will not have a dependency on Azure DevOps and a prerequisite to setup the service connection before deploying.",
		nil)

	azureDevOpsCmdGroup.AddPersistentParameterString(flags.EnvironmentTag, "", true, "e",
		"Tag to use to make all environment names unique, e.g. your env identifier")
	azureDevOpsCmdGroup.AddPersistentParameterString(flags.Organization, "", true, "o",
		"Name of the Azure DevOps organization")
	azureDevOpsCmdGroup.AddPersistentParameterString(flags.Project, "", true, "p",
		"The name of Azure DevOps project")
	azureDevOpsCmdGroup.AddPersistentParameterString(flags.AccessToken, "", true, "t",
		"Your access token to access Azure DevOps (see. Azure DevOps->Profile->Security)")
	azureDevOpsCmdGroup.AddPersistentParameterString(flags.Username, "", true, "u",
		"Your username to access Azure DevOps (see. Azure DevOps->Profile->Security)")
	azureDevOpsCmdGroup.AddPersistentParameterString(flags.Namespace, "", false, "n",
		"Namespace to which the endpoint can be scoped (RBAC rights)")

	createAzureDevopsCreateEndpointCommand(hq, azureDevOpsCmdGroup)
}

func createAzureDevopsCreateEndpointCommand(hq hq.HQ, cmd cli.Command) {

	cmd.AddCommand("create-endpoint", "[DEPRECATED] Connect the Azure DevOps account to the current Kubernetes cluster.",
		"[DEPRECATED] Use this command to provision a service account which will be used to create a service endpoint in the Azure DevOps project. "+
			"This command is idempotent. Check azure-devops command info for deprecation explanation.", func() {
			azure_devops.New(hq).ConfigureEndpoint()
		})
}
