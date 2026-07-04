terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a ceiling: the cluster's managed_storage_configuration
      # and the capacity provider's managed_draining land on the v6 line;
      # `init` resolves the latest release at or above the floor.
      version = ">= 6.0.0"
    }
  }
}

provider "aws" {
  region = var.spec.region
}
