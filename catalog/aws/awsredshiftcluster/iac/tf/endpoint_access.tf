# Redshift-managed VPC endpoints expose the cluster inside other subnet
# groups -- same-account cross-VPC access without peering (RA3 only).
# Each entry renders one aws_redshift_endpoint_access keyed by endpoint
# name; the endpoint's private address is exported per endpoint on the
# outputs contract. The cross-account grantee side (resource_owner) is
# deliberately not modeled here -- a grantee creates its endpoint in its
# own account against an authorization this cluster grants (see
# endpoint_authorization.tf).
resource "aws_redshift_endpoint_access" "this" {
  for_each = { for e in var.spec.endpoint_accesses : e.endpoint_name => e }

  cluster_identifier = aws_redshift_cluster.this.cluster_identifier
  endpoint_name      = each.value.endpoint_name

  # An entry without its own group reuses the cluster's subnet group
  # (managed or referenced), yielding an extra endpoint in the
  # cluster's own VPC.
  subnet_group_name = each.value.subnet_group_name != "" ? each.value.subnet_group_name : local.effective_subnet_group

  # Empty uses the VPC's default security group (the AWS default).
  vpc_security_group_ids = length(each.value.vpc_security_group_ids) > 0 ? each.value.vpc_security_group_ids : null
}

# The endpoint-access resource is not taggable in AWS -- no tags here.
