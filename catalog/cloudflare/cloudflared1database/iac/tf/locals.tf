# locals.tf

locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Base labels
  base_labels = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "cloudflare_d1_database"
  }

  # Organization label only if var.metadata.org is non-empty
  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  # Environment label only if var.metadata.env is non-empty
  env_label = (
    var.metadata.env != null &&
    try(var.metadata.env, "") != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Merge base, org, and environment labels
  final_labels = merge(local.base_labels, local.org_label, local.env_label)

  # The tfvars converter OMITS unset enums entirely (EmitUnpopulated: false),
  # so an unset region simply never appears in tfvars; the try() covers the
  # omission and the sentinel-name guard below is defensive only -- the
  # converter never emits zero-value enum names.
  region = try(var.spec.region, "")
  primary_location_hint = (
    local.region != "" && local.region != "cloudflare_d1_region_unspecified"
  ) ? local.region : null

  # Data-residency jurisdiction (eu/fedramp); omitted -> null -> no constraint.
  jurisdiction = try(var.spec.jurisdiction, "") != "" ? var.spec.jurisdiction : null
}

