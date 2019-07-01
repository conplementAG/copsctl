package azuredevops

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/conplementAG/copsctl/pkg/adapters/azuredevops"
	"github.com/conplementAG/copsctl/pkg/adapters/kubernetes"
	"github.com/conplementAG/copsctl/pkg/common/fileprocessing"
	"github.com/conplementAG/copsctl/pkg/common/logging"

	"github.com/spf13/viper"
)

type AzureDevopsOrchestrator struct {
	Organization       string
	Project            string
	Namespace          string
	globalScope        bool
	serviceAccountName string
	roleName           string
	endpointName       string
}

func NewOrchestrator() *AzureDevopsOrchestrator {
	environmentTag := viper.GetString("environment-tag")
	organization := viper.GetString("organization")
	project := viper.GetString("project")
	namespace := viper.GetString("namespace")

	if organization == "" {
		panic("Organization configuration is mandatory!")
	}

	if project == "" {
		panic("Project configuratifon is mandatory!")
	}

	isGlobalScope := namespace == ""

	if isGlobalScope {
		namespace = "kube-system"
	}

	return &AzureDevopsOrchestrator{
		Organization: organization,
		Project:      project,
		Namespace:    namespace,

		globalScope:        isGlobalScope,
		serviceAccountName: strings.ToLower(organization) + "-" + strings.ToLower(project) + "-azuredevops-account",
		roleName:           strings.ToLower(organization) + "-" + strings.ToLower(project) + "-" + namespace + "-azuredevops-role",
		endpointName:       environmentTag + "-" + namespace,
	}
}

func (orchestrator *AzureDevopsOrchestrator) ConfigureEndpoint() {
	logging.LogProvisioningStepStarted("Connecting the current k8s cluster with an Azure DevOps account...")

	if orchestrator.globalScope {
		log.Println("RBAC will be without limitation, since no 'namespace' was specified, and the RBAC resources will be in kube-system")
	} else {
		log.Println("RBAC will be scoped to namespace " + orchestrator.Namespace)
	}

	logging.LogProvisioningStepStarted("Creating the RBAC resources")

	outputPath := orchestrator.prepareRbacFiles()

	kubernetes.Apply(outputPath)
	fileprocessing.DeletePath(outputPath)

	logging.LogProvisioningStepStarted("Setting up the Azure DevOps connection...")

	// first, get the token, the certificate of the created service account and the master plane FQDN
	serviceAccount := kubernetes.GetServiceAccount(orchestrator.Namespace, orchestrator.serviceAccountName)

	if len(serviceAccount.Secrets) != 1 {
		panic("Expected the service account to contain exactly one secret (where the token and cert are located)")
	}

	serviceAccountSecret := kubernetes.GetServiceAccountSecret(orchestrator.Namespace, serviceAccount.Secrets[0].Name)

	masterPlaneFqdn := kubernetes.GetCurrentMasterPlaneFqdn()

	// now we can create the endpoint (aka. service connection / service endpoint)
	azuredevops.CreateServiceEndpoint(orchestrator.endpointName, orchestrator.Organization, orchestrator.Project, masterPlaneFqdn, serviceAccountSecret.Data.Token, serviceAccountSecret.Data.CaCrt)
}

func (orchestrator *AzureDevopsOrchestrator) prepareRbacFiles() string {
	if orchestrator.globalScope {
		return orchestrator.prepareGlobalRbacFiles()
	} else {
		return orchestrator.prepareScopedRbacFiles()
	}
}

func (orchestrator *AzureDevopsOrchestrator) prepareGlobalRbacFiles() string {
	return fileprocessing.InterpolateStaticFiles(
		"pkg/azuredevops/global",
		map[string]string{
			"{{NAMESPACE}}":       "kube-system",
			"{{ROLE_NAME}}":       orchestrator.roleName,
			"{{BINDING_NAME}}":    orchestrator.roleName + "-binding",
			"{{SERVICE_ACCOUNT}}": orchestrator.serviceAccountName,
		})
}

func (orchestrator *AzureDevopsOrchestrator) prepareScopedRbacFiles() string {
	return fileprocessing.InterpolateFiles(
		filepath.Join("..", "pkg", "azuredevops", "scoped"),
		map[string]string{
			"{{NAMESPACE}}":       orchestrator.Namespace,
			"{{ROLE_NAME}}":       orchestrator.roleName,
			"{{BINDING_NAME}}":    orchestrator.roleName + "-binding",
			"{{SERVICE_ACCOUNT}}": orchestrator.serviceAccountName,
		})
}
