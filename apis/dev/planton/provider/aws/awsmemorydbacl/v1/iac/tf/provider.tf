terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: the MemoryDB family floor is set by the cluster's
      # ip_discovery/network_type arguments (v6.34.0); the ACL resource
      # itself is stable across the v6 line, and one family floor keeps the
      # engines' resolved provider versions aligned.
      version = ">= 6.34.0"
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
