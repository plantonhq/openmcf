terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      # Floor 6.34.0: the workgroup's query_results_s3_access_grants_configuration
      # lands there (identity_center 6.7.0, managed_query_results 6.25.0, and
      # monitoring_configuration 6.28.0 are all below the floor).
      source  = "hashicorp/aws"
      version = ">= 6.34.0"
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
