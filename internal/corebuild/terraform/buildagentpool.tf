data "azuredevops_project" "target" {
  name = var.azure_devops_project_name
}

resource "azuredevops_serviceendpoint_azurerm" "default" {
  project_id                             = data.azuredevops_project.target.id
  service_endpoint_name                  = var.azure_devops_service_connection_name
  description                            = "Managed by Terraform"
  service_endpoint_authentication_scheme = "WorkloadIdentityFederation"
  credentials {
    serviceprincipalid = azurerm_user_assigned_identity.buildagentpool.client_id
  }
  azurerm_spn_tenantid      = data.azurerm_subscription.current.tenant_id
  azurerm_subscription_id   = data.azurerm_subscription.current.subscription_id
  azurerm_subscription_name = data.azurerm_subscription.current.display_name

  depends_on = [
    azurerm_role_assignment.buildagentpool, azurerm_role_assignment.subscription, azurerm_role_assignment.vmss
  ]
}

resource "azurerm_federated_identity_credential" "default" {
  name                = "federated-credential-elastic-pool"
  resource_group_name = data.azurerm_resource_group.buildagentpool.name
  parent_id           = azurerm_user_assigned_identity.buildagentpool.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = azuredevops_serviceendpoint_azurerm.default.workload_identity_federation_issuer
  subject             = azuredevops_serviceendpoint_azurerm.default.workload_identity_federation_subject
}

resource "azuredevops_elastic_pool" "default" {
  name                   = var.azure_devops_pool_name
  service_endpoint_id    = azuredevops_serviceendpoint_azurerm.default.id
  service_endpoint_scope = data.azuredevops_project.target.project_id
  desired_idle           = var.azure_devops_desired_idle
  time_to_live_minutes   = var.azure_devops_ttl
  max_capacity           = var.azure_devops_max_capacity
  azure_resource_id      = azurerm_linux_virtual_machine_scale_set.buildagentpool.id
  depends_on             = [azurerm_federated_identity_credential.default]
}

resource "azuredevops_agent_queue" "default" {
  project_id    = data.azuredevops_project.target.id
  agent_pool_id = azuredevops_elastic_pool.default.id
}