terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.8.0"
    }
  }

  required_version = ">= 1.6"
}

provider "azurerm" {
  features {}
}