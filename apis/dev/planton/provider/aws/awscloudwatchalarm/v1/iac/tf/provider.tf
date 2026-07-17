terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a ceiling. PromQL alarms (evaluation_criteria +
      # evaluation_interval) landed in v6.42.0, but that release shipped a
      # plan-time regression (spurious "One of 'metric_name', 'metric_query',
      # or 'evaluation_criteria' must be set" errors) fixed in v6.43.0 — so
      # the floor deliberately skips the broken release.
      version = ">= 6.43.0"
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
