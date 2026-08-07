terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      # Floor: >= 6.0.0. The newest arguments this module sends — the
      # endpoint configuration's attachment_arn (cross-account attachments)
      # and the accelerator's dual-stack surface — landed on the v5 line
      # (v5.47-v5.48), so the v6 major floor carries the full modeled
      # surface with current provider behavior.
      source  = "hashicorp/aws"
      version = ">= 6.0.0"
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
