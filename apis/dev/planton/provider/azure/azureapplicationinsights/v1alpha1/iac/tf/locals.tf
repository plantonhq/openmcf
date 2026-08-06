locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Base tags for Azure resources
  base_tags = {
    # PARITY-EXCEPTION: resource_kind here is the family-wide snake-case
    # literal and resource_id falls back to metadata.name, while the
    # Pulumi module emits the lowered CloudResourceKind enum string and
    # omits resource_id when metadata.id is empty. Output-neutral (tags
    # never feed stack outputs); aligning the two shapes is a family-wide
    # convention change, not a per-kind fix.
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_application_insights"
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

  # Merge base, org, environment, and user tags -- user tags win on key
  # conflicts (the governance surface belongs to the user).
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # Application-type wire values. Azure's API strings are CASE-SENSITIVE
  # and irregular ("Node.JS", "MobileCenter") -- an unmatched value would be
  # silently treated as ASP.NET by Azure, which is why the spec closes the
  # vocabulary and this map carries the exact wire strings.
  application_type_map = {
    "WEB"           = "web"
    "JAVA"          = "java"
    "NODE_JS"       = "Node.JS"
    "OTHER"         = "other"
    "IOS"           = "ios"
    "PHONE"         = "phone"
    "STORE"         = "store"
    "MOBILE_CENTER" = "MobileCenter"
  }
  application_type = (
    var.spec.application_type != null && var.spec.application_type != ""
    ? local.application_type_map[var.spec.application_type]
    : "web"
  )
}
