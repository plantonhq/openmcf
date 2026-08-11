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
    "resource_kind" = "azure_cognitive_account"
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

  # The spec's bypass enum names to the provider's wire values.
  # Unspecified (the enum's zero value renders as "") maps to null so the
  # property is omitted and ARM applies its default.
  network_acls_bypass_wire = {
    "AZURE_SERVICES" = "AzureServices"
    "NONE"           = "None"
  }

  # The spec's content-severity enum names to the provider's wire values.
  rai_content_level_wire = {
    "LOW"    = "Low"
    "MEDIUM" = "Medium"
    "HIGH"   = "High"
  }

  # The spec's responsible-AI policy modes to the provider's wire values.
  rai_policy_mode_wire = {
    "DEFAULT"             = "Default"
    "BLOCKING"            = "Blocking"
    "ASYNCHRONOUS_FILTER" = "AsynchronousFilter"
    "DEFERRED"            = "Deferred"
  }
}
