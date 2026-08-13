locals {
  # The settled tag convention, matching the Pulumi module key-for-key:
  # user labels merge in FIRST so the Name + planton.ai/* identity keys can
  # never be overridden by a label. Labels reaching AWS as tags is what lets
  # a composition stamp cloud-side discovery tags on a subnet (e.g.
  # Karpenter's karpenter.sh/discovery tag, which its EC2NodeClass subnet
  # selector matches) without a per-kind tags field.
  aws_tags = merge(
    try(var.metadata.labels, {}),
    {
      "Name"                     = var.metadata.name
      "planton.ai/resource"      = "true"
      "planton.ai/organization"  = var.metadata.org
      "planton.ai/environment"   = var.metadata.env
      "planton.ai/resource-kind" = "AwsSubnet"
      "planton.ai/resource-id"   = var.metadata.id
    }
  )

  # The subnet owns a dedicated route table when it declares inline routes
  # and/or VGW propagation (spec CEL forbids either alongside an external
  # route_table_id).
  owns_route_table = length(var.spec.routes) > 0 || length(var.spec.propagating_vgws) > 0
}
