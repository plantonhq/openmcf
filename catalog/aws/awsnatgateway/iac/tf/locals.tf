locals {
  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsNatGateway"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Optional inputs map to null when unset so AWS applies its own defaults and
  # the mode/connectivity-specific attributes are simply omitted.
  allocation_id                      = var.spec.allocation_id != "" ? var.spec.allocation_id : null
  private_ip                         = var.spec.private_ip != "" ? var.spec.private_ip : null
  secondary_allocation_ids           = length(var.spec.secondary_allocation_ids) > 0 ? var.spec.secondary_allocation_ids : null
  secondary_private_ip_addresses     = length(var.spec.secondary_private_ip_addresses) > 0 ? var.spec.secondary_private_ip_addresses : null
  secondary_private_ip_address_count = var.spec.secondary_private_ip_address_count > 0 ? var.spec.secondary_private_ip_address_count : null

  # Regional-mode inputs (spec CEL keeps them off zonal gateways). The empty
  # availability_mode means zonal -- the classic form -- and is sent as null
  # so pre-regional state never diffs.
  availability_mode = var.spec.availability_mode != "" ? var.spec.availability_mode : null
  subnet_id         = var.spec.subnet_id != "" ? var.spec.subnet_id : null
  vpc_id            = var.spec.vpc_id != "" ? var.spec.vpc_id : null
}
