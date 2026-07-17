terraform {
  required_providers {
    aws = {
      # v6 family floor: the full SESv2 email-identity surface modeled here
      # (identity + mail-from + feedback + policy sub-resources) is stable
      # before the v6 line; the floor keeps the SES family on one provider
      # generation.
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
