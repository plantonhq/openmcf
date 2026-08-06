terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Floor 6.1.0: the earliest release whose vendored S3 SDK accepts every
      # enum value the spec allows — the FSx-flavored storage classes
      # (FSX_OPENZFS / FSX_ONTAP) and the aws:fsx encryption value land in
      # 6.1.0's SDK; the provider validates these client-side, so an older
      # provider would reject manifests using them. Everything else on
      # aws_s3_object predates the v6 line.
      version = ">= 6.1.0"
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
