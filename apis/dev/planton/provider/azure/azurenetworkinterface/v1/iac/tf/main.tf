# Create the network interface -- the attachment point that gives a
# virtual machine its presence in a subnet.
#
# Lifecycle notes worth knowing before operating this resource:
# - Name, region, and edge zone are the NIC's identity -- changing any of
#   them replaces the NIC, detaching it from its VM. Everything else
#   (configurations, DNS, acceleration, forwarding, associations, tags)
#   updates in place.
# - The MAC address is assigned when the NIC attaches to a running VM,
#   not at creation.
# - NSG and ASG memberships are separate ARM operations, realized below
#   as association resources (Azure's own model) rather than inline NIC
#   properties -- detaching is just removing the spec field.
resource "azurerm_network_interface" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Most NICs carry exactly one configuration; multiple serve dual-stack
  # and multi-IP scenarios. ARM requires the first to be primary when
  # several are declared (spec-level validation enforces it).
  dynamic "ip_configuration" {
    for_each = local.ip_configurations
    content {
      name                                               = ip_configuration.value.name
      subnet_id                                          = ip_configuration.value.subnet_id
      private_ip_address_allocation                      = ip_configuration.value.allocation
      private_ip_address                                 = ip_configuration.value.private_ip_address
      private_ip_address_version                         = ip_configuration.value.version
      public_ip_address_id                               = ip_configuration.value.public_ip_address_id
      primary                                            = ip_configuration.value.primary
      gateway_load_balancer_frontend_ip_configuration_id = ip_configuration.value.gateway_lb_fip_id
    }
  }

  dns_servers             = var.spec.dns_servers
  internal_dns_name_label = var.spec.internal_dns_name_label

  accelerated_networking_enabled = var.spec.accelerated_networking_enabled
  ip_forwarding_enabled          = var.spec.ip_forwarding_enabled

  # NVA acceleration (preview; subscription must be enrolled). null sends
  # nothing -- the correct shape for every non-appliance NIC.
  auxiliary_mode = local.auxiliary_mode
  auxiliary_sku  = local.auxiliary_sku

  edge_zone = var.spec.edge_zone

  tags = local.final_tags
}

# Attach the NIC-level network security group. Its own ARM operation
# (Azure's model), so filtering can change without touching the NIC.
resource "azurerm_network_interface_security_group_association" "main" {
  count = var.spec.network_security_group_id != null && var.spec.network_security_group_id != "" ? 1 : 0

  network_interface_id      = azurerm_network_interface.main.id
  network_security_group_id = var.spec.network_security_group_id
}

# Join the NIC to its application security groups so NSG rules can target
# workload groups instead of IP ranges.
resource "azurerm_network_interface_application_security_group_association" "main" {
  count = length(var.spec.application_security_group_ids)

  network_interface_id          = azurerm_network_interface.main.id
  application_security_group_id = var.spec.application_security_group_ids[count.index]
}

# Join ip_configurations to load-balancer backend pools. Membership is
# expressed from the member side in Azure's model -- the pool is declared
# on the load balancer, and this association adds THIS NIC to it.
resource "azurerm_network_interface_backend_address_pool_association" "main" {
  for_each = local.lb_pool_associations

  network_interface_id    = azurerm_network_interface.main.id
  ip_configuration_name   = each.value.ip_configuration_name
  backend_address_pool_id = each.value.pool_id
}

# Complete single-target inbound NAT rules: the load balancer declares
# the port forward, this association picks the receiving instance.
resource "azurerm_network_interface_nat_rule_association" "main" {
  for_each = local.lb_nat_rule_associations

  network_interface_id  = azurerm_network_interface.main.id
  ip_configuration_name = each.value.ip_configuration_name
  nat_rule_id           = each.value.nat_rule_id
}

# Join ip_configurations to Application Gateway backend pools (the L7
# counterpart of the load-balancer membership above).
resource "azurerm_network_interface_application_gateway_backend_address_pool_association" "main" {
  for_each = local.appgw_pool_associations

  network_interface_id    = azurerm_network_interface.main.id
  ip_configuration_name   = each.value.ip_configuration_name
  backend_address_pool_id = each.value.pool_id
}
