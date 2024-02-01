output "public_egress_ip" {
  value = azurerm_public_ip.buildagentpool.ip_address
}

output "managed_identity_name" {
  value = azurerm_user_assigned_identity.buildagentpool.name
}
