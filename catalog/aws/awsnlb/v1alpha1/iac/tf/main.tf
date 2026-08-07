# The Network Load Balancer carries no routing configuration by design:
# listeners and target groups are separate resources that attach to it by
# ARN, so this module owns only what is truly load-balancer-wide -- node
# placement with optional static IPs, security groups, and traffic
# distribution behavior. Changing "internal" replaces the load balancer.
resource "aws_lb" "this" {
  name               = local.nlb_name
  load_balancer_type = "network"
  internal           = var.spec.internal

  enable_deletion_protection = var.spec.delete_protection_enabled

  # Cross-zone distribution is a real cost decision on NLB (inter-AZ data
  # transfer is billed), which is why AWS defaults it off and the spec makes
  # it an explicit opt-in.
  enable_cross_zone_load_balancing = var.spec.cross_zone_load_balancing_enabled

  # Only explicitly set attributes are sent, so AWS keeps its own defaults
  # for the rest -- the module never bakes in opinions the spec does not
  # express.
  ip_address_type                  = var.spec.ip_address_type != "" ? var.spec.ip_address_type : null
  dns_record_client_routing_policy = var.spec.dns_record_client_routing_policy != "" ? var.spec.dns_record_client_routing_policy : null
  enable_zonal_shift               = var.spec.zonal_shift_enabled ? true : null

  enforce_security_group_inbound_rules_on_private_link_traffic = var.spec.enforce_security_group_inbound_rules_on_private_link_traffic != "" ? var.spec.enforce_security_group_inbound_rules_on_private_link_traffic : null

  # Each mapping pins one NLB node to a subnet, optionally with a static
  # Elastic IP (internet-facing) or a fixed private address (internal) --
  # the static-IP story that differentiates NLB from ALB.
  dynamic "subnet_mapping" {
    for_each = var.spec.subnet_mappings
    content {
      subnet_id            = subnet_mapping.value.subnet_id
      allocation_id        = subnet_mapping.value.allocation_id != "" ? subnet_mapping.value.allocation_id : null
      private_ipv4_address = subnet_mapping.value.private_ipv4_address != "" ? subnet_mapping.value.private_ipv4_address : null
    }
  }

  # Optional for NLB (unlike ALB) -- and once attached, AWS never allows
  # removing the last one, so attaching any group is a one-way door.
  security_groups = length(var.spec.security_groups) > 0 ? var.spec.security_groups : null

  # NLB access logs only capture TLS-listener traffic (an AWS limitation);
  # "enabled" is implied by the block's presence in the spec. The bucket must
  # carry the ELB log-delivery bucket policy or delivery fails silently.
  dynamic "access_logs" {
    for_each = var.spec.access_logs != null ? [var.spec.access_logs] : []
    content {
      bucket  = access_logs.value.bucket
      prefix  = access_logs.value.prefix != "" ? access_logs.value.prefix : null
      enabled = true
    }
  }

  tags = local.aws_tags
}

# Optional Route53 records for each hostname when DNS is enabled.
# allow_overwrite adopts an existing alias record (e.g. left by a prior partial apply,
# or one already pointing at this NLB) instead of failing the apply on a CREATE
# collision -- this alias record is owned by the NLB module.
resource "aws_route53_record" "this" {
  # toset([]) rather than [] in the false arm: both arms of a conditional
  # must have one type, and a bare [] is a tuple -- for_each rejects it.
  for_each = local.create_dns_records ? toset(var.spec.dns.hostnames) : toset([])

  allow_overwrite = true
  zone_id         = var.spec.dns.route53_zone_id
  name            = each.value
  type            = "A"

  alias {
    name    = aws_lb.this.dns_name
    zone_id = aws_lb.this.zone_id
    # false on purpose: target-health evaluation only changes behavior under
    # failover/weighted routing policies, and a simple alias should not pay
    # for health evaluation. Must stay identical in the Pulumi module
    # (cross-engine parity).
    evaluate_target_health = false
  }
}
