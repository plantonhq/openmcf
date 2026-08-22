locals {
  # Reference fields (droplet IDs, kubernetes IDs, load balancer UIDs) arrive
  # flattened as plain strings: the Planton orchestrator resolves valueFrom
  # references before Terraform runs. Droplet IDs are numeric in the
  # DigitalOcean API, so they convert with tonumber here.
  droplet_ids = [for id in var.spec.droplet_ids : tonumber(id)]

  inbound_rules = [
    for rule in var.spec.inbound_rules : {
      protocol                  = rule.protocol
      port_range                = rule.port_range != null ? rule.port_range : ""
      source_addresses          = rule.source_addresses != null ? rule.source_addresses : []
      source_droplet_ids        = rule.source_droplet_ids != null ? [for id in rule.source_droplet_ids : tonumber(id)] : []
      source_tags               = rule.source_tags != null ? rule.source_tags : []
      source_kubernetes_ids     = rule.source_kubernetes_ids != null ? rule.source_kubernetes_ids : []
      source_load_balancer_uids = rule.source_load_balancer_uids != null ? rule.source_load_balancer_uids : []
    }
  ]

  outbound_rules = [
    for rule in var.spec.outbound_rules : {
      protocol                       = rule.protocol
      port_range                     = rule.port_range != null ? rule.port_range : ""
      destination_addresses          = rule.destination_addresses != null ? rule.destination_addresses : []
      destination_droplet_ids        = rule.destination_droplet_ids != null ? [for id in rule.destination_droplet_ids : tonumber(id)] : []
      destination_tags               = rule.destination_tags != null ? rule.destination_tags : []
      destination_kubernetes_ids     = rule.destination_kubernetes_ids != null ? rule.destination_kubernetes_ids : []
      destination_load_balancer_uids = rule.destination_load_balancer_uids != null ? rule.destination_load_balancer_uids : []
    }
  ]
}
