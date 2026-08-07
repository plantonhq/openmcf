terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: DocumentDB Serverless
      # (serverless_v2_scaling_configuration, v6.10) and network_type
      # (v6.23) are v6-line additions -- the v6.23 floor keeps the module
      # on the modern major where the whole modeled surface is present.
      version = ">= 6.23.0"
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
