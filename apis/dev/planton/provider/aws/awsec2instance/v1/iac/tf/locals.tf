locals {
  # EC2 instances have no name argument -- the Name tag IS the display
  # identity -- so metadata.name travels in the tag set and a manifest
  # deploys identically on either engine.
  instance_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEc2Instance"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
