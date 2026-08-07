terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a ceiling: the full aws_cloudwatch_composite_alarm surface
      # (incl. actions_suppressor) is stable across the v6 line; the family
      # floor keeps sibling AWS modules on the same resolved major.
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
