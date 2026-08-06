# The group and its rules are one resource: rules are authored INLINE on
# aws_security_group. AWS forbids mixing inline rules with standalone
# aws_vpc_security_group_*_rule resources on the same group (each apply would
# fight to own the rule set), so this module must never emit standalone rule
# resources -- the inline blocks below are the single owner of the group's
# rules.
resource "aws_security_group" "this" {
  # The group name is metadata.name -- create-only in AWS, and the basis both
  # engines share so a manifest deploys identically on either.
  name        = local.security_group_name
  description = local.description
  vpc_id      = local.vpc_id

  # Forcibly revoke this group's rules -- and rules in OTHER groups that
  # reference it -- before delete. Without it, destroying a group that a
  # sibling group still references fails with a DependencyViolation and
  # requires manual rule surgery.
  revoke_rules_on_delete = var.spec.revoke_rules_on_delete

  # Each spec rule maps to exactly one inline rule block. A rule may carry
  # several sources at once (CIDRs, prefix lists, other groups, self); AWS
  # expands them into individual permissions server-side.
  dynamic "ingress" {
    for_each = var.spec.ingress
    content {
      description      = ingress.value.description
      protocol         = ingress.value.protocol
      from_port        = ingress.value.from_port
      to_port          = ingress.value.to_port
      cidr_blocks      = ingress.value.ipv4_cidrs
      ipv6_cidr_blocks = ingress.value.ipv6_cidrs
      # Managed prefix lists let a rule target a named CIDR set (an AWS
      # service like S3, or a customer-maintained office/partner range) by
      # stable ID instead of hardcoding addresses.
      prefix_list_ids = ingress.value.prefix_list_ids
      security_groups = ingress.value.source_security_group_ids
      self            = ingress.value.self_reference
    }
  }

  # With inline egress, an empty list means DENY ALL outbound: the provider
  # revokes the allow-all egress rule AWS adds to every new group, so the
  # manifest is the complete statement of what the group permits.
  dynamic "egress" {
    for_each = var.spec.egress
    content {
      description      = egress.value.description
      protocol         = egress.value.protocol
      from_port        = egress.value.from_port
      to_port          = egress.value.to_port
      cidr_blocks      = egress.value.ipv4_cidrs
      ipv6_cidr_blocks = egress.value.ipv6_cidrs
      prefix_list_ids  = egress.value.prefix_list_ids
      security_groups  = egress.value.destination_security_group_ids
      self             = egress.value.self_reference
    }
  }

  tags = local.aws_tags
}
