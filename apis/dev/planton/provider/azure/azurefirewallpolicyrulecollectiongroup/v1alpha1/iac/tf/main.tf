# Create the rule collection group -- an ordered document of application,
# network, and DNAT rule collections nested under its parent firewall
# policy. A policy carries many groups, each deployed independently, which
# is why the group (not the rule or the policy) is the deployment unit.
#
# Azure evaluates groups by group priority, collections by collection
# priority WITHIN a type -- and across types always DNAT -> network ->
# application, regardless of priorities. Lower numbers run first.
#
# The provider serializes writes to the same parent policy (named locks),
# so concurrent group deployments against one policy queue rather than
# conflict.
resource "azurerm_firewall_policy_rule_collection_group" "main" {
  name               = var.spec.name
  firewall_policy_id = var.spec.firewall_policy_id
  priority           = var.spec.priority

  dynamic "application_rule_collection" {
    for_each = var.spec.application_rule_collections
    content {
      name     = application_rule_collection.value.name
      priority = application_rule_collection.value.priority
      action   = lookup(local.filter_action_wire, application_rule_collection.value.action, application_rule_collection.value.action)

      dynamic "rule" {
        for_each = application_rule_collection.value.rules
        content {
          name        = rule.value.name
          description = rule.value.description

          dynamic "protocols" {
            for_each = rule.value.protocols
            content {
              type = lookup(local.application_protocol_type_wire, protocols.value.type, protocols.value.type)
              port = protocols.value.port
            }
          }

          dynamic "http_headers" {
            for_each = rule.value.http_headers
            content {
              name  = http_headers.value.name
              value = http_headers.value.value
            }
          }

          source_addresses      = rule.value.source_addresses
          source_ip_groups      = rule.value.source_ip_groups
          destination_addresses = rule.value.destination_addresses
          destination_fqdns     = rule.value.destination_fqdns
          destination_urls      = rule.value.destination_urls
          destination_fqdn_tags = rule.value.destination_fqdn_tags
          # terminate_tls decrypts for inspection (Premium + policy TLS
          # certificate); required for URL (path-level) filtering of
          # HTTPS traffic.
          terminate_tls  = rule.value.terminate_tls
          web_categories = rule.value.web_categories
        }
      }
    }
  }

  dynamic "network_rule_collection" {
    for_each = var.spec.network_rule_collections
    content {
      name     = network_rule_collection.value.name
      priority = network_rule_collection.value.priority
      action   = lookup(local.filter_action_wire, network_rule_collection.value.action, network_rule_collection.value.action)

      dynamic "rule" {
        for_each = network_rule_collection.value.rules
        content {
          name        = rule.value.name
          description = rule.value.description
          protocols = [
            for protocol in rule.value.protocols :
            lookup(local.rule_protocol_wire, protocol, protocol)
          ]
          source_addresses      = rule.value.source_addresses
          source_ip_groups      = rule.value.source_ip_groups
          destination_addresses = rule.value.destination_addresses
          destination_ip_groups = rule.value.destination_ip_groups
          # FQDN destinations require the policy's DNS proxy so the
          # firewall resolves names the same way clients do.
          destination_fqdns = rule.value.destination_fqdns
          destination_ports = rule.value.destination_ports
        }
      }
    }
  }

  dynamic "nat_rule_collection" {
    for_each = var.spec.nat_rule_collections
    content {
      name     = nat_rule_collection.value.name
      priority = nat_rule_collection.value.priority
      # The DNAT action vocabulary has exactly one value -- a constant,
      # not a knob -- so the module sends the provider's schema literal
      # "Dnat" unconditionally (ARM normalizes case).
      action = "Dnat"

      dynamic "rule" {
        for_each = nat_rule_collection.value.rules
        content {
          name        = rule.value.name
          description = rule.value.description
          protocols = [
            for protocol in rule.value.protocols :
            lookup(local.rule_protocol_wire, protocol, protocol)
          ]
          source_addresses    = rule.value.source_addresses
          source_ip_groups    = rule.value.source_ip_groups
          destination_address = rule.value.destination_address
          # ARM caps DNAT destination ports at ONE entry (a port or a
          # range, no wildcard); the provider models the same cap.
          destination_ports = rule.value.destination_ports
          # Exactly one translation target (spec-validated); the unset
          # one stays null and is omitted from the payload.
          translated_address = rule.value.translated_address
          translated_fqdn    = rule.value.translated_fqdn
          translated_port    = rule.value.translated_port
        }
      }
    }
  }
}
