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
    "resource_kind" = "azure_monitor_metric_alert"
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

  # Enum wire maps. The tfvars wire format carries the FULL proto enum
  # value names -- these maps must match them verbatim.
  aggregation_map = {
    "AVERAGE" = "Average"
    "COUNT"   = "Count"
    "MINIMUM" = "Minimum"
    "MAXIMUM" = "Maximum"
    "TOTAL"   = "Total"
  }

  # One shared operator vocabulary speaks both criteria families' wire
  # dialects (the spec CELs keep each family to its legal subset).
  operator_map = {
    "EQUALS"                = "Equals"
    "GREATER_THAN"          = "GreaterThan"
    "GREATER_THAN_OR_EQUAL" = "GreaterThanOrEqual"
    "LESS_THAN"             = "LessThan"
    "LESS_THAN_OR_EQUAL"    = "LessThanOrEqual"
    "GREATER_OR_LESS_THAN"  = "GreaterOrLessThan"
  }

  dimension_operator_map = {
    "INCLUDE"     = "Include"
    "EXCLUDE"     = "Exclude"
    "STARTS_WITH" = "StartsWith"
  }

  sensitivity_map = {
    "LOW"    = "Low"
    "MEDIUM" = "Medium"
    "HIGH"   = "High"
  }
}
