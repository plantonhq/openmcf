terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

# Credentials are injected by the runtime as ARM_* environment variables --
# the empty block is what enables keyless (OIDC) auth. Keys are data-plane
# objects: beyond ARM permissions, the deploying credential needs key
# permissions on the target vault (the "Key Vault Administrator" or "Key
# Vault Crypto Officer" RBAC role, or key permissions in a legacy access
# policy). The default feature behavior purges a soft-deleted key on
# destroy so the name frees up immediately.
provider "azurerm" {
  features {}
}
