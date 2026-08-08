# Create the ExpressRoute circuit -- the dedicated private connection
# between your infrastructure and Microsoft. The circuit is the billing
# and identity object: creation issues the service key the connectivity
# provider needs to provision the physical cross-connect, and Azure
# meters the circuit FROM THIS MOMENT, even while the provider side is
# unprovisioned.
resource "azurerm_express_route_circuit" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  sku {
    tier   = lookup(local.sku_tier_wire, var.spec.sku_tier, var.spec.sku_tier)
    family = lookup(local.sku_family_wire, var.spec.sku_family, var.spec.sku_family)
  }

  # Exactly one provisioning mode (spec-validated): the service-provider
  # trio, or the ExpressRoute Direct pair. ARM treats the two property
  # sets as mutually exclusive shapes of the same create call.
  service_provider_name = var.spec.service_provider_name != "" ? var.spec.service_provider_name : null
  peering_location      = var.spec.peering_location != "" ? var.spec.peering_location : null
  # ARM can grow a provisioned circuit's bandwidth in place but never
  # shrink it -- the provider forces replacement on a decrease.
  bandwidth_in_mbps = var.spec.bandwidth_in_mbps > 0 ? var.spec.bandwidth_in_mbps : null

  express_route_port_id = var.spec.express_route_port_id != "" ? var.spec.express_route_port_id : null
  bandwidth_in_gbps     = var.spec.bandwidth_in_gbps > 0 ? var.spec.bandwidth_in_gbps : null

  # Rate limiting maps to ARM's EnableDirectPortRateLimit -- meaningful
  # on ExpressRoute Direct circuits only.
  rate_limiting_enabled    = var.spec.rate_limiting_enabled
  allow_classic_operations = var.spec.allow_classic_operations

  # The key this circuit REDEEMS (capacity someone else owns) -- not the
  # keys it issues. ARM never returns it on reads; the provider writes it
  # in a follow-up call after the circuit exists.
  authorization_key = var.spec.authorization_key != "" ? var.spec.authorization_key : null

  tags = local.final_tags
}

# The composed authorizations: standalone ARM children of the circuit,
# one per spec entry, keyed by name. ARM GENERATES each key; the
# name-keyed authorization_keys output surfaces them (sensitive) so a
# far-side gateway in another subscription can redeem one.
resource "azurerm_express_route_circuit_authorization" "authorizations" {
  for_each = { for authorization in var.spec.authorizations : authorization.name => authorization }

  name                       = each.value.name
  express_route_circuit_name = azurerm_express_route_circuit.main.name
  resource_group_name        = var.spec.resource_group
}
