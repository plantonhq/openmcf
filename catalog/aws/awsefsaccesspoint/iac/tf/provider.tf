terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a ceiling: the access point surface is stable across the v6
      # line; the family floor matches the AwsElasticFileSystem module so the
      # two EFS kinds always resolve the same provider major.
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
