terraform {
  required_providers {
    # Pessimistic minor float on the current major: every GCP module tracks the
    # same 7.x line so behavior is uniform across the catalog and upgrades to a
    # future major happen provider-wide in one decision, never per-kind.
    google = {
      source  = "hashicorp/google"
      version = "~> 7.0"
    }
  }
}

provider "google" {
}
