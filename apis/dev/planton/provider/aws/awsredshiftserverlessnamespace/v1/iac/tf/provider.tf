terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: the v6 line carries the full
      # aws_redshiftserverless_namespace surface this module uses,
      # including the Secrets-Manager-managed admin password
      # (manage_admin_password + admin_password_secret_kms_key_id).
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
