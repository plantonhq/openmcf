locals {
  ###########################################################################
  # Common resource metadata logic (consistent with other modules)
  ###########################################################################

  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Base labels
  base_labels = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "aws_vpc"
  }

  # Organization label only if var.metadata.org is non-empty
  org_label = (
  var.metadata.org != null && var.metadata.org != ""
  ) ? {
    "organization" = var.metadata.org
  } : {}

  # Environment label only if var.metadata.env is non-empty
  env_label = (
  var.metadata.env != null && var.metadata.env != ""
  ) ? {
    "organization" = var.metadata.env
  } : {}
  # Merge base, org, and environment labels into final_labels
  final_labels = merge(local.base_labels, local.org_label, local.env_label)


  ###########################################################################
  # VPC-specific logic
  ###########################################################################

  # List of availability zones from var.spec
  azs      = var.spec.availability_zones
  az_count = length(local.azs)

  # Public subnets: generate one subnet per availability zone index
  public_subnets_list = flatten([
    for az_i, az in local.azs : [
      for i in range(var.spec.subnets_per_availability_zone) : {
        key = "public-${az}-${i}"
        value = {
          availability_zone = az
          cidr_block = cidrsubnet(
            var.spec.vpc_cidr,
            var.spec.subnet_size,
            (az_i * var.spec.subnets_per_availability_zone) + i
          )
        }
      }
    ]
  ])

  public_subnets = {
    for s in local.public_subnets_list : s.key => s.value
  }

  # Private subnets: offset the subnet index by the total number of public subnets
  private_subnets_list = flatten([
    for az_i, az in local.azs : [
      for i in range(var.spec.subnets_per_availability_zone) : {
        key = "private-${az}-${i}"
        value = {
          availability_zone = az
          cidr_block = cidrsubnet(
            var.spec.vpc_cidr,
            var.spec.subnet_size,
            (az_i * var.spec.subnets_per_availability_zone) + i + (
            var.spec.subnets_per_availability_zone * local.az_count
            )
          )
        }
      }
    ]
  ])

  private_subnets = {
    for s in local.private_subnets_list : s.key => s.value
  }

  # NAT Gateway subnets: pick the first public subnet in each AZ
  nat_gateway_subnets = {
    for az_i, az in local.azs :
    az => "public-${az}-0"
  }

  # Mapping from each AZ to the NAT Gateway created in that AZ
  az_to_nat_gw = {
    for az, subnet_key in local.nat_gateway_subnets :
    az => (var.spec.is_nat_gateway_enabled ? aws_nat_gateway.this[az].id : null)
  }
}
