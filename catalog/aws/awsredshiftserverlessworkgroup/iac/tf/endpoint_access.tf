# VPC endpoints expose the workgroup inside other subnets --
# same-account cross-VPC access without peering. Each entry renders one
# aws_redshiftserverless_endpoint_access keyed by endpoint name; the
# endpoint's private address is exported per endpoint on the outputs
# contract. The cross-account grantee side (owner_account) is
# deliberately not modeled -- it lives in the grantee's credential
# domain.
resource "aws_redshiftserverless_endpoint_access" "this" {
  for_each = { for e in var.spec.endpoint_accesses : e.endpoint_name => e }

  # Endpoint-access CREATE and DELETE both answer 400 ConflictException
  # ("An operation is running on the serverless workgroup") unless the
  # workgroup is idle, and the provider carries no ConflictException
  # retry on either (only on the workgroup's own delete/update).
  # Endpoint accesses therefore apply straight after the workgroup
  # (idle from the provider's own wait-for-available; the custom
  # domain's window is unproven -- its live arm is deferred), with the
  # conflict-immune usage limits chained behind them across the settle
  # sleep that protects the destroy direction (the full live-probed
  # contract lives in satellite_settle.tf). Entries WITHIN this group
  # still dispatch concurrently -- a durable fix for many-entry specs
  # is a provider-side retry, recorded upstream.
  depends_on = [aws_redshiftserverless_custom_domain_association.this]

  workgroup_name = aws_redshiftserverless_workgroup.this.workgroup_name
  endpoint_name  = each.value.endpoint_name

  # An entry without its own subnets reuses the workgroup's (the spec
  # CEL guarantees the fallback exists).
  subnet_ids = length(each.value.subnet_ids) > 0 ? each.value.subnet_ids : var.spec.subnet_ids

  # Empty uses the VPC's default security group (the AWS default).
  vpc_security_group_ids = length(each.value.vpc_security_group_ids) > 0 ? each.value.vpc_security_group_ids : null
}

# The endpoint-access resource is not taggable in AWS -- no tags here.
