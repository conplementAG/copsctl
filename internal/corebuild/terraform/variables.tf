####################################
// Common
####################################
variable "resource_group_name" {}
variable "region" {}

####################################
// Identity
####################################
variable "managed_identity_name" {}
variable "role_assignments" {
  type = map(object({
    scope = string
    role_definition_name  = string
  }))
  default = {}
}

####################################
// Network
####################################
variable "vnet_name" {}
variable "build_agent_pool_public_ip_name" {}
variable "build_agent_pool_lb_name" {}

####################################
// Buildagent Pool
####################################
variable "build_agent_pool_name" {}
variable "build_agent_pool_node_sku" {}
variable "build_agent_pool_data_disk_enabled" {}
variable "build_agent_pool_data_disk_size_gb" {}
variable "build_agent_pool_data_disk_type" {}

####################################
// Azure Devops Buildagent
####################################
variable "azure_devops_project_name" {}
variable "azure_devops_service_connection_name" {}
variable "azure_devops_pool_name" {}
variable "azure_devops_desired_idle" {}
variable "azure_devops_max_capacity" {}
variable "azure_devops_ttl" {}