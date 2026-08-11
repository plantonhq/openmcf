terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# Credentials are injected by the runtime as ARM_* environment variables --
# the empty block is what enables keyless (OIDC) auth. Destroying a backup
# instance stops protection and deletes its backup data; when the vault's
# soft delete is on, the data lingers as a soft-deleted item for 14 days
# and HOLDS the vault's own deletion for that window.
provider "azurerm" {
  features {}
}
