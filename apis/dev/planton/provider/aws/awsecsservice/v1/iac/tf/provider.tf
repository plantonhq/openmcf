terraform {
  required_providers {
    aws = {
      # Floor, not a cap: native blue/green deployments
      # (deployment_configuration, load_balancer.advanced_configuration),
      # alarm-gated rollbacks, Service Connect access logs, and managed EBS
      # task volumes all landed across the provider's v6 line -- this floor
      # is the first release carrying the full surface this module drives.
      source  = "hashicorp/aws"
      version = ">= 6.50.0"
    }
  }

  required_version = ">= 1.0"
}

provider "aws" {
  # Region and credentials are injected by the runtime as environment variables
  # (AWS_REGION + AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / AWS_SESSION_TOKEN), resolved
  # from the stack input's provider_config. For keyless (oidc)
  # connections the runtime performs the STS web-identity exchange and injects the resulting
  # short-lived credentials. Keep this block empty -- do not wire region or static keys here.
}
