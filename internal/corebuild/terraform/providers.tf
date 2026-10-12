terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "5.9.0"
    }
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = "1.16.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.9.1"
    }
  }
  backend "azurerm" {
    container_name = "tfstate"
    key            = "terraform.tfstate"
  }

  required_version = ">= 1.14"
}

provider "azurerm" {
  // azurerm v5 no longer registers any resource providers by default (v4 registered the "legacy" set).
  // Register only what corebuild deploys, so bootstrapping into a fresh subscription keeps working.
  resource_providers_to_register = [
    "Microsoft.Compute",
    "Microsoft.ManagedIdentity",
    "Microsoft.Network",
  ]

  features {}
}

provider "azuredevops" {
  // required parameters org_service_url and personal_access_token are set via environment variables
}