# Transit Gateway VPC attachment -- one spoke of the hub.
#
# AWS provisions an ENI in each listed subnet; traffic between the VPC and
# the gateway flows through those ENIs. The gateway and the VPC are
# create-time immutable (changing either replaces the attachment), while the
# subnet set updates in place -- AZs can be added or removed without
# replacement.
resource "aws_ec2_transit_gateway_vpc_attachment" "this" {
  # Referenced IDs arrive pre-resolved as plain strings.
  transit_gateway_id = var.spec.transit_gateway_id
  vpc_id             = var.spec.vpc_id
  subnet_ids         = var.spec.subnet_ids

  # Tri-state options: null falls through to the provider default ("enable"
  # for DNS support) or, for security group referencing, to the value the
  # attachment inherits from the gateway (AWS computes it).
  dns_support                        = local.dns_support
  security_group_referencing_support = local.security_group_referencing_support

  # Plain-bool options whose spec default (false) matches the AWS default
  # ("disable"): always sent explicitly. Appliance mode keeps a flow's
  # return traffic in the AZ it entered through -- required for stateful
  # inspection VPCs, harmless noise everywhere else.
  ipv6_support           = var.spec.ipv6_support ? "enable" : "disable"
  appliance_mode_support = var.spec.appliance_mode_support ? "enable" : "disable"

  # Default route table membership. Left null, AWS derives the value from
  # the GATEWAY's own default-association/propagation dials -- the provider
  # marks both Optional+Computed for exactly this inheritance. Only a spec
  # that pins them sends a value: false detaches this spoke from the default
  # table (the segmented-topology posture, where a custom
  # AwsTransitGatewayRouteTable owns the association instead).
  transit_gateway_default_route_table_association = var.spec.default_route_table_association
  transit_gateway_default_route_table_propagation = var.spec.default_route_table_propagation

  tags = local.aws_tags
}
