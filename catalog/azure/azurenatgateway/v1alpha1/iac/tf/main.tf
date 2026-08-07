# Create the NAT gateway -- the managed SNAT service that gives every
# workload in its attached subnets stable, scalable outbound connectivity.
#
# Lifecycle notes worth knowing before operating this resource:
# - Idle timeout, tags, and the IP/prefix associations update IN PLACE.
#   Name, SKU, and zone are the gateway's identity -- changing any of them
#   replaces it, briefly interrupting egress for every attached subnet.
# - The gateway is just the gateway: its addresses are referenced
#   first-class resources, and the subnets it serves attach themselves
#   (AzureSubnet's nat_gateway_id, matching Azure's model). A gateway with
#   no associated addresses deploys but cannot translate anything.
resource "azurerm_nat_gateway" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # null lets Azure apply its default SKU (Standard); mapped in locals.
  sku_name = local.sku_name

  idle_timeout_in_minutes = var.spec.idle_timeout_in_minutes

  # A STANDARD gateway is zonal; STANDARD_V2 is zone-redundant and forbids
  # zones (spec-level validation enforces the pairing).
  zones = length(var.spec.zones) > 0 ? var.spec.zones : null

  tags = local.final_tags
}

# Associate the referenced addresses and ranges. Each association is its
# own ARM operation; creating them here (rather than inside the public IP
# modules) keeps the addresses reusable and makes the gateway the owner of
# which addresses it SNATs through.
resource "azurerm_nat_gateway_public_ip_association" "main" {
  count = length(var.spec.public_ip_ids)

  nat_gateway_id       = azurerm_nat_gateway.main.id
  public_ip_address_id = var.spec.public_ip_ids[count.index]
}

resource "azurerm_nat_gateway_public_ip_prefix_association" "main" {
  count = length(var.spec.public_ip_prefix_ids)

  nat_gateway_id      = azurerm_nat_gateway.main.id
  public_ip_prefix_id = var.spec.public_ip_prefix_ids[count.index]
}
