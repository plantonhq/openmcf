terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: dns_options.private_dns_preference and
      # private_dns_specified_domains landed in provider 6.28.0, and the
      # Resource / ServiceNetwork endpoint types are v6-era -- an older
      # floor would silently reject parts of the modeled surface.
      version = ">= 6.28.0"
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
