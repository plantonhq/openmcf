# Endpoint authorizations grant OTHER AWS accounts permission to create
# managed VPC endpoints to this cluster -- the grantor side of
# cross-account access, living in this cluster's own credential domain.
# Each entry renders one aws_redshift_endpoint_authorization keyed by
# the grantee account (AWS keeps one authorization per account).
resource "aws_redshift_endpoint_authorization" "this" {
  for_each = { for a in var.spec.endpoint_authorizations : a.account => a }

  cluster_identifier = aws_redshift_cluster.this.cluster_identifier
  account            = each.value.account

  # Empty authorizes ALL of the grantee account's VPCs.
  vpc_ids = length(each.value.vpc_ids) > 0 ? each.value.vpc_ids : null

  # When true, delete revokes the grant even while the grantee still
  # has live endpoints (deleting them too). AWS's default refuses.
  force_delete = each.value.force_delete
}
