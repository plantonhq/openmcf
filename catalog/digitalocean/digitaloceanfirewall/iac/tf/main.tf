# DigitalOcean Cloud Firewall: a stateful, default-deny rule set applied to
# Droplets directly by ID or dynamically by Droplet tag. Empty rule leaves are
# sent as null (not empty sets): the provider omits absent collections when it
# reads state back, so sending [] would create permanent diffs on the
# set-hashed rule blocks.
resource "digitalocean_firewall" "firewall" {
  name = var.spec.firewall_name

  droplet_ids = length(local.droplet_ids) > 0 ? local.droplet_ids : null
  tags        = length(var.spec.tags) > 0 ? var.spec.tags : null

  dynamic "inbound_rule" {
    for_each = local.inbound_rules
    content {
      protocol                  = inbound_rule.value.protocol
      port_range                = inbound_rule.value.port_range != "" ? inbound_rule.value.port_range : null
      source_addresses          = length(inbound_rule.value.source_addresses) > 0 ? inbound_rule.value.source_addresses : null
      source_droplet_ids        = length(inbound_rule.value.source_droplet_ids) > 0 ? inbound_rule.value.source_droplet_ids : null
      source_tags               = length(inbound_rule.value.source_tags) > 0 ? inbound_rule.value.source_tags : null
      source_kubernetes_ids     = length(inbound_rule.value.source_kubernetes_ids) > 0 ? inbound_rule.value.source_kubernetes_ids : null
      source_load_balancer_uids = length(inbound_rule.value.source_load_balancer_uids) > 0 ? inbound_rule.value.source_load_balancer_uids : null
    }
  }

  dynamic "outbound_rule" {
    for_each = local.outbound_rules
    content {
      protocol                       = outbound_rule.value.protocol
      port_range                     = outbound_rule.value.port_range != "" ? outbound_rule.value.port_range : null
      destination_addresses          = length(outbound_rule.value.destination_addresses) > 0 ? outbound_rule.value.destination_addresses : null
      destination_droplet_ids        = length(outbound_rule.value.destination_droplet_ids) > 0 ? outbound_rule.value.destination_droplet_ids : null
      destination_tags               = length(outbound_rule.value.destination_tags) > 0 ? outbound_rule.value.destination_tags : null
      destination_kubernetes_ids     = length(outbound_rule.value.destination_kubernetes_ids) > 0 ? outbound_rule.value.destination_kubernetes_ids : null
      destination_load_balancer_uids = length(outbound_rule.value.destination_load_balancer_uids) > 0 ? outbound_rule.value.destination_load_balancer_uids : null
    }
  }
}
