terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: the v6 line is where aws_redshift_cluster's
      # encrypted defaults to true, publicly_accessible defaults to
      # false, and the inline logging/snapshot_copy blocks gave way to
      # the standalone aws_redshift_logging and aws_redshift_snapshot_copy
      # resources this module uses.
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
