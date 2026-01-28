##############################
# one EIP per AZ for NAT gateway IP
##############################

resource "aws_eip" "nat" {
  for_each = var.spec.is_nat_gateway_enabled ? local.nat_gateway_subnets : {}
  domain = "vpc"

  # Use metadata.labels instead of metadata.tags
  tags = merge(
      var.metadata.labels != null ? var.metadata.labels : {},
    { "Name" = "${each.key}-nat-eip" }
  )
}

##############################
# NAT Gateway (one per AZ)
##############################

resource "aws_nat_gateway" "this" {
  for_each  = var.spec.is_nat_gateway_enabled ? local.nat_gateway_subnets : {}
  subnet_id = aws_subnet.public[each.value].id
  allocation_id = aws_eip.nat[each.key].id

  # Use metadata.labels instead of metadata.tags
  tags = merge(
      var.metadata.labels != null ? var.metadata.labels : {},
    { "Name" = "${each.key}-nat-gw" }
  )

  depends_on = [aws_subnet.public]
}
