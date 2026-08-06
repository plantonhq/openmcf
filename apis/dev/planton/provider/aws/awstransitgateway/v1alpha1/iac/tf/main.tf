# AWS Transit Gateway -- the pure hub.
#
# The gateway owns the BGP/routing defaults, the feature dials, and the
# default route table pair. Everything that composes onto it -- VPC
# attachments, custom route tables with their associations/propagations/
# routes -- is its own resource kind referencing this gateway's outputs, so
# spokes and routing domains come and go without touching the hub.
resource "aws_ec2_transit_gateway" "this" {
  # Empty string means "not set": send null so both engines omit the
  # attribute identically instead of pinning an empty description.
  description = var.spec.description == "" ? null : var.spec.description

  # 0 means "not set" (proto3 scalar zero value): fall through to the
  # provider default of 64512 rather than sending an out-of-range ASN.
  # Changing the ASN after creation replaces the gateway.
  amazon_side_asn = var.spec.amazon_side_asn == 0 ? null : var.spec.amazon_side_asn

  # Tri-state dials: null falls through to the provider default ("enable"
  # for DNS/ECMP and the default-table pair). AWS QUIRK on the default-table
  # pair: flipping disable -> enable REPLACES the gateway (the provider's
  # asymmetric ForceNew), while enable -> disable updates in place -- state
  # a topology's posture up front rather than tightening later.
  default_route_table_association = local.default_route_table_association
  default_route_table_propagation = local.default_route_table_propagation
  dns_support                     = local.dns_support
  vpn_ecmp_support                = local.vpn_ecmp_support

  # Left null unless the spec pins it: AWS computes the effective in-transit
  # encryption posture on its own when unset.
  encryption_support = local.encryption_support

  # Plain-bool dials whose spec default (false) matches the AWS default
  # ("disable"): always sent explicitly, so the applied state is exactly the
  # spec with no inheritance ambiguity. multicast_support is create-time
  # immutable -- changing it replaces the gateway.
  auto_accept_shared_attachments     = var.spec.auto_accept_shared_attachments ? "enable" : "disable"
  security_group_referencing_support = var.spec.security_group_referencing_support ? "enable" : "disable"
  multicast_support                  = var.spec.multicast_support ? "enable" : "disable"

  # Only needed for TGW Connect (GRE appliance integration); empty for the
  # overwhelming majority of hubs.
  transit_gateway_cidr_blocks = var.spec.transit_gateway_cidr_blocks

  tags = local.aws_tags
}
