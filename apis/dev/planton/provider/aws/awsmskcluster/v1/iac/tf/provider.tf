terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: connectivity_info.network_type landed in 6.40.0 and
      # 6.41.0 fixed a vpc_connectivity update regression introduced there;
      # the rebalancing block (Express brokers) is older (6.22.1). Everything
      # else this module uses predates the v6 line.
      version = ">= 6.41.0"
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
