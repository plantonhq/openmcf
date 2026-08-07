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
    "resource_kind" = "azure_dns_record"
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

  # The platform materializes the proto default (300) before the module
  # runs; the coalesce is a same-value safety net for direct stack-input
  # paths, never a different fallback.
  ttl = coalesce(var.spec.ttl_seconds, 300)

  # The CAA tag arrives as the proto enum value name (ISSUE, ISSUEWILD,
  # IODEF, CONTACTEMAIL); Azure's wire vocabulary is the lowercase form
  # and the provider validates it case-sensitively.
  caa_tag_map = {
    "ISSUE"        = "issue"
    "ISSUEWILD"    = "issuewild"
    "IODEF"        = "iodef"
    "CONTACTEMAIL" = "contactemail"
  }
}
