terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# Azure services leave nested resources this module never created
# (Application Insights writes a "Smart Detection" action group into
# the group; similar furniture shows up for other services). The
# provider default prevent_deletion_if_contains_resources = true is a
# Terraform-state seatbelt that then refuses to delete the group,
# stranding an otherwise-empty RG. Destroy of a Planton resource group
# means ARM-delete the group -- the same contract as `az group delete`.
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}
