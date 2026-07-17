terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      # Family floor: >= 6.26.0. The gateway's encryption_support argument
      # landed in 6.25.0 with a crash regression that 6.26.0 fixed, so the
      # floor deliberately clears the broken release. The whole Transit
      # Gateway family (gateway, VPC attachment, route table) shares this
      # floor so composed deployments resolve one provider build.
      source  = "hashicorp/aws"
      version = ">= 6.26.0"
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
