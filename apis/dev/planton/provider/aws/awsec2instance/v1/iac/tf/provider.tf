terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a ceiling: cpu_options.nested_virtualization (6.33) is
      # the newest argument this module uses; secondary_network_interface
      # landed in 6.32 and primary_network_interface in 6.10. `init`
      # resolves the latest release at or above the floor.
      version = ">= 6.33.0"
    }
  }
}

provider "aws" {
  region = var.spec.region
}
