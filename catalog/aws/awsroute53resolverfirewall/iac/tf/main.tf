# One DNS Firewall rule group with its domain lists, rules, and VPC
# associations managed in-line.
#
# Lifecycle facts the render below depends on:
#   - the group itself is a name-and-tags container (name ForceNew,
#     update path is tags-only);
#   - domain list contents push through a separate update call after
#     create (ADD on create, REPLACE on update) - a partially failed
#     import surfaces as a retry error, not silent success;
#   - a rule's match source (domain list vs threat protection) is
#     ForceNew either way; action, priority, and the block-response
#     shape update in place;
#   - the provider enforces the standard-XOR-advanced rule shape at
#     plan time (mirrored in the spec's CELs);
#   - block_override_dns_type has exactly one legal value (CNAME), so
#     the module pins it whenever the response is OVERRIDE - a dead
#     knob is never surfaced;
#   - associations update priority/mutation_protection in place.

resource "aws_route53_resolver_firewall_rule_group" "this" {
  name = var.metadata.name
  tags = local.aws_tags
}

# Owned domain lists, keyed by list name.
resource "aws_route53_resolver_firewall_domain_list" "this" {
  for_each = { for list in var.spec.domain_lists : list.name => list }

  name    = each.value.name
  domains = length(each.value.domains) > 0 ? each.value.domains : null

  tags = local.aws_tags
}

# Filtering rules, keyed by rule name. The match source resolves to an
# owned list's generated ID, a literal external/managed list ID, or
# the threat-protection arm.
resource "aws_route53_resolver_firewall_rule" "this" {
  for_each = { for rule in var.spec.rules : rule.name => rule }

  name                   = each.value.name
  firewall_rule_group_id = aws_route53_resolver_firewall_rule_group.this.id
  priority               = each.value.priority
  action                 = each.value.action

  firewall_domain_list_id = each.value.domain_list_name != "" ? aws_route53_resolver_firewall_domain_list.this[each.value.domain_list_name].id : (each.value.domain_list_id != "" ? each.value.domain_list_id : null)

  dns_threat_protection = each.value.dns_threat_protection != "" ? each.value.dns_threat_protection : null
  confidence_threshold  = each.value.confidence_threshold != "" ? each.value.confidence_threshold : null

  block_response        = each.value.block_response != "" ? each.value.block_response : null
  block_override_domain = each.value.block_override_domain != "" ? each.value.block_override_domain : null
  block_override_ttl    = each.value.block_override_ttl
  # The one legal override record type - module-owned constant.
  block_override_dns_type = each.value.block_response == "OVERRIDE" ? "CNAME" : null

  firewall_domain_redirection_action = each.value.firewall_domain_redirection_action != "" ? each.value.firewall_domain_redirection_action : null
  q_type                             = each.value.q_type != "" ? each.value.q_type : null
}

# VPC associations, keyed by association name.
resource "aws_route53_resolver_firewall_rule_group_association" "this" {
  for_each = { for assoc in var.spec.vpc_associations : assoc.name => assoc }

  name                   = each.value.name
  firewall_rule_group_id = aws_route53_resolver_firewall_rule_group.this.id
  vpc_id                 = each.value.vpc_id
  priority               = each.value.priority

  mutation_protection = each.value.mutation_protection != "" ? each.value.mutation_protection : null

  tags = local.aws_tags
}
