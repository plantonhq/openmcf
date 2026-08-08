# Create the Private Link Service -- the PROVIDER side of Azure Private
# Link: consumers in other virtual networks (or other tenants) reach the
# service behind it through private endpoints, over the Microsoft
# backbone, with no peering and no public exposure.
resource "azurerm_private_link_service" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # The NAT configurations consumer traffic is source-NATed through.
  # Every subnet here must have private-link-service network policies
  # DISABLED -- ARM rejects the configuration otherwise. Setting a static
  # private address flips ARM's allocation method to Static; the provider
  # derives the method from the address's presence, so the module only
  # carries the address itself.
  dynamic "nat_ip_configuration" {
    for_each = var.spec.nat_ip_configurations
    content {
      name                       = nat_ip_configuration.value.name
      subnet_id                  = nat_ip_configuration.value.subnet_id
      primary                    = nat_ip_configuration.value.primary
      private_ip_address         = nat_ip_configuration.value.private_ip_address != "" ? nat_ip_configuration.value.private_ip_address : null
      private_ip_address_version = coalesce(nat_ip_configuration.value.private_ip_address_version, "IPv4")
    }
  }

  # Exactly one traffic destination (spec-validated): the Standard load
  # balancer frontends the service fronts, or one fixed private IP.
  load_balancer_frontend_ip_configuration_ids = (
    length(var.spec.load_balancer_frontend_ip_configuration_ids) > 0
    ? var.spec.load_balancer_frontend_ip_configuration_ids
    : null
  )
  destination_ip_address = var.spec.destination_ip_address != "" ? var.spec.destination_ip_address : null

  # PROXY protocol v2 headers give the backend the consumer's original
  # source IP -- only when the backend parses them.
  proxy_protocol_enabled = var.spec.proxy_protocol_enabled

  # Discoverability and approval: visibility gates who can SEE the
  # service (UUIDs, or the "*" wildcard); auto-approval skips the manual
  # connection-approval queue for listed subscriptions.
  auto_approval_subscription_ids = length(var.spec.auto_approval_subscription_ids) > 0 ? var.spec.auto_approval_subscription_ids : null
  visibility_subscription_ids    = length(var.spec.visibility_subscription_ids) > 0 ? var.spec.visibility_subscription_ids : null

  fqdns = length(var.spec.fqdns) > 0 ? var.spec.fqdns : null

  tags = local.final_tags
}
