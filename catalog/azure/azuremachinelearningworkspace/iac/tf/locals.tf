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
    "resource_kind" = "azure_machine_learning_workspace"
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

  # The spec's identity flavors to the provider's comma-joined wire values.
  identity_type_wire = {
    "SYSTEM_ASSIGNED"          = "SystemAssigned"
    "USER_ASSIGNED"            = "UserAssigned"
    "SYSTEM_AND_USER_ASSIGNED" = "SystemAssigned, UserAssigned"
  }

  # The spec's workspace flavors to the provider's wire values.
  # Unspecified (the enum's zero value renders as "") maps to null so the
  # provider applies its default, "Default".
  kind_wire = {
    "DEFAULT"       = "Default"
    "FEATURE_STORE" = "FeatureStore"
  }

  # The spec's isolation modes to the provider's wire values. Unspecified
  # maps to null so the property is omitted and the value is read back.
  isolation_mode_wire = {
    "DISABLED"                     = "Disabled"
    "ALLOW_INTERNET_OUTBOUND"      = "AllowInternetOutbound"
    "ALLOW_ONLY_APPROVED_OUTBOUND" = "AllowOnlyApprovedOutbound"
  }

  # The spec's storage-access types to the provider's wire values.
  # Unspecified maps to null so the provider applies its default,
  # "AccessKey".
  storage_account_access_type_wire = {
    "ACCESS_KEY" = "AccessKey"
    "IDENTITY"   = "Identity"
  }
}
