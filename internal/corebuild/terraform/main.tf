data "azurerm_subscription" "current" {
}

data "azurerm_resource_group" "buildagentpool" {
  name     = var.resource_group_name
}
