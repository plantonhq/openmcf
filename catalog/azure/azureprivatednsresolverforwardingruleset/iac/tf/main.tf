# Create the DNS forwarding ruleset -- the rule book that steers DNS
# queries for chosen domains out of Azure -- and its forwarding rules
# as child resources. The ruleset binds the outbound endpoint(s) of a
# DNS Private Resolver (at most 2, both from the SAME resolver --
# Azure enforces it at deploy time) and takes effect in a network only
# once that network is linked to it
# (AzurePrivateDnsResolverVirtualNetworkLink). The endpoint list and
# tags update in place; name, resource group, and region replace the
# ruleset. Rulesets and rules are free at rest.
resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "main" {
  name                                       = var.spec.name
  resource_group_name                        = var.spec.resource_group
  location                                   = var.spec.region
  private_dns_resolver_outbound_endpoint_ids = var.spec.outbound_endpoint_ids
  tags                                       = local.final_tags
}

# The forwarding rules -- one per captured domain. Keyed by rule name
# (spec-validated unique) so adding or removing one rule never touches
# its siblings. Everything on a rule updates in place EXCEPT
# domain_name, which replaces that rule.
resource "azurerm_private_dns_resolver_forwarding_rule" "main" {
  for_each = { for rule in var.spec.forwarding_rules : rule.name => rule }

  name                      = each.value.name
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.main.id
  # ARM stores domains as fully qualified names WITH the trailing dot
  # ("corp.contoso.com.") -- write them that way in the spec.
  domain_name = each.value.domain_name
  enabled     = coalesce(each.value.enabled, true)

  # Azure tries the targets in order (up to 6 per rule). The port is
  # always sent explicitly -- 53 (the standard DNS port and ARM's
  # default) when the spec leaves it unset -- so plans stay
  # deterministic.
  dynamic "target_dns_servers" {
    for_each = each.value.target_dns_servers
    content {
      ip_address = target_dns_servers.value.ip_address
      port       = coalesce(target_dns_servers.value.port, 53)
    }
  }

  # ARM's free-form annotation map on the rule itself (rules carry no
  # tags).
  metadata = length(each.value.metadata) > 0 ? each.value.metadata : null
}
