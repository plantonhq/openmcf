# The DB subnet group is a named list of subnets -- pure glue with no
# independent lifecycle, so it lives inside the module. The referenced
# subnets themselves are first-class AwsSubnet nodes this module never
# modifies. AWS requires the group to cover two availability zones even
# for a single-AZ instance (the CEL contract enforces two subnets).
resource "aws_db_subnet_group" "this" {
  count = local.manage_subnet_group ? 1 : 0

  name       = local.instance_identifier
  subnet_ids = var.spec.subnet_ids
  tags       = local.aws_tags
}
