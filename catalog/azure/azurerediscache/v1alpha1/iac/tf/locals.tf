locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  base_tags = {
    # PARITY-EXCEPTION: resource_kind here is the family-wide snake-case
    # literal and resource_id falls back to metadata.name, while the
    # Pulumi module emits the lowered CloudResourceKind enum string and
    # omits resource_id when metadata.id is empty. Output-neutral (tags
    # never feed stack outputs); aligning the two shapes is a family-wide
    # convention change, not a per-kind fix.
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_redis_cache"
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

  # The spec's sku enum arrives as the FULL proto value name; absent means
  # STANDARD (the tfvars wire format drops zero-valued proto fields, so the
  # module materializes the spec's documented default).
  sku_map = {
    "BASIC"    = "Basic"
    "STANDARD" = "Standard"
    "PREMIUM"  = "Premium"
  }
  sku_name = local.sku_map[coalesce(var.spec.sku_name, "STANDARD")]

  # Azure spells the size as family letter + capacity number; the family
  # letter is fully determined by the tier ("C" for Basic/Standard, "P"
  # for Premium), so the spec never spells it twice.
  family = local.sku_name == "Premium" ? "P" : "C"

  # Patch-schedule days arrive as the spec enum's name string; ARM wants
  # the capitalized English day name.
  day_of_week_map = {
    "MONDAY"    = "Monday"
    "TUESDAY"   = "Tuesday"
    "WEDNESDAY" = "Wednesday"
    "THURSDAY"  = "Thursday"
    "FRIDAY"    = "Friday"
    "SATURDAY"  = "Saturday"
    "SUNDAY"    = "Sunday"
  }

  # Persistence auth method: spec enum name -> ARM's value. null when
  # unset so Azure applies its default (SAS).
  persistence_auth_map = {
    "SAS"              = "SAS"
    "MANAGED_IDENTITY" = "ManagedIdentity"
  }

  # Identity type: spec enum name -> ARM's value.
  identity_type_map = {
    "SYSTEM_ASSIGNED"          = "SystemAssigned"
    "USER_ASSIGNED"            = "UserAssigned"
    "SYSTEM_AND_USER_ASSIGNED" = "SystemAssigned, UserAssigned"
  }

  redis_configuration = var.spec.redis_configuration
}
