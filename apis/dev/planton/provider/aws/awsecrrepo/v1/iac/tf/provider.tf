terraform {
  required_providers {
    aws = {
      # Floor 6.12.0: the IMMUTABLE_WITH_EXCLUSION / MUTABLE_WITH_EXCLUSION
      # mutability modes and image_tag_mutability_exclusion_filter landed
      # across 6.10.0–6.12.0 (aws_ecr_repository_policy and
      # aws_ecr_lifecycle_policy are far older).
      source  = "hashicorp/aws"
      version = ">= 6.12.0"
    }
  }
}

provider "aws" {
  # Region and credentials are injected by the runtime as environment variables
  # (AWS_REGION + AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / AWS_SESSION_TOKEN), resolved
  # from the stack input's provider_config. For keyless (oidc)
  # connections the runtime performs the STS web-identity exchange and injects the resulting
  # short-lived credentials. Keep this block empty -- do not wire region or static keys here.
}
