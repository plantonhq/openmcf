# One network ACL with its rules and subnet associations managed
# in-line.
#
# Lifecycle facts the render below depends on:
#   - vpc_id is fixed for life (replace-on-change);
#   - in-line ingress/egress rules are the single declarative owner;
#     the standalone rule/association resources are identical payloads
#     and fight this form, so this module never uses them;
#   - AWS stores protocols as NUMBERS and the provider normalizes
#     names in its rule hash, so "tcp" here never causes a perpetual
#     diff; AWS's own catch-all rules (32767/32768) are invisible and
#     unmanageable;
#   - a subnet listed here is atomically REPLACED onto this ACL (AWS
#     has no attach - only ReplaceNetworkAclAssociation); removing it
#     hands it back to the VPC's default NACL;
#   - destroy tears down all subnet associations first, then deletes
#     (retrying on DependencyViolation).

resource "aws_network_acl" "this" {
  vpc_id = var.spec.vpc_id

  dynamic "ingress" {
    for_each = var.spec.ingress != null ? var.spec.ingress : []
    content {
      rule_no         = ingress.value.rule_no
      action          = ingress.value.action
      protocol        = ingress.value.protocol
      cidr_block      = ingress.value.cidr_block != "" ? ingress.value.cidr_block : null
      ipv6_cidr_block = ingress.value.ipv6_cidr_block != "" ? ingress.value.ipv6_cidr_block : null
      from_port       = ingress.value.from_port
      to_port         = ingress.value.to_port
      icmp_type       = ingress.value.icmp_type != null ? ingress.value.icmp_type : null
      icmp_code       = ingress.value.icmp_code != null ? ingress.value.icmp_code : null
    }
  }

  dynamic "egress" {
    for_each = var.spec.egress != null ? var.spec.egress : []
    content {
      rule_no         = egress.value.rule_no
      action          = egress.value.action
      protocol        = egress.value.protocol
      cidr_block      = egress.value.cidr_block != "" ? egress.value.cidr_block : null
      ipv6_cidr_block = egress.value.ipv6_cidr_block != "" ? egress.value.ipv6_cidr_block : null
      from_port       = egress.value.from_port
      to_port         = egress.value.to_port
      icmp_type       = egress.value.icmp_type != null ? egress.value.icmp_type : null
      icmp_code       = egress.value.icmp_code != null ? egress.value.icmp_code : null
    }
  }

  subnet_ids = var.spec.subnet_ids != null && length(var.spec.subnet_ids) > 0 ? var.spec.subnet_ids : null

  tags = local.aws_tags
}
