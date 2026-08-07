resource "aws_wafv2_ip_set" "this" {
  # The set's AWS name is the Planton resource name -- the stable identity
  # web ACL statements and operators see. Name, scope, and address version
  # are all create-time immutable (ForceNew).
  name               = var.metadata.name
  scope              = var.spec.scope
  ip_address_version = var.spec.ip_address_version

  # CIDR entries only (a bare address is rejected by AWS -- the spec's CEL
  # enforces the /nn suffix up front). An empty list is a valid placeholder
  # set that matches nothing until ranges are added.
  addresses = var.spec.addresses

  description = local.description

  tags = local.aws_tags
}
