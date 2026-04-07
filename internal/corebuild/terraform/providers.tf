terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.67.0"
    }
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = "1.15.0"
    }
  }
  backend "azurerm" {
    container_name = "tfstate"
    key            = "terraform.tfstate"
  }

  required_version = ">= 1.14"
}

provider "azurerm" {
  features {}
}

provider "azuredevops" {
  // required parameters org_service_url and personal_access_token are set via environment variables
}