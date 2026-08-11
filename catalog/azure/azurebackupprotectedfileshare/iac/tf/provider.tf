terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# The features block stays at its defaults deliberately: destroying
# the protection stops it AND deletes the backup data (vault soft
# delete -- always on since Azure's secure-by-default policy -- may
# hold the deleted item's data 14 days). Teams that need backup data
# to outlive the binding manage protection state out of band --
# documented in the GUIDE.
provider "azurerm" {
  features {}
}
