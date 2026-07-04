locals {
  # The cluster name is metadata.name -- create-only in AWS (max 64 chars), and
  # the basis both engines share so a manifest deploys identically on either.
  cluster_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Tags are the
  # ONLY mutable surface on a serverless MSK cluster -- everything else is
  # create-time (ForceNew).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsMskServerlessCluster"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
