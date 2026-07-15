terraform {
  required_providers {
    aws = {
      # Family floor: >= 6.22.1 carries the full OpenZFS surface used by this
      # module — the INTELLIGENT_TIERING storage type with
      # read_cache_configuration (landed in 6.22.1), plus backup_id,
      # delete_options, and final_backup_tags.
      source  = "hashicorp/aws"
      version = ">= 6.22.1"
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
