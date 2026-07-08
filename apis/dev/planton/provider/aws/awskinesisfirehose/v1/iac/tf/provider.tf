terraform {
  required_providers {
    aws = {
      # Feature-driven floor: iceberg_configuration.append_only landed in
      # 6.8.0 (the newest attribute this module renders; the Snowflake,
      # Iceberg, MSK-source, and secrets-manager blocks themselves all
      # predate the v6 line).
      source  = "hashicorp/aws"
      version = ">= 6.8.0"
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
