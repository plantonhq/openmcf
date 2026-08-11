terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# The features block stays at its defaults deliberately: destroying the
# protection deletes the backup data (the retain-on-destroy switches
# vm_backup_stop_protection_and_retain_data_on_destroy /
# vm_backup_suspend_protection_and_retain_data_on_destroy default OFF).
# Teams that need backup data to outlive the binding flip those
# switches in their own engine configuration -- documented in the
# GUIDE.
provider "azurerm" {
  features {}
}
