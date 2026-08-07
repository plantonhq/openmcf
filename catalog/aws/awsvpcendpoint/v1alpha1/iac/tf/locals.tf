locals {
  # VPC endpoints carry no name parameter in AWS -- identity lives
  # entirely in tags, so the Name tag is what the console displays.
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsVpcEndpoint"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
