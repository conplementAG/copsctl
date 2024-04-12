terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.98.0"
    }
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = "1.0.1"
    }
  }
  backend "azurerm" {
    container_name = "tfstate"
    key            = "terraform.tfstate"
  }

  required_version = ">= 1.7"
}

provider "azurerm" {
  features {}
}

provider "azuredevops" {
  // required parameters org_service_url and personal_access_token are set via environment variables
}