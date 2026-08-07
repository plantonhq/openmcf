# Create the subnet -- the workload segment that partitions a virtual
# network's address space, and the composition hub where route tables,
# network security groups, and NAT gateways attach.
#
# Lifecycle notes worth knowing before operating this resource:
# - Address prefixes, service endpoints, policies, delegations, and the
#   route-table/NSG/NAT attachments all update IN PLACE; name and the
#   parent network are the subnet's ARM identity, so changing either
#   replaces the subnet and everything deployed into it.
# - A prefix in use by deployed resources cannot shrink, and toggling
#   default_outbound_access_enabled requires the subnet to be empty of
#   VMs -- plan address space and egress posture before workloads land.
# - Subnets are not tracked ARM resources, so they carry no tags.
resource "azurerm_subnet" "main" {
  name                 = var.spec.name
  resource_group_name  = local.resource_group_name
  virtual_network_name = local.virtual_network_name

  # Exactly one address source is set (spec-level validation enforces the
  # XOR): self-managed CIDR blocks, or delegated allocation from a Network
  # Manager IPAM pool that provisions the actual range at deploy time.
  address_prefixes = length(var.spec.address_prefixes) > 0 ? var.spec.address_prefixes : null

  dynamic "ip_address_pool" {
    for_each = var.spec.ip_address_pool != null ? [var.spec.ip_address_pool] : []
    content {
      id                     = ip_address_pool.value.id
      number_of_ip_addresses = ip_address_pool.value.number_of_ip_addresses
    }
  }

  service_endpoints           = length(var.spec.service_endpoints) > 0 ? var.spec.service_endpoints : null
  service_endpoint_policy_ids = length(var.spec.service_endpoint_policy_ids) > 0 ? var.spec.service_endpoint_policy_ids : null

  # Delegations hand the subnet to a PaaS service. An explicit action list
  # is only sent when the spec carries one; otherwise Azure applies the
  # service's default action set.
  dynamic "delegation" {
    for_each = var.spec.delegations
    content {
      name = delegation.value.name
      service_delegation {
        name    = delegation.value.service_name
        actions = length(delegation.value.actions) > 0 ? delegation.value.actions : null
      }
    }
  }

  private_endpoint_network_policies             = local.private_endpoint_network_policies
  private_link_service_network_policies_enabled = var.spec.private_link_service_network_policies_enabled
  default_outbound_access_enabled               = var.spec.default_outbound_access_enabled

  # ARM only accepts sharing_scope alongside disabled default outbound
  # access (spec-level validation enforces the pairing).
  sharing_scope = local.sharing_scope
}

# The attach seams. Azure models route-table/NSG/NAT attachment as writes
# to the subnet, declared subnet-side because one table, group, or gateway
# serves many subnets. Each association is its own ARM operation; creating
# them here (rather than inside the referenced resources' modules) keeps
# those resources reusable across subnets.

resource "azurerm_subnet_route_table_association" "main" {
  count = var.spec.route_table_id != null && var.spec.route_table_id != "" ? 1 : 0

  subnet_id      = azurerm_subnet.main.id
  route_table_id = var.spec.route_table_id
}

resource "azurerm_subnet_network_security_group_association" "main" {
  count = var.spec.network_security_group_id != null && var.spec.network_security_group_id != "" ? 1 : 0

  subnet_id                 = azurerm_subnet.main.id
  network_security_group_id = var.spec.network_security_group_id
}

resource "azurerm_subnet_nat_gateway_association" "main" {
  count = var.spec.nat_gateway_id != null && var.spec.nat_gateway_id != "" ? 1 : 0

  subnet_id      = azurerm_subnet.main.id
  nat_gateway_id = var.spec.nat_gateway_id
}
