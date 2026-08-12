# VPC endpoints expose the workgroup inside other subnets --
# same-account cross-VPC access without peering. Each entry renders one
# aws_redshiftserverless_endpoint_access keyed by endpoint name; the
# endpoint's private address is exported per endpoint on the outputs
# contract. The cross-account grantee side (owner_account) is
# deliberately not modeled -- it lives in the grantee's credential
# domain.
resource "aws_redshiftserverless_endpoint_access" "this" {
  for_each = { for e in var.spec.endpoint_accesses : e.endpoint_name => e }

  workgroup_name = aws_redshiftserverless_workgroup.this.workgroup_name
  endpoint_name  = each.value.endpoint_name

  # An entry without its own subnets reuses the workgroup's (the spec
  # CEL guarantees the fallback exists).
  subnet_ids = length(each.value.subnet_ids) > 0 ? each.value.subnet_ids : var.spec.subnet_ids

  # Empty uses the VPC's default security group (the AWS default).
  vpc_security_group_ids = length(each.value.vpc_security_group_ids) > 0 ? each.value.vpc_security_group_ids : null
}

# The endpoint-access resource is not taggable in AWS -- no tags here.
