package corebuild

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/conplementAG/copsctl/internal/adapters/azure"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/common"
	"github.com/conplementAG/copsctl/internal/common/file_processing"
	"github.com/conplementAG/copsctl/internal/corebuild/configuration"
	"github.com/conplementAG/copsctl/internal/corebuild/security"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/naming"
	"github.com/conplementag/cops-hq/v2/pkg/naming/resources"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/azure_login"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/terraform"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed terraform/*
var terraformDirectory embed.FS

type Orchestrator struct {
	hq                 hq.HQ
	executor           commands.Executor
	config             configuration.SourceConfig
	shortNamingService naming.Service
	longNamingService  naming.Service
	resourceGroupName  string
	autoApprove        bool
	cleanupActions     []func() error
}

func New(hq hq.HQ) (*Orchestrator, error) {
	configFile := viper.GetString(flags.ConfigFile)
	autoApprove := viper.GetBool(flags.AutoApproveFlag)

	sopsConfigFile := viper.GetString(flags.SopsConfigFile)
	if sopsConfigFile == "" {
		sopsConfigFile = filepath.Join(filepath.Dir(configFile), ".sops.yaml")
	}
	cryptographer, err := security.NewSopsCryptographer(sopsConfigFile)
	if err != nil {
		return nil, err
	}

	config, err := file_processing.LoadEncryptedFile[configuration.SourceConfig](configFile, cryptographer)
	if err != nil {
		return nil, err
	}
	shortNamingService, err := naming.New("cctl", config.Environment.Region, config.Environment.Name, "cb", "")
	if err != nil {
		return nil, err
	}
	longNamingService, err := naming.New("cctl", config.Environment.Region, config.Environment.Name, "corebuild", "")
	if err != nil {
		return nil, err
	}

	resourceGroupName, err := longNamingService.GenerateResourceName(resources.ResourceGroup, "build")
	if err != nil {
		return nil, err
	}

	return &Orchestrator{
		hq:                 hq,
		executor:           hq.GetExecutor(),
		config:             *config,
		shortNamingService: *shortNamingService,
		longNamingService:  *longNamingService,
		resourceGroupName:  resourceGroupName,
		autoApprove:        autoApprove,
	}, nil
}

func (o *Orchestrator) CreateInfrastructure() {
	logrus.Info("================== 🛀 🛀 🛀 Creating build agent pool 🛀 🛀 🛀  ====================")

	err := o.login()
	common.FatalOnError(err)

	tf, err := o.initializeTerraform()
	common.FatalOnError(err)

	err = tf.DeployFlow(false, false, o.autoApprove)
	common.FatalOnError(err)

	publicEgressIp, err := tf.Output(common.ToPtr("public_egress_ip"))
	common.FatalOnError(err)

	managedIdentityName, err := tf.Output(common.ToPtr("managed_identity_name"))
	common.FatalOnError(err)

	err = o.runCleanupActions()
	common.FatalOnError(err)

	logrus.Info("================== Build agent pool created  ====================")
	logrus.Infof("Make sure you add public egress ip %s to all resources firewall access lists build agent needs access", publicEgressIp)
	logrus.Infof("Make sure you add build agent managed identity %s to all resources permissions needed", managedIdentityName)
	logrus.Infof("Check configuration here: %s", fmt.Sprintf("https://dev.azure.com/%s/%s/_settings/agentqueues", o.config.AzureDevops.OrganisationName, o.config.AzureDevops.ProjectName))
	logrus.Info("=================================================================")
}

func (o *Orchestrator) DestroyInfrastructure() {
	logrus.Info("==================  Destroying build agent pool  ====================")

	err := o.login()
	common.FatalOnError(err)

	tf, err := o.initializeTerraform()
	common.FatalOnError(err)

	err = tf.DestroyFlow(false, false, o.autoApprove)
	common.FatalOnError(err)

	err = o.cleanup()
	common.FatalOnError(err)

	err = o.runCleanupActions()
	common.FatalOnError(err)

	logrus.Info("================== Build agent pool destroyed  ====================")
}

func (o *Orchestrator) login() error {
	azureLogin := azure_login.New(o.executor)
	err := azureLogin.Login()
	if err != nil {
		return fmt.Errorf("Login failed. %w", err)
	}

	err = azureLogin.SetSubscription(o.config.Environment.SubscriptionID)
	if err != nil {
		return fmt.Errorf("Failed to set current subscription %s. %w", o.config.Environment.SubscriptionID, err)
	}

	return nil
}

func (o *Orchestrator) initializeTerraform() (terraform.Terraform, error) {

	terraformStorageAccountName, err := o.shortNamingService.GenerateResourceName(resources.StorageAccount, "tf")
	if err != nil {
		return nil, err
	}

	backendStorageSettings := terraform.DefaultBackendStorageSettings
	backendStorageSettings.AllowedIpAddresses = o.config.Security.AuthorizedIPRanges.Cidrs
	backendStorageSettings.ContainerCreateRetryCount = 20

	// terraform azure devops provider authentication
	err = os.Setenv("AZDO_ORG_SERVICE_URL", fmt.Sprintf("https://dev.azure.com/%s", o.config.AzureDevops.OrganisationName))
	err = os.Setenv("AZDO_PERSONAL_ACCESS_TOKEN", o.config.AzureDevops.PatSecret)
	if err != nil {
		return nil, err
	}

	tempDir, err := file_processing.CreateTempDirectory(terraformDirectory, "terraform")
	o.cleanupActions = append(o.cleanupActions, func() error { return file_processing.DeletePath(tempDir) })
	if err != nil {
		return nil, err
	}

	tf := terraform.New(o.executor, "core-build",
		o.config.Environment.SubscriptionID,
		o.config.Environment.TenantID,
		o.config.Environment.Region,
		o.resourceGroupName,
		terraformStorageAccountName,
		tempDir,
		backendStorageSettings,
		terraform.DefaultDeploymentSettings)
	err = tf.Init()
	if err != nil {
		return nil, err
	}

	managedIdentityName, err := o.longNamingService.GenerateResourceName(resources.UserAssignedIdentity, "")
	if err != nil {
		return nil, err
	}
	vnetName, err := o.longNamingService.GenerateResourceName(resources.VirtualNetwork, "")
	if err != nil {
		return nil, err
	}
	publicIpName, err := o.longNamingService.GenerateResourceName(resources.PublicIp, "")
	if err != nil {
		return nil, err
	}
	loadBalancerName, err := o.longNamingService.GenerateResourceName(resources.LoadBalancer, "")
	if err != nil {
		return nil, err
	}
	vmss, err := o.longNamingService.GenerateResourceName(resources.VirtualMachineScaleSetLinux, "")
	if err != nil {
		return nil, err
	}

	vars := make(map[string]interface{})

	vars["resource_group_name"] = o.resourceGroupName
	vars["region"] = o.config.Environment.Region

	vars["managed_identity_name"] = managedIdentityName
	vars["vnet_name"] = vnetName

	roleAssignments := make(map[string]roleAssignment)
	for _, assignment := range o.config.Security.RoleAssignments {
		roleAssignments[fmt.Sprintf("%s-%s", assignment.Scope, assignment.RoleDefinitionName)] = roleAssignment{
			Scope:              assignment.Scope,
			RoleDefinitionName: assignment.RoleDefinitionName,
		}
	}
	vars["role_assignments"] = roleAssignments

	vars["build_agent_pool_public_ip_name"] = publicIpName
	vars["build_agent_pool_lb_name"] = loadBalancerName
	vars["build_agent_pool_name"] = vmss
	vars["build_agent_pool_node_sku"] = o.config.Environment.NodeSku

	vars["build_agent_pool_data_disk_enabled"] = o.config.Environment.DataDisk.Enabled
	vars["build_agent_pool_data_disk_size_gb"] = o.config.Environment.DataDisk.SizeGb
	vars["build_agent_pool_data_disk_type"] = o.config.Environment.DataDisk.Type

	vars["azure_devops_project_name"] = o.config.AzureDevops.ProjectName
	vars["azure_devops_service_connection_name"] = fmt.Sprintf("%s-federated-serviceconnection", o.config.Environment.Name)
	vars["azure_devops_pool_name"] = o.config.Environment.Name
	vars["azure_devops_desired_idle"] = o.config.AzureDevops.PoolSettings.DesiredIdle
	vars["azure_devops_max_capacity"] = o.config.AzureDevops.PoolSettings.MaxCapacity
	vars["azure_devops_ttl"] = o.config.AzureDevops.PoolSettings.TimeToLiveMinutes

	err = tf.SetVariables(vars)

	return tf, nil
}

func (o *Orchestrator) cleanup() error {
	azureAdapter, err := azure.New(o.config.Environment.SubscriptionID)
	if err != nil {
		return err
	}

	return azureAdapter.RemoveResourceGroup(o.resourceGroupName)
}

func (o *Orchestrator) runCleanupActions() error {
	var errs []string

	for _, f := range o.cleanupActions {
		if err := f(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New("combined error: " + fmt.Sprint(errs))
	}

	return nil
}

type roleAssignment struct {
	Scope              string `mapstructure:"scope" json:"scope" yaml:"scope"`
	RoleDefinitionName string `mapstructure:"role_definition_name" json:"role_definition_name" yaml:"role_definition_name"`
}
