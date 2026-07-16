# The Neptune subnet group is a named list of subnets -- pure glue with
# no independent lifecycle, so it lives inside the module. The
# referenced subnets themselves are first-class AwsSubnet nodes this
# module never modifies. Changing the group's subnet membership updates
# in place; moving the CLUSTER to a different group replaces the cluster
# (AWS create-time constraint).
resource "aws_neptune_subnet_group" "this" {
  count = local.manage_subnet_group ? 1 : 0

  name       = local.cluster_identifier
  subnet_ids = var.spec.subnet_ids
  tags       = local.aws_tags
}
