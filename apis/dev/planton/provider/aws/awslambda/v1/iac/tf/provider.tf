terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a pin: 6.0 carries the full surface this module
      # renders (logging_config, snap_start, ipv6_allowed_for_dual_stack,
      # source_kms_key_arn, the recursion and runtime-management
      # resources all predate the v6 line).
      version = ">= 6.0.0"
    }
  }
}

provider "aws" {
  # Region and credentials are injected by the runtime as environment variables
  # (AWS_REGION + AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / AWS_SESSION_TOKEN), resolved
  # from the stack input's provider_config. For keyless (oidc / cross_account_trust)
  # connections the runtime performs the STS web-identity exchange and injects the resulting
  # short-lived credentials. Keep this block empty -- do not wire region or static keys here.
}
