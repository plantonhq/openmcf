terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      # Floor 6.0.0: the whole V2 surface this module renders (pipeline_type,
      # execution_mode, triggers with file-path filters, per-action
      # timeout_in_minutes, and stage conditions with rules) landed on the
      # 5.x line (<= 5.93.0), so any 6.x release carries it.
      source  = "hashicorp/aws"
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
