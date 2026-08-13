# Custom cluster endpoints: one resource per spec.custom_endpoints
# entry, keyed by the entry's name so endpoints come and go without
# touching the cluster or its instances. Members reference the folded
# instances BY THEIR SPEC NAMES -- indexing the instance resources maps
# each name to its full AWS identifier and makes a typo'd member fail
# the plan instead of silently fronting nothing.
resource "aws_rds_cluster_endpoint" "this" {
  for_each = { for endpoint in var.spec.custom_endpoints : endpoint.name => endpoint }

  cluster_identifier          = aws_rds_cluster.this.cluster_identifier
  cluster_endpoint_identifier = "${local.cluster_identifier}-${each.value.name}"
  custom_endpoint_type        = each.value.type

  static_members   = length(each.value.static_members) > 0 ? [for m in each.value.static_members : aws_rds_cluster_instance.this[m].identifier] : null
  excluded_members = length(each.value.excluded_members) > 0 ? [for m in each.value.excluded_members : aws_rds_cluster_instance.this[m].identifier] : null

  tags = local.aws_tags
}
