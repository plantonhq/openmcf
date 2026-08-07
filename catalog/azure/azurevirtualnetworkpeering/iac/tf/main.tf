# Create one direction of a virtual network peering -- private connectivity
# between two networks over the Microsoft backbone.
#
# Lifecycle notes worth knowing before operating this resource:
# - This resource is ONE DIRECTION of a peering; traffic only flows once
#   the reciprocal peering exists on the remote network. Azure retries
#   internally while the far side catches up, so the two directions can
#   deploy concurrently.
# - The access/forwarding/gateway flags and the subnet-name lists update
#   IN PLACE. Name, the two networks, and the complete-vs-subnet-scoped
#   and IPv6-only choices are the peering's identity -- changing any of
#   them replaces it (a brief connectivity gap for this direction only).
# - Peerings are not tracked ARM resources, so they carry no tags.
# - The local network's name and resource group are derived from the
#   referenced network's ARM ID (see locals), so the module never asks the
#   user to restate derivable state that could then disagree with the
#   network.
resource "azurerm_virtual_network_peering" "main" {
  name                      = var.spec.name
  resource_group_name       = local.resource_group_name
  virtual_network_name      = local.virtual_network_name
  remote_virtual_network_id = var.spec.remote_virtual_network_id

  # The four connectivity flags; defaults mirror Azure's (access on;
  # forwarding, gateway transit, and remote gateways off), so both engines
  # always send the same effective values.
  allow_virtual_network_access = var.spec.allow_virtual_network_access
  allow_forwarded_traffic      = var.spec.allow_forwarded_traffic
  allow_gateway_transit        = var.spec.allow_gateway_transit
  use_remote_gateways          = var.spec.use_remote_gateways

  peer_complete_virtual_networks_enabled = var.spec.peer_complete_virtual_networks_enabled

  # Subnet-scoped peering: the subnet-name lists are only meaningful when
  # complete-network peering is off (spec-level validation enforces the
  # pairing); an empty list is simply not sent.
  local_subnet_names  = length(var.spec.local_subnet_names) > 0 ? var.spec.local_subnet_names : null
  remote_subnet_names = length(var.spec.remote_subnet_names) > 0 ? var.spec.remote_subnet_names : null

  # Only sent when explicitly chosen: ARM treats the field as a
  # creation-time property of the peering.
  only_ipv6_peering_enabled = var.spec.only_ipv6_peering_enabled ? true : null
}
