terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.0"
    }
  }
}

# Deleting an ML workspace is a SOFT delete: the ghost keeps holding the
# workspace NAME until purged, and the provider default
# (purge_soft_deleted_workspace_on_destroy = false) leaves the ghost
# behind — a delete-recreate cycle under the same name then fails on a
# collision with a resource that no longer appears anywhere. Purging on
# destroy makes destroy mean destroy; a soft-delete recovery window is
# not part of this module's contract (recreate from the manifest instead).
provider "azurerm" {
  features {
    machine_learning {
      purge_soft_deleted_workspace_on_destroy = true
    }
  }
}
