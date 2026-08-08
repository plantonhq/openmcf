terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.43"
    }
  }
}

provider "google" {
  credentials = var.provider_config.service_account_key != "" ? var.provider_config.service_account_key : null
}
