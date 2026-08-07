locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # PARITY-EXCEPTION: resource_kind here is the family-wide snake-case
  # literal and resource_id falls back to metadata.name, while the Pulumi
  # module emits the lowered CloudResourceKind enum string and omits
  # resource_id when metadata.id is empty. Output-neutral (tags never feed
  # stack outputs); aligning the two shapes is a family-wide convention
  # change, not a per-kind fix.
  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_public_ip_prefix"
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

  # Map the spec enums' name strings to ARM's values. null lets Azure apply
  # its defaults (IPv4 / Standard / Regional); only an explicit choice is
  # ever sent, so an unspecified spec and Azure's default deploy identically
  # on both engines.
  ip_version = (
    var.spec.ip_version == "IPV4" ? "IPv4" :
    var.spec.ip_version == "IPV6" ? "IPv6" : null
  )
  sku = (
    var.spec.sku == "STANDARD" ? "Standard" :
    var.spec.sku == "STANDARD_V2" ? "StandardV2" : null
  )
  sku_tier = (
    var.spec.sku_tier == "REGIONAL" ? "Regional" :
    var.spec.sku_tier == "GLOBAL" ? "Global" : null
  )
}
