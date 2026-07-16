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
    "resource_kind" = "azure_front_door_profile"
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
    "STANDARD" = "Standard_AzureFrontDoor"
    "PREMIUM"  = "Premium_AzureFrontDoor"
  }
  sku_name = local.sku_map[coalesce(var.spec.sku, "STANDARD")]

  # Identity type: spec enum name -> ARM's value.
  identity_type_map = {
    "SYSTEM_ASSIGNED"          = "SystemAssigned"
    "USER_ASSIGNED"            = "UserAssigned"
    "SYSTEM_AND_USER_ASSIGNED" = "SystemAssigned, UserAssigned"
  }

  # Log-scrubbing match variables: spec enum names -> ARM's values.
  log_scrubbing_variable_map = {
    "QUERY_STRING_ARG_NAMES" = "QueryStringArgNames"
    "REQUEST_IP_ADDRESS"     = "RequestIPAddress"
    "REQUEST_URI"            = "RequestUri"
  }
  log_scrubbing_variables = coalesce(var.spec.log_scrubbing_variables, [])
}
