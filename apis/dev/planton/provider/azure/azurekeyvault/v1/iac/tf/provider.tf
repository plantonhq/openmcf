terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

# Credentials are injected by the runtime as ARM_* environment variables --
# the empty block is what enables keyless (OIDC) auth. The default feature
# behavior purges a soft-deleted vault on destroy (unless purge protection
# is on, which turns destroy into a scheduled deletion) and auto-recovers a
# soft-deleted vault on a name-colliding create.
provider "azurerm" {
  features {}
}
