resource "aws_nat_gateway" "this" {
  connectivity_type = var.spec.connectivity_type

  # Placement: zonal gateways live in one subnet; regional gateways span a
  # VPC (availability_mode = "regional" + vpc_id). The spec's CEL rules
  # guarantee exactly one placement form arrives here.
  availability_mode = local.availability_mode
  subnet_id         = local.subnet_id
  vpc_id            = local.vpc_id

  # Public gateways: Elastic IP allocation(s) (zonal form).
  allocation_id            = local.allocation_id
  secondary_allocation_ids = local.secondary_allocation_ids

  # Private gateways: private IPv4 addressing (zonal form).
  private_ip                         = local.private_ip
  secondary_private_ip_addresses     = local.secondary_private_ip_addresses
  secondary_private_ip_address_count = local.secondary_private_ip_address_count

  # Regional gateways: explicit per-AZ layout (manual mode). An empty set
  # means AWS picks the zones and manages addresses itself (auto mode);
  # switching a live gateway between auto and manual replaces it (provider
  # behavior, via its internal auto-mode marker attribute).
  dynamic "availability_zone_address" {
    for_each = var.spec.availability_zone_addresses
    content {
      availability_zone    = availability_zone_address.value.availability_zone != "" ? availability_zone_address.value.availability_zone : null
      availability_zone_id = availability_zone_address.value.availability_zone_id != "" ? availability_zone_address.value.availability_zone_id : null
      allocation_ids       = length(availability_zone_address.value.allocation_ids) > 0 ? availability_zone_address.value.allocation_ids : null
    }
  }

  tags = local.aws_tags
}
