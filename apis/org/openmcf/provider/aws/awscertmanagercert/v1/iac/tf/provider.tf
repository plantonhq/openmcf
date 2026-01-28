terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

# AWS provider configuration
# Region and other settings should be configured by the caller
provider "aws" {}


