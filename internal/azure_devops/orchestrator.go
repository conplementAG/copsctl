package azure_devops

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"strings"
	"time"

	"net/url"

	"github.com/conplementAG/copsctl/internal/adapters/azure_devops"
	"github.com/conplementAG/copsctl/internal/adapters/kubernetes"
	"github.com/conplementAG/copsctl/internal/common/file_processing"
	"github.com/conplementAG/copsctl/internal/common/logging"

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
	username           string
	accesstoken        string
}

func NewOrchestrator() *AzureDevopsOrchestrator {
	environmentTag := viper.GetString(flags.EnvironmentTag)
	organization := viper.GetString(flags.Organization)
	project := viper.GetString(flags.Project)
	namespace := viper.GetString(flags.Namespace)
	username := viper.GetString(flags.Username)
	accessToken := viper.GetString(flags.AccessToken)

	isGlobalScope := namespace == ""

	if isGlobalScope {
		namespace = "kube-system"
	}

	return &AzureDevopsOrchestrator{
		Organization: trim(organization),
		Project:      trim(project),
		Namespace:    trim(namespace),

		globalScope:        isGlobalScope,
		serviceAccountName: strings.ToLower(urlDecode(organization)) + "-" + strings.ToLower(urlDecode(project)) + "-azuredevops-account",
		roleName:           strings.ToLower(urlDecode(organization)) + "-" + strings.ToLower(urlDecode(project)) + "-" + trim(namespace) + "-azuredevops-role",
		endpointName:       trim(environmentTag) + "-" + trim(namespace),
		username:           trim(username),
		accesstoken:        trim(accessToken),
	}
}

func urlDecode(source string) string {
	urlDecoded, _ := url.QueryUnescape(source)
	return trim(urlDecoded)
}

func trim(source string) string {
	return strings.Replace(source, " ", "", -1)
}

func (orchestrator *AzureDevopsOrchestrator) ConfigureEndpoint() {
	logging.Info("Connecting the current k8s cluster with an Azure DevOps account...")

	if orchestrator.globalScope {
		logging.Info("RBAC will be without limitation, since no 'namespace' was specified, and the RBAC resources will be in kube-system")
	} else {
		logging.Info("RBAC will be scoped to namespace " + orchestrator.Namespace)
	}

	logging.Info("Creating the RBAC resources")

	outputPath := orchestrator.prepareRbacFiles()

	_, err := kubernetes.Apply(outputPath)

	if err != nil {
		panic("Apply failed: " + err.Error())
	}

	file_processing.DeletePath(outputPath)

	logging.Info("Setting up the Azure DevOps connection...")

	// sleep a bit to make sure the secret is created
	time.Sleep(3 * time.Second)

	// first, get the token, the certificate of the created service account and the master plane FQDN
	serviceAccount := kubernetes.GetServiceAccount(orchestrator.Namespace, orchestrator.serviceAccountName)

	if len(serviceAccount.Secrets) != 1 {
		panic("Expected the service account to contain exactly one secret (where the token and cert are located)")
	}

	serviceAccountSecret := kubernetes.GetServiceAccountSecret(orchestrator.Namespace, serviceAccount.Secrets[0].Name)

	masterPlaneFqdn, err := kubernetes.GetCurrentMasterPlaneFqdn()

	if err != nil {
		panic("Could not get the master plane fqdn " + err.Error())
	}

	// now we can create the endpoint (aka. service connection / service endpoint)
	azure_devops.CreateServiceEndpoint(
		orchestrator.username,
		orchestrator.accesstoken,
		orchestrator.endpointName,
		orchestrator.Organization,
		orchestrator.Project,
		masterPlaneFqdn,
		serviceAccountSecret.Data.Token,
		serviceAccountSecret.Data.CaCrt)
}

func (orchestrator *AzureDevopsOrchestrator) prepareRbacFiles() string {
	if orchestrator.globalScope {
		return orchestrator.prepareGlobalRbacFiles()
	} else {
		return orchestrator.prepareScopedRbacFiles()
	}
}

func (orchestrator *AzureDevopsOrchestrator) prepareGlobalRbacFiles() string {
	return file_processing.InterpolateStaticFiles(
		"internal/azure_devops/global",
		map[string]string{
			"{{NAMESPACE}}":       "kube-system",
			"{{BINDING_NAME}}":    orchestrator.roleName + "-binding",
			"{{SERVICE_ACCOUNT}}": orchestrator.serviceAccountName,
		})
}

func (orchestrator *AzureDevopsOrchestrator) prepareScopedRbacFiles() string {
	return file_processing.InterpolateStaticFiles(
		"internal/azure_devops/scoped",
		map[string]string{
			"{{NAMESPACE}}":       orchestrator.Namespace,
			"{{BINDING_NAME}}":    orchestrator.roleName + "-binding",
			"{{SERVICE_ACCOUNT}}": orchestrator.serviceAccountName,
		})
}
