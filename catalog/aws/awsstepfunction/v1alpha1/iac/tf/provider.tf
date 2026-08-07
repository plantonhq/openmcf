terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      # Family floor on the v6 line. Everything this module uses (publish +
      # versioning outputs, encryption_configuration, logging plan-time
      # validation) predates 6.0, so the floor is the family convention, not
      # a feature gate.
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
