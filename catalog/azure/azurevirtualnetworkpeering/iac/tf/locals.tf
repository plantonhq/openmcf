locals {
  # The peering is an ARM child of its LOCAL network: the network's ARM ID
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
}
