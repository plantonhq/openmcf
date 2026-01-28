locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Generate NAT Gateway name
  nat_gateway_name = "natgw-${var.metadata.name}"

  # Parse subnet ID to extract resource group and location
  # Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/{vnet}/subnets/{subnet}
  subnet_id_parts = split("/", var.spec.subnet_id)

  # Extract resource group name from subnet ID
  resource_group = (
    length(local.subnet_id_parts) >= 5 ?
    element([for i, v in local.subnet_id_parts : local.subnet_id_parts[i + 1] if v == "resourceGroups"], 0) :
    ""
  )

  # Base tags for Azure resources
  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_nat_gateway"
    "resource_name" = var.metadata.name
  }

  # Organization tag only if var.metadata.org is non-empty
  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  # Environment tag only if var.metadata.env is non-empty
  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Merge base, org, environment, and user-provided tags
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # Determine if we're creating a prefix or individual IP
  use_ip_prefix = var.spec.public_ip_prefix_length != null && var.spec.public_ip_prefix_length > 0
}

# Data source to get the subnet details for location
data "azurerm_subnet" "target" {
  name                 = element(reverse(split("/", var.spec.subnet_id)), 0)
  virtual_network_name = element(reverse(split("/", var.spec.subnet_id)), 2)
  resource_group_name  = local.resource_group
}

