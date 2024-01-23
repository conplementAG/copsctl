####################################
// Common
####################################
variable "resource_group_name" {}
variable "region" {}

####################################
// Identity
####################################
variable "managed_identity_name" {}

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