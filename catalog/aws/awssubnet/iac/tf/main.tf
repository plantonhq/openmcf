resource "aws_subnet" "this" {
  vpc_id = var.spec.vpc_id

  # Exactly one placement form is present (spec CEL): the zone name or the
  # account-stable zone id.
  availability_zone    = var.spec.availability_zone != "" ? var.spec.availability_zone : null
  availability_zone_id = var.spec.availability_zone_id != "" ? var.spec.availability_zone_id : null

  # IPv4 addressing: an explicit CIDR, an IPAM allocation, or neither
  # (ipv6-native subnets). The spec CEL guarantees the arms never mix.
  cidr_block          = var.spec.cidr_block != "" ? var.spec.cidr_block : null
  ipv4_ipam_pool_id   = var.spec.ipv4_ipam_pool_id != "" ? var.spec.ipv4_ipam_pool_id : null
  ipv4_netmask_length = var.spec.ipv4_netmask_length

  map_public_ip_on_launch                        = var.spec.map_public_ip_on_launch
  assign_ipv6_address_on_creation                = var.spec.assign_ipv6_address_on_creation
  enable_dns64                                   = var.spec.enable_dns64
  enable_resource_name_dns_a_record_on_launch    = var.spec.enable_resource_name_dns_a_record_on_launch
  enable_resource_name_dns_aaaa_record_on_launch = var.spec.enable_resource_name_dns_aaaa_record_on_launch

  # IPv6 addressing: an explicit /44-/64 CIDR or an IPAM allocation (spec CEL
  # keeps them exclusive); ipv6_native drops IPv4 entirely (IPv6-only subnet).
  ipv6_cidr_block     = var.spec.ipv6_cidr_block != "" ? var.spec.ipv6_cidr_block : null
  ipv6_ipam_pool_id   = var.spec.ipv6_ipam_pool_id != "" ? var.spec.ipv6_ipam_pool_id : null
  ipv6_netmask_length = var.spec.ipv6_netmask_length
  ipv6_native         = var.spec.ipv6_native

  private_dns_hostname_type_on_launch = var.spec.private_dns_hostname_type_on_launch != "" ? var.spec.private_dns_hostname_type_on_launch : null

  tags = local.aws_tags
}
