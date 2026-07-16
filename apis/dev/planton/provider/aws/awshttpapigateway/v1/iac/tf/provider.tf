terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      # Floor 6.29.0: the API/domain ip_address_type argument landed on the
      # 5.97 line, but 6.29.0 is where the apigatewayv2 family last changed
      # shape (domain routing_mode / routing rules) -- pinning the family
      # floor there keeps every argument this module and its siblings use
      # resolvable from one provider build.
      source  = "hashicorp/aws"
      version = ">= 6.29.0"
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
