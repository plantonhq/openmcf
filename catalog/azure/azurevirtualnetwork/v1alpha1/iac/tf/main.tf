# Create the virtual network -- the isolated private address space every
# network-attached Azure workload lives inside.
#
# Lifecycle notes worth knowing before operating this resource:
# - Address space, DNS servers, BGP community, DDoS attachment, encryption,
#   flow timeout, and tags all update IN PLACE. Name, region, resource
#   group, and edge zone are the network's ARM identity -- changing any of
#   them replaces the network and everything inside it.
# - Address-space blocks can be added or removed live, but a block that
#   subnets are carved from cannot shrink below them.
# - The network is deliberately just the network: subnets live in
#   AzureSubnet, outbound NAT in AzureNatGateway, and private DNS
#   attachments in AzurePrivateDnsZoneVirtualNetworkLink, each referencing
#   this network's outputs.
resource "azurerm_virtual_network" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Exactly one address source is set (spec-level validation): either
  # self-managed CIDR blocks or delegated allocation from Azure Network
  # Manager IPAM pools. azurerm treats an empty set as "not configured", so
  # passing both fields through directly preserves the exactly-one contract.
  address_space = length(var.spec.address_spaces) > 0 ? var.spec.address_spaces : null

  # With IPAM pools the actual CIDR ranges are provisioned at deploy time
  # and surface through the address_space output attribute.
  dynamic "ip_address_pool" {
    for_each = var.spec.ip_address_pools
    content {
      id                     = ip_address_pool.value.id
      number_of_ip_addresses = ip_address_pool.value.number_of_ip_addresses
    }
  }

  # Empty means Azure's default resolver (168.63.129.16) serves the network
  # -- required for private DNS zone resolution to work directly, so custom
  # servers are only sent when explicitly configured.
  dns_servers = length(var.spec.dns_servers) > 0 ? var.spec.dns_servers : null

  bgp_community = var.spec.bgp_community

  # The DDoS plan is a separate, shared (and billed) resource; this block
  # only attaches an existing plan. ARM keeps attachment and activation
  # distinct so a plan can stay attached with protection toggled off.
  dynamic "ddos_protection_plan" {
    for_each = var.spec.ddos_protection_plan != null ? [var.spec.ddos_protection_plan] : []
    content {
      id     = ddos_protection_plan.value.id
      enable = ddos_protection_plan.value.enable
    }
  }

  # Absent means ARM's default (encryption off); only an explicit
  # enforcement mode sends the block, so an unspecified spec and Azure's
  # default deploy identically on both engines. Note ARM currently accepts
  # only AllowUnencrypted -- DropUnencrypted is modeled because the API
  # defines it but is not yet generally available.
  dynamic "encryption" {
    for_each = local.encryption_enforcement != null ? [local.encryption_enforcement] : []
    content {
      enforcement = encryption.value
    }
  }

  # null lets Azure apply its 4-minute default; the spec constrains the
  # value to ARM's accepted 4-30 range.
  flow_timeout_in_minutes = var.spec.flow_timeout_in_minutes

  # null means ARM's default ("Disabled"); only the opt-in "Basic" mode is
  # ever sent.
  private_endpoint_vnet_policies = local.private_endpoint_vnet_policies

  edge_zone = var.spec.edge_zone

  tags = local.final_tags
}
