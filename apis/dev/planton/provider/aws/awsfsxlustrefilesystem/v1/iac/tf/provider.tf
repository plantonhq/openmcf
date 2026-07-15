terraform {
  required_providers {
    aws = {
      # Family floor: >= 6.8.0 carries the full Lustre surface used by this
      # module — INTELLIGENT_TIERING with throughput_capacity and
      # data_read_cache_configuration (including the read-cache size
      # validation fixed in 6.8.0), efa_enabled, root_squash_configuration,
      # final_backup_tags, and the legacy S3 link arm.
      source  = "hashicorp/aws"
      version = ">= 6.8.0"
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
