# Folded custom endpoints: an endpoint's identity IS this cluster (it
# fronts a subset of the cluster's own instances), so endpoints live here
# rather than as their own kind. Each entry is keyed by its name --
# adding, renaming, or removing one endpoint never touches its siblings
# or the cluster. Member lists name spec.instances entries and are
# translated to the derived AWS instance identifiers the same way the
# instance resources build them.
resource "aws_neptune_cluster_endpoint" "this" {
  for_each = { for endpoint in var.spec.custom_endpoints : endpoint.name => endpoint }

  cluster_identifier          = aws_neptune_cluster.this.id
  cluster_endpoint_identifier = each.value.name
  endpoint_type               = each.value.endpoint_type

  # Pin to exactly these instances (by their spec.instances entry names,
  # mapped to the derived identifiers)...
  static_members = length(each.value.static_members) > 0 ? [
    for member in each.value.static_members : aws_neptune_cluster_instance.this[member].identifier
  ] : null

  # ...or front every instance of the endpoint's type EXCEPT these
  # (mutually exclusive with static_members -- CEL-enforced).
  excluded_members = length(each.value.excluded_members) > 0 ? [
    for member in each.value.excluded_members : aws_neptune_cluster_instance.this[member].identifier
  ] : null

  tags = local.aws_tags
}
