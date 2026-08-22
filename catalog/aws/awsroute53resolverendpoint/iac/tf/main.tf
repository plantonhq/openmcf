# One Route 53 Resolver endpoint with its forwarding rules and their
# VPC associations managed in-line.
#
# Lifecycle facts the render below depends on:
#   - direction and security_group_ids are ForceNew on the endpoint;
#     ip_addresses churn in place (the provider adds before it removes
#     so the 2-address floor never breaks, one waiter round-trip per
#     change);
#   - a rule's resolver_endpoint_id updates in place EXCEPT detaching
#     it (endpoint id -> empty), which the provider forces to replace;
#     SYSTEM rules never carry one;
#   - rule associations are pure joins (rule, vpc) - every argument
#     ForceNew, no update path;
#   - AWS strips a trailing dot from rule domain names (the provider
#     normalizes both ways, so no drift);
#   - tags on a RAM-shared rule are never read back (the provider
#     skips them for SHARED_WITH_ME rules).

resource "aws_route53_resolver_endpoint" "this" {
  name               = var.metadata.name
  direction          = var.spec.direction
  security_group_ids = var.spec.security_group_ids

  resolver_endpoint_type = var.spec.endpoint_type != "" ? var.spec.endpoint_type : null
  protocols              = length(var.spec.protocols) > 0 ? var.spec.protocols : null

  # Tri-state metrics toggles: null leaves AWS's default, an explicit
  # value (true or false) is sent as stated.
  rni_enhanced_metrics_enabled       = var.spec.rni_enhanced_metrics_enabled
  target_name_server_metrics_enabled = var.spec.target_name_server_metrics_enabled

  dynamic "ip_address" {
    for_each = var.spec.ip_addresses
    content {
      subnet_id = ip_address.value.subnet_id
      ip        = ip_address.value.ip != "" ? ip_address.value.ip : null
      ipv6      = ip_address.value.ipv6 != "" ? ip_address.value.ipv6 : null
    }
  }

  tags = local.aws_tags
}

# Forwarding rules, keyed by rule name. SYSTEM rules carry no endpoint
# binding (they restore recursive resolution); FORWARD and DELEGATE
# rules bind to this endpoint.
resource "aws_route53_resolver_rule" "this" {
  for_each = { for rule in var.spec.rules : rule.name => rule }

  name        = each.value.name
  domain_name = each.value.domain_name
  rule_type   = each.value.rule_type

  resolver_endpoint_id = each.value.rule_type == "SYSTEM" ? null : aws_route53_resolver_endpoint.this.id

  dynamic "target_ip" {
    for_each = each.value.target_ips
    content {
      ip       = target_ip.value.ip != "" ? target_ip.value.ip : null
      ipv6     = target_ip.value.ipv6 != "" ? target_ip.value.ipv6 : null
      port     = target_ip.value.port > 0 ? target_ip.value.port : null
      protocol = target_ip.value.protocol != "" ? target_ip.value.protocol : null
    }
  }

  tags = local.aws_tags
}

locals {
  # One association per (rule, vpc) pair, keyed "rule//vpc" - the "//"
  # separator is the import bridge's address-key-segment convention.
  rule_vpc_pairs = flatten([
    for rule in var.spec.rules : [
      for vpc_id in rule.vpc_ids : {
        key       = "${rule.name}//${vpc_id}"
        rule_name = rule.name
        vpc_id    = vpc_id
      }
    ]
  ])
}

resource "aws_route53_resolver_rule_association" "this" {
  for_each = { for pair in local.rule_vpc_pairs : pair.key => pair }

  # Deterministic cosmetic name, identical across both engines and
  # capped at the provider's 64-character wall (Pulumi would otherwise
  # auto-generate one from the resource URN).
  name = substr("${each.value.rule_name}-${each.value.vpc_id}", 0, 64)

  resolver_rule_id = aws_route53_resolver_rule.this[each.value.rule_name].id
  vpc_id           = each.value.vpc_id
}
