terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: Serverless v2 automatic pause
      # (serverlessv2_scaling_configuration.seconds_until_auto_pause) and
      # database_insights_mode are late-v5 additions -- the v6 floor keeps
      # the module on the modern major where the whole modeled surface
      # (including the iam-db-auth-error log type) is present.
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
