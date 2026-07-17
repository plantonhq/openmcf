terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a pin: 6.37.0 carries the full surface this module
      # renders -- global_table_witness (6.22), multi-attribute GSI
      # key_schema (6.29), and the 6.37 fix for GSI removal under the
      # key_schema syntax deleting every index on the table.
      version = ">= 6.37.0"
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
