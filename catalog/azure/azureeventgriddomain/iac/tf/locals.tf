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
    "resource_kind" = "azure_eventgrid_domain"
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

  # The spec enum's value names mapped to the provider's identity tokens.
  identity_type_map = {
    "SYSTEM_ASSIGNED" = "SystemAssigned"
    "USER_ASSIGNED"   = "UserAssigned"
  }

  # An input mapping block is sent only when it carries at least one
  # non-empty field -- the provider treats an empty block and an absent
  # one identically, and the built-in schemas need no mapping at all.
  input_mapping_fields_set = (
    var.spec.input_mapping_fields != null && length(compact([
      var.spec.input_mapping_fields.id,
      var.spec.input_mapping_fields.topic,
      var.spec.input_mapping_fields.event_time,
      var.spec.input_mapping_fields.event_type,
      var.spec.input_mapping_fields.subject,
      var.spec.input_mapping_fields.data_version,
    ])) > 0
  )

  input_mapping_default_values_set = (
    var.spec.input_mapping_default_values != null && length(compact([
      var.spec.input_mapping_default_values.event_type,
      var.spec.input_mapping_default_values.subject,
      var.spec.input_mapping_default_values.data_version,
    ])) > 0
  )
}
