terraform {
  required_providers {
    # Pessimistic minor float on the current major: every GCP module tracks the
    # same 6.x line so behavior is uniform across the catalog and upgrades to a
    # future major happen provider-wide in one decision, never per-kind.
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    # The subnetwork resource runs on google-beta: allow_subnet_cidr_routes_overlap
    # is beta-only on the released 6.x line, and selecting the provider
    # per-resource keeps the full surface available without a per-field
    # retrofit later. All other fields are GA and identical in beta.
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
