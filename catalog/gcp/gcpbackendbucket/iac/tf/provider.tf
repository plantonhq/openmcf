terraform {
  required_providers {
    # Pessimistic minor float on the current major: every GCP module tracks the
    # same 6.x line so behavior is uniform across the catalog and upgrades to a
    # future major happen provider-wide in one decision, never per-kind.
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    # Declared alongside google so preview-stage resources or fields can opt in
    # with an explicit `provider = google-beta` (none needed by this module
    # today; the backend bucket surface is GA and identical in beta).
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 6.0"
    }
  }
}

provider "google" {
}

provider "google-beta" {
}
