terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# Credentials are injected by the runtime as ARM_* environment variables
# (service principal or keyless OIDC). Keep this block empty -- wiring
# static credentials here would break the keyless path.
provider "azurerm" {
  features {}
}
