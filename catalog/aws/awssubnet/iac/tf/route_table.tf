# A subnet-owned route table is created when inline routes are supplied
# and/or when VGW route propagation is requested (propagation without inline
# routes is a legal shape: the table's contents then come entirely from the
# gateway). Each route maps its destination and its typed target onto the
# matching AWS route attributes (others left null).
resource "aws_route_table" "this" {
  count  = local.owns_route_table ? 1 : 0
  vpc_id = var.spec.vpc_id

  # Site-to-Site VPN / Direct Connect routes propagate in from these virtual
  # private gateways instead of being declared route by route.
  propagating_vgws = length(var.spec.propagating_vgws) > 0 ? var.spec.propagating_vgws : null

  dynamic "route" {
    for_each = var.spec.routes
    content {
      cidr_block                 = route.value.destination_cidr_block != "" ? route.value.destination_cidr_block : null
      ipv6_cidr_block            = route.value.destination_ipv6_cidr_block != "" ? route.value.destination_ipv6_cidr_block : null
      destination_prefix_list_id = route.value.destination_prefix_list_id != "" ? route.value.destination_prefix_list_id : null

      gateway_id                = route.value.target_type == "internet_gateway" ? route.value.target_id : null
      nat_gateway_id            = route.value.target_type == "nat_gateway" ? route.value.target_id : null
      transit_gateway_id        = route.value.target_type == "transit_gateway" ? route.value.target_id : null
      vpc_peering_connection_id = route.value.target_type == "vpc_peering_connection" ? route.value.target_id : null
      vpc_endpoint_id           = route.value.target_type == "vpc_endpoint" ? route.value.target_id : null
      network_interface_id      = route.value.target_type == "network_interface" ? route.value.target_id : null
      egress_only_gateway_id    = route.value.target_type == "egress_only_internet_gateway" ? route.value.target_id : null
      carrier_gateway_id        = route.value.target_type == "carrier_gateway" ? route.value.target_id : null
      # core_network and odb_network targets are ARNs, not ids. For ODB routes
      # AWS reports a gateway id alongside the ARN; providers before 6.53.0
      # surfaced that as perpetual drift -- the pinned provider handles it.
      core_network_arn = route.value.target_type == "core_network" ? route.value.target_id : null
      odb_network_arn  = route.value.target_type == "odb_network" ? route.value.target_id : null
      local_gateway_id = route.value.target_type == "local_gateway" ? route.value.target_id : null
    }
  }

  tags = local.aws_tags
}

# Associate the subnet with its inline-created table, or with the externally
# referenced route_table_id. When neither is set, the subnet stays on the VPC
# main route table and no association is created.
resource "aws_route_table_association" "this" {
  count = (local.owns_route_table || var.spec.route_table_id != "") ? 1 : 0

  subnet_id = aws_subnet.this.id
  route_table_id = (
    local.owns_route_table
    ? aws_route_table.this[0].id
    : var.spec.route_table_id
  )
}
