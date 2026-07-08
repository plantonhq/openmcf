terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      # The database surface (resource links, federation, Lake Formation
      # create-table grants) is stable across the v6 line; the floor keeps the
      # kind on the provider family the rest of the AWS catalog targets.
      source  = "hashicorp/aws"
      version = ">= 6.0.0"
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
