locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_route_table"
    "resource_name" = var.metadata.name
  }

  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Metadata-derived tags first, then the user's spec tags merged over them:
  # user tags deliberately win so an org's governance conventions (cost
  # center, owner) can override the derived values where they collide.
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # Map the spec enum's name strings to ARM's RouteNextHopType values. The
  # spec enum requires a defined, non-unspecified value, so every route
  # resolves to exactly one ARM hop type.
  next_hop_type_to_arm = {
    "VIRTUAL_NETWORK_GATEWAY" = "VirtualNetworkGateway"
    "VNET_LOCAL"              = "VnetLocal"
    "INTERNET"                = "Internet"
    "VIRTUAL_APPLIANCE"       = "VirtualAppliance"
    "NONE"                    = "None"
  }
}
