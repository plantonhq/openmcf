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
    "resource_kind" = "azure_firewall"
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

  # Enum wire maps. tfvars carries FULL proto enum value names; the maps
  # translate them to azurerm's exact (case-sensitive) vocabulary. Absent
  # sku fields fall back to the production defaults, sent explicitly so
  # both engines produce identical payloads.
  sku_name_wire = {
    "AZFW_VNET" = "AZFW_VNet"
    "AZFW_HUB"  = "AZFW_Hub"
  }
  sku_name = lookup(local.sku_name_wire, coalesce(var.spec.sku_name, "AZFW_VNET"), "AZFW_VNet")

  sku_tier_wire = {
    "BASIC"    = "Basic"
    "STANDARD" = "Standard"
    "PREMIUM"  = "Premium"
  }
  sku_tier = lookup(local.sku_tier_wire, coalesce(var.spec.sku_tier, "STANDARD"), "Standard")

  # Sent only when specified -- the ARM field is server-defaulted (Alert)
  # and Computed in the provider, so omission lets Azure own the default.
  threat_intel_wire = {
    "ALERT" = "Alert"
    "DENY"  = "Deny"
    "OFF"   = "Off"
  }
  threat_intel_mode = (
    var.spec.threat_intel_mode != null
  ) ? lookup(local.threat_intel_wire, var.spec.threat_intel_mode, null) : null
}
