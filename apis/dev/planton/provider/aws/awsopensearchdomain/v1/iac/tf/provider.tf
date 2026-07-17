terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: aiml_options landed in 6.15.0, identity_center_options
      # in 6.20.0, and the aiml serverless_vector_acceleration block this module
      # emits arrived in 6.31.0 -- the newest argument in use.
      version = ">= 6.31.0"
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
