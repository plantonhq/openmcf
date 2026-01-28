terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

# AWS provider configuration
# Credentials are passed via environment variables or provider_config
provider "aws" {
  # Region, access_key, and secret_key are configured via:
  # - Environment variables (AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
  # - Or passed via OpenMCF CLI
}
