terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      # Family floor: >= 6.29.0 carries the full Windows surface used by this
      # module — the self-managed AD domain_join_service_account_secret arm
      # with optional username/password (landed in 6.29.0), plus backup_id
      # and final_backup_tags.
      source  = "hashicorp/aws"
      version = ">= 6.29.0"
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
