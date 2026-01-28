# Auth0 Provider Configuration
# This file configures the Auth0 Terraform provider
# Documentation: https://registry.terraform.io/providers/auth0/auth0/latest/docs

# Variables for Auth0 authentication
variable "auth0_credential" {
  description = "Auth0 API authentication credentials"
  type = object({
    # Auth0 tenant domain (e.g., your-tenant.auth0.com or custom domain)
    domain = string

    # Auth0 Machine-to-Machine application Client ID
    # Create in Auth0 Dashboard: Applications -> Create Application -> Machine to Machine
    client_id = string

    # Auth0 Machine-to-Machine application Client Secret
    # Shown in the application settings after creation
    client_secret = string
  })
  sensitive = true
}

# Configure the Auth0 Provider
terraform {
  required_version = ">= 1.0"

  required_providers {
    auth0 = {
      source  = "auth0/auth0"
      version = "~> 1.0"
    }
  }
}

# Provider configuration
# The provider uses domain, client_id, and client_secret for API authentication
provider "auth0" {
  domain        = var.auth0_credential.domain
  client_id     = var.auth0_credential.client_id
  client_secret = var.auth0_credential.client_secret
}


