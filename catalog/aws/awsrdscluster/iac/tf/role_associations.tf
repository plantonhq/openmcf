# Engine feature roles: one association resource per spec.iam_roles
# entry, keyed by the role ARN so roles attach and detach without
# touching the cluster or each other. AWS keys associations by
# (cluster, role) -- one association per role -- and the optional
# feature_name links the role to a specific engine capability
# (s3Import, Lambda, SageMaker, ...). The inline iam_roles argument on
# aws_rds_cluster is deliberately unused: it cannot carry feature names
# and, per the provider's own warning, mixing it with association
# resources overwrites the associations.
resource "aws_rds_cluster_role_association" "this" {
  for_each = { for entry in var.spec.iam_roles : entry.role => entry }

  db_cluster_identifier = aws_rds_cluster.this.cluster_identifier
  role_arn              = each.value.role
  feature_name          = each.value.feature_name != "" ? each.value.feature_name : null
}
