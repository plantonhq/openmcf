# Engine feature roles: one association resource per spec.iam_roles
# entry, keyed by the role ARN so roles attach and detach without
# touching the instance. AWS keys associations by (instance, role) and
# REQUIRES the feature name here (unlike the cluster-side association,
# where it is optional) -- the spec's CEL mirrors that asymmetry.
resource "aws_db_instance_role_association" "this" {
  for_each = { for entry in var.spec.iam_roles : entry.role => entry }

  db_instance_identifier = aws_db_instance.this.identifier
  role_arn               = each.value.role
  feature_name           = each.value.feature_name
}
