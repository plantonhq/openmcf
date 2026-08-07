terraform {
  required_providers {
    aws = {
      # Family floor: the complete volume surface used by this module
      # predates the v6 line (the last addition — final_backup_tags — landed
      # in 5.59.0), so the floor is the v6 major itself, keeping the FSx
      # family on one provider line.
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
