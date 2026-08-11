terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# Credentials are injected by the runtime as ARM_* environment variables --
# the empty block is what enables keyless (OIDC) auth. Secrets are
# data-plane objects: beyond ARM permissions, the deploying credential
# needs secret permissions on the target vault (the "Key Vault
# Administrator" or "Key Vault Secrets Officer" RBAC role, or secret
# permissions in a legacy access policy). The default feature behavior
# purges a soft-deleted secret on destroy so the name frees up
# immediately (skipped automatically when the vault has purge
# protection), and recovers a soft-deleted secret on a name collision
# at create.
provider "azurerm" {
  features {}
}
