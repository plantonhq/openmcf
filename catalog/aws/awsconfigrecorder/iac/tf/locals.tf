locals {
  # Resource-identity tags match the Pulumi module key-for-key. None of
  # the four resources this module manages carries tags in the provider
  # (the recorder family is untagged AWS surface) -- the map exists for
  # convention and future taggable satellites.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsConfigRecorder"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
