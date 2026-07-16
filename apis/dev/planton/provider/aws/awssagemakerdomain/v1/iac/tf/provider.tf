terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor 6.33.0: the newest argument this module sends
      # (domain_settings.trusted_identity_propagation_settings) landed in
      # provider 6.33.0; every other block in the module predates it.
      version = ">= 6.33.0"
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
