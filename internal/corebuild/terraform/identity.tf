resource "azurerm_user_assigned_identity" "buildagentpool" {
  location            = data.azurerm_resource_group.buildagentpool.location
  name                = var.managed_identity_name
  resource_group_name = data.azurerm_resource_group.buildagentpool.name
}

resource "azurerm_role_assignment" "subscription" {
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.buildagentpool.principal_id
}

resource "azurerm_role_assignment" "vmss" {
  scope                = azurerm_linux_virtual_machine_scale_set.buildagentpool.id
  role_definition_name = "Virtual Machine Contributor"
  principal_id         = azurerm_user_assigned_identity.buildagentpool.principal_id
}

resource "azurerm_role_assignment" "buildagentpool" {
  for_each = var.role_assignments

  scope                = each.value.scope
  role_definition_name = each.value.role_definition_name
  principal_id         = azurerm_user_assigned_identity.buildagentpool.principal_id
}