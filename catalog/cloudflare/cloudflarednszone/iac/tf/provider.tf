terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.23"
    }
  }
}

provider "cloudflare" {
  # Cloudflare provider automatically uses CLOUDFLARE_API_TOKEN environment variable
}

