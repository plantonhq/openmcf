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

  # The provider names the single association resource PUTs onto the
  # cluster: the Fargate built-ins plus every folded EC2 capacity
  # provider. Managed-instances providers are DELIBERATELY excluded:
  # PutClusterCapacityProviders does not manage them (AWS binds an MI
  # provider to its cluster at CreateCapacityProvider, and a PUT neither
  # attaches nor detaches it -- live-verified: a PUT omitting an attached
  # MI provider leaves it attached). Listing one also races its async
  # provisioning: the PUT fails 400 "not in an ACTIVE state" seconds
  # after create, and the pinned provider's retry list does not cover
  # that error.
  associated_capacity_providers = concat(
    var.spec.capacity_providers,
    [for p in var.spec.ec2_capacity_providers : p.name],
  )

  # The full capacity vocabulary services can name in a
  # capacity_provider_strategy -- every provider including the
  # managed-instances ones ECS attaches itself. This is the outputs
  # contract (matching the Pulumi module element-for-element), NOT the
  # association PUT's list.
  all_capacity_provider_names = concat(
    local.associated_capacity_providers,
    [for p in var.spec.managed_instances_capacity_providers : p.name],
  )
}
