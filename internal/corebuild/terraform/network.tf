resource "azurerm_virtual_network" "buildagentpool" {
  name                = var.vnet_name
  resource_group_name = data.azurerm_resource_group.buildagentpool.name
  location            = data.azurerm_resource_group.buildagentpool.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "buildagentpool" {
  name                 = "internal"
  resource_group_name  = data.azurerm_resource_group.buildagentpool.name
  virtual_network_name = azurerm_virtual_network.buildagentpool.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "buildagentpool" {
  name                = var.build_agent_pool_public_ip_name
  location            = data.azurerm_resource_group.buildagentpool.location
  resource_group_name = data.azurerm_resource_group.buildagentpool.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "buildagentpool" {
  name                = var.build_agent_pool_lb_name
  location            = data.azurerm_resource_group.buildagentpool.location
  resource_group_name = data.azurerm_resource_group.buildagentpool.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "public-ip-config"
    public_ip_address_id = azurerm_public_ip.buildagentpool.id
  }
}

resource "azurerm_lb_backend_address_pool" "buildagentpool" {
  loadbalancer_id = azurerm_lb.buildagentpool.id
  name            = "backend-pool"
}

resource "azurerm_lb_outbound_rule" "buildagentpool_outbound_default" {
  name                    = "outbound-rule"
  loadbalancer_id         = azurerm_lb.buildagentpool.id
  protocol                = "All"
  backend_address_pool_id = azurerm_lb_backend_address_pool.buildagentpool.id
  frontend_ip_configuration {
    name = "public-ip-config"
  }
}