locals {
  # The cluster name is metadata.name -- create-only in AWS (changing it
  # replaces the cluster), and the basis both engines share so a manifest
  # deploys identically on either.
  cluster_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEcsCluster"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The union of associated provider names: the Fargate built-ins plus
  # every folded EC2 capacity provider. This is the vocabulary services
  # can name in a capacity_provider_strategy, and what the single
  # association resource PUTs onto the cluster.
  associated_capacity_providers = concat(
    var.spec.capacity_providers,
    [for p in var.spec.ec2_capacity_providers : p.name],
  )
}
