####################################
// Common
####################################

resource "azurerm_resource_group" "buildagentpool" {
  name     = var.resource_group_name
  location = var.region
}

####################################
// Identity
####################################

resource "azurerm_user_assigned_identity" "buildagentpool" {
  location            = azurerm_resource_group.buildagentpool.location
  name                = var.managed_identity_name
  resource_group_name = azurerm_resource_group.buildagentpool.name
}

####################################
// Network
####################################

resource "azurerm_virtual_network" "buildagentpool" {
  name                = var.vnet_name
  resource_group_name = azurerm_resource_group.buildagentpool.name
  location            = azurerm_resource_group.buildagentpool.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "buildagentpool" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.buildagentpool.name
  virtual_network_name = azurerm_virtual_network.buildagentpool.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "buildagentpool" {
  name                = var.build_agent_pool_public_ip_name
  location            = azurerm_resource_group.buildagentpool.location
  resource_group_name = azurerm_resource_group.buildagentpool.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "buildagentpool" {
  name                = var.build_agent_pool_lb_name
  location            = azurerm_resource_group.buildagentpool.location
  resource_group_name = azurerm_resource_group.buildagentpool.name
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

####################################
// Buildagent Pool
####################################
resource "random_password" "password" {
  length           = 64
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "azurerm_linux_virtual_machine_scale_set" "buildagentpool" {
  name                = var.build_agent_pool_name
  resource_group_name = azurerm_resource_group.buildagentpool.name
  location            = azurerm_resource_group.buildagentpool.location
  sku                 = "Standard_B2s"
  instances           = 1

  // either password or sshkey is required. push admin password in future to
  // devops keyvault when available
  disable_password_authentication = false
  admin_username                  = "corebuildadm"
  admin_password                  = random_password.password.result

  overprovision               = false
  upgrade_mode                = "Manual"
  single_placement_group      = false
  platform_fault_domain_count = 1
  custom_data                 = filebase64("${path.module}/config/cloud-config.txt")

  # https://learn.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-automatic-upgrade#supported-os-images
  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "StandardSSD_LRS"
    caching              = "ReadWrite"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.buildagentpool.id]
  }

  network_interface {
    name    = "default"
    primary = true

    ip_configuration {
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.buildagentpool.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.buildagentpool.id]
    }
  }
}