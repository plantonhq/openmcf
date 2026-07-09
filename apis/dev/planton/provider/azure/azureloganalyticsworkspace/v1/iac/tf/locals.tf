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
    "resource_kind" = "azure_log_analytics_workspace"
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

  # SKU wire values. The tfvars wire format carries the FULL proto enum
  # value name; an absent sku deploys Azure's recommended PerGB2018.
  # Standard/Premium/LACluster/Unlimited are deliberately not mapped:
  # Azure blocks creating workspaces on them (see the spec enum comment).
  sku_map = {
    "PER_GB_2018"          = "PerGB2018"
    "CAPACITY_RESERVATION" = "CapacityReservation"
    "PER_NODE"             = "PerNode"
    "STANDALONE"           = "Standalone"
  }
  sku = (
    var.spec.sku != null && var.spec.sku != ""
    ? local.sku_map[var.spec.sku]
    : "PerGB2018"
  )

  # Identity type wire values. Workspaces accept exactly SystemAssigned or
  # UserAssigned -- the combined model does not exist on this resource.
  identity_type_map = {
    "SYSTEM_ASSIGNED" = "SystemAssigned"
    "USER_ASSIGNED"   = "UserAssigned"
  }
}
