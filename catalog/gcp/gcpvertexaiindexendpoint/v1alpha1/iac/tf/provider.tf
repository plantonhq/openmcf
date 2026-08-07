terraform {
  required_providers {
    # Pessimistic minor float on the current major: every GCP module tracks the
    # same 6.x line so behavior is uniform across the catalog and upgrades to a
    # future major happen provider-wide in one decision, never per-kind. Every
    # field this module sends is GA on the released 6.x line (verified against
    # the installed 6.50.0 provider schema), so no google-beta dependency
    # exists to drift.
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
}

provider "google" {
  credentials = var.provider_config.service_account_key != "" ? var.provider_config.service_account_key : null
}
