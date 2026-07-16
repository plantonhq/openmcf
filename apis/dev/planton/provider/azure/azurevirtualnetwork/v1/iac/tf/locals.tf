locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_virtual_network"
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

  # Map the spec enum's name string to ARM's enforcement value. null keeps
  # the encryption block absent entirely (ARM's default: encryption off),
  # so an unspecified spec and Azure's default deploy identically on both
  # engines.
  encryption_enforcement = (
    var.spec.encryption == "ALLOW_UNENCRYPTED" ? "AllowUnencrypted" :
    var.spec.encryption == "DROP_UNENCRYPTED" ? "DropUnencrypted" : null
  )

  # Map the spec enum's name string to ARM's policy value. null lets
  # azurerm apply ARM's default ("Disabled"); only the opt-in "Basic" mode
  # is ever sent.
  private_endpoint_vnet_policies = (
    var.spec.private_endpoint_vnet_policies == "BASIC" ? "Basic" : null
  )
}
