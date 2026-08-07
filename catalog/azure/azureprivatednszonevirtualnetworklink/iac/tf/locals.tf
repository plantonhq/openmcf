locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_private_dns_zone_virtual_network_link"
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

  # The link is an ARM child of the zone: the zone's ARM ID
  # (/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/privateDnsZones/{zone})
  # carries the zone name and resource group, and the module derives both
  # rather than modeling redundant fields that could contradict the
  # referenced zone. regex() fails the plan loudly on a malformed ID --
  # better than sending a half-parsed request to ARM.
  zone_id_parts = regex(
    "(?i)/subscriptions/[^/]+/resourceGroups/(?P<resource_group>[^/]+)/providers/Microsoft\\.Network/privateDnsZones/(?P<zone_name>[^/]+)",
    var.spec.private_dns_zone_id
  )
  zone_resource_group_name = local.zone_id_parts.resource_group
  zone_name                = local.zone_id_parts.zone_name

  # Map the spec enum's name string to ARM's ResolutionPolicy value. null
  # lets Azure apply its per-zone-type default; only an explicit policy is
  # ever sent, so an unspecified spec and Azure's default deploy
  # identically on both engines.
  resolution_policy = (
    var.spec.resolution_policy == "DEFAULT" ? "Default" :
    var.spec.resolution_policy == "NX_DOMAIN_REDIRECT" ? "NxDomainRedirect" : null
  )
}
