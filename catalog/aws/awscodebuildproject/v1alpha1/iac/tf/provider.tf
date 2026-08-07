terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      # Floor 6.16.0: the project's auto_retry_limit lands there (the
      # webhook's pull_request_build_policy 6.13.0, the environment's
      # docker_server 6.2.0, and everything else are below the floor).
      source  = "hashicorp/aws"
      version = ">= 6.16.0"
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
