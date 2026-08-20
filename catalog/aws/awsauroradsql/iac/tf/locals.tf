locals {
  # Resource-identity tags match the Pulumi module key-for-key. DSQL
  # generates its own cluster identifier, so the Name tag is how
  # humans find this cluster in the console.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsAuroraDsql"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
