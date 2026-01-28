# MongoDB Atlas Provider Configuration
# This file configures the MongoDB Atlas Terraform provider
# Documentation: https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs

# Variables for MongoDB Atlas authentication
variable "mongodbatlas_credential" {
  description = "MongoDB Atlas authentication credentials"
  type = object({
    # MongoDB Atlas Public API Key
    # Create API keys in Atlas UI: Project Settings -> Access Manager -> API Keys
    public_key = string

    # MongoDB Atlas Private API Key
    # This key is shown only once when created and should be stored securely
    private_key = string
  })
  sensitive = true
}

# Configure the MongoDB Atlas Provider
terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.14"
    }
  }
}

# Provider configuration
# The provider uses the public and private keys for API authentication
provider "mongodbatlas" {
  public_key  = var.mongodbatlas_credential.public_key
  private_key = var.mongodbatlas_credential.private_key
}

