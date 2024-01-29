output "public_egress_ip" {
  value = azurerm_public_ip.buildagentpool.ip_address
}
