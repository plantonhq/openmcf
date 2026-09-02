locals {
  # The AWS secret name is create-time immutable -- spec.secret_name when
  # set (hierarchical paths and service-required prefixes like
  # "ecr-pullthroughcache/..." that metadata.name cannot carry), else
  # metadata.name. Both engines share this resolution so a manifest
  # deploys identically on either.
  secret_name = var.spec.secret_name != "" ? var.spec.secret_name : var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSecretsManagerSecret"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Exactly one value arm may be set (CEL-enforced); either creates the
  # managed version. A shell secret (no value) is legal -- an application
  # or rotation function writes the first version then.
  create_version = var.spec.string_value != "" || var.spec.binary_value != ""

  # A policy is rendered only when the manifest declares one.
  create_policy = var.spec.policy != null

  # Rotation is configured only when the manifest declares the block.
  create_rotation = var.spec.rotation != null
}
