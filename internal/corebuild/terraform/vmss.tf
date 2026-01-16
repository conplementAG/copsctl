resource "random_password" "password" {
  length           = 64
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "azurerm_linux_virtual_machine_scale_set" "buildagentpool" {
  name                = var.build_agent_pool_name
  resource_group_name = data.azurerm_resource_group.buildagentpool.name
  location            = data.azurerm_resource_group.buildagentpool.location
  sku                 = var.build_agent_pool_node_sku == "" ? "Standard_B2s" : var.build_agent_pool_node_sku
  instances           = 0

  // either password or sshkey is required. push admin password in future to
  // devops keyvault when available
  disable_password_authentication = false
  admin_username                  = "corebuildadm"
  admin_password                  = random_password.password.result

  overprovision               = false
  upgrade_mode                = "Manual"
  single_placement_group      = false
  platform_fault_domain_count = 1
  custom_data = base64encode(templatefile("${path.module}/config/cloud-config.tpl", {
    use_data_disk = var.build_agent_pool_data_disk_enabled
  }))

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

  dynamic "data_disk" {
    for_each = var.build_agent_pool_data_disk_enabled ? [1] : []

    content {
      lun                  = 0
      disk_size_gb         = var.build_agent_pool_data_disk_size_gb
      caching              = "ReadWrite"
      storage_account_type = var.build_agent_pool_data_disk_type
    }
  }

  identity {
    type = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.buildagentpool.id]
  }

  network_interface {
    name    = "default"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.buildagentpool.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.buildagentpool.id]
    }
  }

  lifecycle {
    ignore_changes = [
      // ignore changes to instances as autoscaling is enabled
      instances,
      // ignored as azure devops agentpool will add tags
      tags,
    ]
  }
}