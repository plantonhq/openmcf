# Transit Gateway route table -- one isolated routing domain.
#
# The table itself is a tiny resource; the routing domain it defines lives
# in the folded members below. Everything is keyed by stable identifiers
# (attachment ID, destination CIDR, prefix list ID) so membership changes
# surgically add or remove one provider resource instead of reshuffling the
# whole domain.
resource "aws_ec2_transit_gateway_route_table" "this" {
  # Create-time immutable: moving the table to another gateway replaces it
  # (and with it, every association, propagation, and route below).
  transit_gateway_id = var.spec.transit_gateway_id

  tags = local.aws_tags
}

# Associations: the attachments whose OUTBOUND traffic is looked up in this
# table. AWS allows an attachment at most ONE association across the whole
# gateway -- an attachment listed here must have default_route_table_
# association turned off, and must not appear in another table's
# associations. Documents cannot see each other, so that gateway-wide
# uniqueness is enforced by AWS at apply time, not by validation.
resource "aws_ec2_transit_gateway_route_table_association" "this" {
  for_each = local.associations

  transit_gateway_attachment_id  = each.value
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.this.id
}

# Propagations: the attachments that ADVERTISE their routes into this table.
# VPC attachments propagate their VPC CIDRs; VPN and Direct Connect
# attachments propagate BGP-learned routes. Unlike associations, an
# attachment can propagate to any number of tables.
resource "aws_ec2_transit_gateway_route_table_propagation" "this" {
  for_each = local.propagations

  transit_gateway_attachment_id  = each.value
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.this.id
}

# Static routes: longest-prefix match beats propagated routes on ties, so
# these steer specific prefixes (an inspection detour, a default route
# toward an egress VPC) or blackhole traffic that must never cross the hub.
# The spec CEL guarantees exactly one of attachment_id/blackhole per route;
# the provider expects the attachment argument to be ABSENT (not empty) for
# a blackhole, hence the null fallback.
resource "aws_ec2_transit_gateway_route" "this" {
  for_each = local.routes_map

  destination_cidr_block         = each.value.destination_cidr_block
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.this.id
  transit_gateway_attachment_id  = each.value.blackhole ? null : each.value.attachment_id
  blackhole                      = each.value.blackhole
}

# Prefix list references: route a managed prefix list's whole CIDR set via
# one attachment (or blackhole it), tracking the list's membership as it
# changes -- one reference instead of N hand-maintained static routes.
resource "aws_ec2_transit_gateway_prefix_list_reference" "this" {
  for_each = local.prefix_list_references_map

  prefix_list_id                 = each.value.prefix_list_id
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.this.id
  transit_gateway_attachment_id  = each.value.blackhole ? null : each.value.attachment_id
  blackhole                      = each.value.blackhole
}
