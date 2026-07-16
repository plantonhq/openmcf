terraform {
  required_providers {
    # Pessimistic minor float on the current major: every GCP module tracks the
    # same 6.x line so behavior is uniform across the catalog and upgrades to a
    # future major happen provider-wide in one decision, never per-kind.
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    # Declared for catalog-wide consistency; every field this module sends is
    # GA on the released 6.x line, so the resource selects the plain google
    # provider. (The only beta-only region-NEG surface — serverless_deployment
    # for API Gateway — is deliberately not modeled.)
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
