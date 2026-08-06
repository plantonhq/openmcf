terraform {
  required_providers {
    aws = {
      # Floor 6.48.0: warm_throughput_mib_ps landed there (max_record_size_in_kib
      # in 6.20.0, aws_kinesis_resource_policy well before the v6 line).
      source  = "hashicorp/aws"
      version = ">= 6.48.0"
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
