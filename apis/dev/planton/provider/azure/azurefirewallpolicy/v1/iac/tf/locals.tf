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
    "resource_kind" = "azure_firewall_policy"
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
  # translate them to azurerm's exact (case-sensitive) vocabulary. An
  # absent sku/mode falls back to Azure's own default, sent explicitly so
  # both engines produce identical payloads.
  sku_wire = {
    "BASIC"    = "Basic"
    "STANDARD" = "Standard"
    "PREMIUM"  = "Premium"
  }
  sku = lookup(local.sku_wire, coalesce(var.spec.sku, "STANDARD"), "Standard")

  threat_intel_wire = {
    "ALERT" = "Alert"
    "DENY"  = "Deny"
    "OFF"   = "Off"
  }
  threat_intelligence_mode = lookup(
    local.threat_intel_wire,
    coalesce(var.spec.threat_intelligence_mode, "ALERT"),
    "Alert"
  )

  # IDPS states carry an IDPS_ prefix in the proto (package scoping); the
  # wire values are Azure's plain Off/Alert/Deny.
  idps_state_wire = {
    "IDPS_OFF"   = "Off"
    "IDPS_ALERT" = "Alert"
    "IDPS_DENY"  = "Deny"
  }

  # Bypass protocols: the proto names ARE the wire vocabulary (the
  # provider validates ANY/ICMP/TCP/UDP case-insensitively and
  # case-suppresses ARM's mixed-case echo).
  identity_type_wire = {
    "SYSTEM_ASSIGNED"          = "SystemAssigned"
    "USER_ASSIGNED"            = "UserAssigned"
    "SYSTEM_AND_USER_ASSIGNED" = "SystemAssigned, UserAssigned"
  }
}
