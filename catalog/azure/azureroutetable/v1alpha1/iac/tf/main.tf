# Create the route table -- a reusable set of user-defined routes (UDRs)
# that overrides Azure's default system routing for every subnet attached
# to it.
#
# Lifecycle notes worth knowing before operating this resource:
# - Routes, BGP propagation, and tags all update IN PLACE -- and take
#   effect immediately for EVERY subnet attached to the table. Name,
#   region, and resource group are the table's ARM identity; changing any
#   of them replaces the table, detaching it from every subnet until the
#   replacement is re-attached.
# - The subnet-side attachment is deliberately not modeled here: a subnet
#   declares which route table it uses (matching Azure's model), so one
#   table serves many subnets without listing them.
resource "azurerm_route_table" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Routes are managed inline as part of the table (a route has no life of
  # its own in Azure). An empty list is meaningful: it removes the last
  # route rather than leaving existing routes externally managed.
  route = [
    for route in var.spec.routes : {
      name           = route.name
      address_prefix = route.address_prefix
      next_hop_type  = local.next_hop_type_to_arm[route.next_hop_type]
      # Only VirtualAppliance routes carry a forwarding IP (spec-level
      # validation enforces the pairing); ARM rejects it on any other hop
      # type.
      next_hop_in_ip_address = route.next_hop_in_ip_address
    }
  ]

  # Azure defaults to propagating BGP-learned routes; disabling is the
  # forced-tunneling hardening that keeps learned routes from bypassing
  # the user-defined ones.
  bgp_route_propagation_enabled = var.spec.bgp_route_propagation_enabled

  tags = local.final_tags
}
