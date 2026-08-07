locals {
  # The subnet is an ARM child of the virtual network: the network's ARM ID
  # (/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/{name})
  # carries the network name and resource group, and the module derives both
  # rather than modeling redundant fields that could contradict the
  # referenced network. regex() fails the plan loudly on a malformed ID --
  # better than sending a half-parsed request to ARM.
  virtual_network_id_parts = regex(
    "(?i)/subscriptions/[^/]+/resourceGroups/(?P<resource_group>[^/]+)/providers/Microsoft\\.Network/virtualNetworks/(?P<network_name>[^/]+)",
    var.spec.virtual_network_id
  )
  resource_group_name  = local.virtual_network_id_parts.resource_group
  virtual_network_name = local.virtual_network_id_parts.network_name

  # Map the spec enum's name string to ARM's PrivateEndpointNetworkPolicies
  # value. null lets Azure apply its default (Disabled); only an explicit
  # mode is ever sent, so an unspecified spec and Azure's default deploy
  # identically on both engines.
  private_endpoint_network_policies = (
    var.spec.private_endpoint_network_policies == "ENABLED" ? "Enabled" :
    var.spec.private_endpoint_network_policies == "NETWORK_SECURITY_GROUP_ENABLED" ? "NetworkSecurityGroupEnabled" :
    var.spec.private_endpoint_network_policies == "ROUTE_TABLE_ENABLED" ? "RouteTableEnabled" : null
  )

  # Map the spec enum's name string to ARM's SharingScope value. ARM only
  # accepts it alongside disabled default outbound access (spec-level
  # validation enforces the pairing).
  sharing_scope = var.spec.sharing_scope == "TENANT" ? "Tenant" : null
}
