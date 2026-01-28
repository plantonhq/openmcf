##############################
# Public Subnets
##############################

resource "aws_subnet" "public" {
  for_each          = local.public_subnets
  vpc_id            = aws_vpc.this.id
  cidr_block        = each.value.cidr_block
  availability_zone = each.value.availability_zone
  map_public_ip_on_launch = true

  # Use metadata.labels instead of metadata.tags
  tags = merge(
      var.metadata.labels != null ? var.metadata.labels : {},
    { "Name" = each.key }
  )
}

##############################
# Public Route Table
##############################

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id

  # Use metadata.labels instead of metadata.tags
  tags = merge(
      var.metadata.labels != null ? var.metadata.labels : {},
    { "Name" = "${var.metadata.name}-public-RT" }
  )

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
  }
}

##############################
# Public Route Table Associations
##############################

resource "aws_route_table_association" "public" {
  for_each = aws_subnet.public

  route_table_id = aws_route_table.public.id
  subnet_id      = each.value.id
}

##############################
# Private Subnets
##############################

resource "aws_subnet" "private" {
  for_each   = local.private_subnets
  vpc_id     = aws_vpc.this.id
  cidr_block = each.value.cidr_block
  availability_zone = each.value.availability_zone

  # Use metadata.labels instead of metadata.tags
  tags = merge(
      var.metadata.labels != null ? var.metadata.labels : {},
    { "Name" = each.key }
  )
}

#################################################
# Private Route Tables with route to NAT gateway
#################################################

resource "aws_route_table" "private" {
  for_each = aws_subnet.private

  vpc_id = aws_vpc.this.id

  # Use metadata.labels instead of metadata.tags
  tags = merge(
      var.metadata.labels != null ? var.metadata.labels : {},
    { "Name" = "${each.key}-RT" }
  )

  # Conditionally add the NAT route only if NAT is enabled
  dynamic "route" {
    for_each = var.spec.is_nat_gateway_enabled ? [true] : []
    content {
      cidr_block = "0.0.0.0/0"
      nat_gateway_id = lookup(local.az_to_nat_gw, each.value.availability_zone, null)
    }
  }
}

##############################
# Private Route Table Associations
##############################

resource "aws_route_table_association" "private" {
  for_each = aws_subnet.private

  route_table_id = aws_route_table.private[each.key].id
  subnet_id      = each.value.id
}
