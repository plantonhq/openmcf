package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpSubnetworkSelfLink = "subnetwork_self_link"
	OpSubnetworkName     = "subnetwork_name"
	OpRegion             = "region"
	OpIpCidrRange        = "ip_cidr_range"
	OpSecondaryRanges    = "secondary_ranges"
	OpGatewayAddress     = "gateway_address"
	OpSubnetworkId       = "subnetwork_id"
	OpInternalIpv6Prefix = "internal_ipv6_prefix"
	OpExternalIpv6Prefix = "external_ipv6_prefix"
)
