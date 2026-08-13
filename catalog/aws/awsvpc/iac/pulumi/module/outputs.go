package module

const (
	OpVpcId                  = "vpc_id"
	OpVpcArn                 = "vpc_arn"
	OpCidrBlock              = "cidr_block"
	OpIpv6CidrBlock          = "ipv6_cidr_block"
	OpOwnerId                = "owner_id"
	OpMainRouteTableId       = "main_route_table_id"
	OpDefaultSecurityGroupId = "default_security_group_id"
	OpDefaultNetworkAclId    = "default_network_acl_id"
	OpDefaultRouteTableId    = "default_route_table_id"
	OpRegion                 = "region"
	// Keyed identically to the module's association resources (literal CIDR,
	// else ipam-<index> / ipv6-<index>) -- the import recipes resolve
	// per-instance association IDs through these keys.
	OpSecondaryIpv4CidrAssociationIds = "secondary_ipv4_cidr_association_ids"
	OpSecondaryIpv6CidrAssociationIds = "secondary_ipv6_cidr_association_ids"
	OpEncryptionControlId             = "encryption_control_id"
)
