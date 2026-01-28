# Auth0ResourceServer Provider Configuration
# This file configures the Auth0 provider for Terraform

terraform {
  required_version = ">= 1.0"

  required_providers {
    auth0 = {
      source  = "auth0/auth0"
      version = ">= 1.0"
    }
  }
}

# Note: Provider configuration is expected to be set via environment variables:
# - AUTH0_DOMAIN
# - AUTH0_CLIENT_ID
# - AUTH0_CLIENT_SECRET
#
# Or passed via provider configuration in the root module.
