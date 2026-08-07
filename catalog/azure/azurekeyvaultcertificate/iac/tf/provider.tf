terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# Credentials are injected by the runtime as ARM_* environment variables --
# the empty block is what enables keyless (OIDC) auth. Certificates are
# data-plane objects: beyond ARM permissions, the deploying credential
# needs certificate permissions on the target vault (the "Key Vault
# Administrator" or "Key Vault Certificates Officer" RBAC role, or
# certificate permissions in a legacy access policy). The default feature
# behavior purges a soft-deleted certificate on destroy so the name frees
# up immediately.
provider "azurerm" {
  features {}
}
