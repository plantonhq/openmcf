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

  # Map the spec enum's name string to ARM's ResolutionPolicy value. null
  # lets Azure apply its per-zone-type default; only an explicit policy is
  # ever sent, so an unspecified spec and Azure's default deploy
  # identically on both engines.
  resolution_policy = (
    var.spec.resolution_policy == "DEFAULT" ? "Default" :
    var.spec.resolution_policy == "NX_DOMAIN_REDIRECT" ? "NxDomainRedirect" : null
  )
}
