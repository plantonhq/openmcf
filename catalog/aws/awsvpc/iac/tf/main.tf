resource "aws_vpc" "this" {
  # The IPAM pool ids are StringValueOrRef spec fields, pre-resolved to plain
  # strings by the orchestrator before the module runs (no IPAM pool catalog
  # kind exists yet, so today they carry literal pool ids).
  cidr_block          = var.spec.cidr_block != "" ? var.spec.cidr_block : null
  ipv4_ipam_pool_id   = var.spec.ipv4_ipam_pool_id != "" ? var.spec.ipv4_ipam_pool_id : null
  ipv4_netmask_length = var.spec.ipv4_netmask_length != 0 ? var.spec.ipv4_netmask_length : null

  instance_tenancy                     = var.spec.instance_tenancy != "" ? var.spec.instance_tenancy : null
  enable_dns_support                   = var.spec.enable_dns_support
  enable_dns_hostnames                 = var.spec.enable_dns_hostnames
  enable_network_address_usage_metrics = var.spec.enable_network_address_usage_metrics

  assign_generated_ipv6_cidr_block     = var.spec.assign_generated_ipv6_cidr_block
  ipv6_cidr_block                      = var.spec.ipv6_cidr_block != "" ? var.spec.ipv6_cidr_block : null
  ipv6_cidr_block_network_border_group = var.spec.ipv6_cidr_block_network_border_group != "" ? var.spec.ipv6_cidr_block_network_border_group : null
  ipv6_ipam_pool_id                    = var.spec.ipv6_ipam_pool_id != "" ? var.spec.ipv6_ipam_pool_id : null
  ipv6_netmask_length                  = var.spec.ipv6_netmask_length != 0 ? var.spec.ipv6_netmask_length : null

  tags = local.aws_tags
}

# Each secondary IPv4 CIDR is its own association so it can be added or removed
# without recreating the VPC. Literal-CIDR entries key by their CIDR (stable
# across list reordering); IPAM-sized entries key by list position.
resource "aws_vpc_ipv4_cidr_block_association" "secondary" {
  for_each = {
    for idx, entry in var.spec.secondary_ipv4_cidrs :
    (try(entry.cidr_block, "") != "" ? entry.cidr_block : "ipam-${idx}") => entry
  }

  vpc_id = aws_vpc.this.id

  cidr_block          = try(each.value.cidr_block, "") != "" ? each.value.cidr_block : null
  ipv4_ipam_pool_id   = try(each.value.ipam_pool_id, "") != "" ? each.value.ipam_pool_id : null
  ipv4_netmask_length = try(each.value.netmask_length, 0) != 0 ? each.value.netmask_length : null
}

# Each secondary IPv6 CIDR is its own association. Exactly one source per
# entry (spec CEL): an Amazon-provided block, a BYOIP public pool, or an IPAM
# pool. Pinned-CIDR entries key by their CIDR; the rest key by list position.
resource "aws_vpc_ipv6_cidr_block_association" "secondary" {
  for_each = {
    for idx, entry in var.spec.secondary_ipv6_cidrs :
    (try(entry.cidr_block, "") != "" ? entry.cidr_block : "ipv6-${idx}") => entry
  }

  vpc_id = aws_vpc.this.id

  assign_generated_ipv6_cidr_block = try(each.value.assign_generated, false) ? true : null
  ipv6_pool                        = try(each.value.ipv6_pool, "") != "" ? each.value.ipv6_pool : null
  ipv6_ipam_pool_id                = try(each.value.ipam_pool_id, "") != "" ? each.value.ipam_pool_id : null
  ipv6_cidr_block                  = try(each.value.cidr_block, "") != "" ? each.value.cidr_block : null
  ipv6_netmask_length              = try(each.value.netmask_length, 0) != 0 ? each.value.netmask_length : null
}

# VPC Encryption Control: AWS's VPC-wide monitor/enforce switch for
# encryption in transit. Rendered only when configured; exclusions are sent
# enable/disable per service and only apply in enforce mode (spec CEL keeps
# monitor-mode exclusions out).
resource "aws_vpc_encryption_control" "this" {
  count = var.spec.encryption_control != null ? 1 : 0

  vpc_id = aws_vpc.this.id
  mode   = var.spec.encryption_control.mode

  internet_gateway_exclusion             = try(var.spec.encryption_control.exclude_internet_gateway, false) ? "enable" : "disable"
  egress_only_internet_gateway_exclusion = try(var.spec.encryption_control.exclude_egress_only_internet_gateway, false) ? "enable" : "disable"
  nat_gateway_exclusion                  = try(var.spec.encryption_control.exclude_nat_gateway, false) ? "enable" : "disable"
  virtual_private_gateway_exclusion      = try(var.spec.encryption_control.exclude_virtual_private_gateway, false) ? "enable" : "disable"
  vpc_peering_exclusion                  = try(var.spec.encryption_control.exclude_vpc_peering, false) ? "enable" : "disable"
  vpc_lattice_exclusion                  = try(var.spec.encryption_control.exclude_vpc_lattice, false) ? "enable" : "disable"
  lambda_exclusion                       = try(var.spec.encryption_control.exclude_lambda, false) ? "enable" : "disable"
  elastic_file_system_exclusion          = try(var.spec.encryption_control.exclude_elastic_file_system, false) ? "enable" : "disable"

  tags = local.aws_tags
}
