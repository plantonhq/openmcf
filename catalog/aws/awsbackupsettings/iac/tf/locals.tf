locals {
  # Resource-identity tags match the Pulumi module key-for-key. Neither
  # settings resource supports tags at the provider - the map exists so
  # the module keeps the catalog-wide shape if AWS ever adds them.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBackupSettings"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
