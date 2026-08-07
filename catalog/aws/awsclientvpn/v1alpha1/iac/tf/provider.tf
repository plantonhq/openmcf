terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor, not a cap: endpoint_ip_address_type / traffic_ip_address_type
      # (and client_cidr_block becoming optional for IPv6 traffic) land in
      # v6.11.0 -- an older floor would silently reject the dual-stack
      # surface. (client_route_enforcement_options and
      # disconnect_on_session_timeout predate it.)
      version = ">= 6.11.0"
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
