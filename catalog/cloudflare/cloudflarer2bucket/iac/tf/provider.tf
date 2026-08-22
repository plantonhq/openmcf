terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.23"
    }
  }
}

provider "cloudflare" {
  # Cloudflare provider configuration
  # API token should be provided via CLOUDFLARE_API_TOKEN environment variable
}

