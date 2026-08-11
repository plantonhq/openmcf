terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# The features block stays at its defaults deliberately. The three
# recovery-service destroy switches it can carry
# (vm_backup_stop_protection_and_retain_data_on_destroy,
# vm_backup_suspend_protection_and_retain_data_on_destroy,
# purge_protected_items_from_vault_on_destroy) all default OFF: destroy
# deletes backup data with the protection, and a vault with remaining
# protected items refuses to delete -- the honest, least-surprising
# posture for a managed module.
provider "azurerm" {
  features {}
}
