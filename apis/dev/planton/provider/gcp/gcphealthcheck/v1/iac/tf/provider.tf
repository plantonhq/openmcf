terraform {
  required_providers {
    # Pessimistic minor float on the current major: every GCP module tracks the
    # same 6.x line so behavior is uniform across the catalog and upgrades to a
    # future major happen provider-wide in one decision, never per-kind.
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    # The health check resources run on google-beta: the grpc_tls_health_check
    # block is beta-only on the released 6.x line, and selecting the provider
    # per-resource keeps the whole protocol surface available without a
    # per-field retrofit later. All other fields are GA and identical in beta.
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
