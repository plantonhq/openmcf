locals {
  # Honor the spec contract: an empty project_id falls back to the
  # provider's default project.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # disk_name falls back to metadata.name — explicit conditional, so both
  # engines derive the identical cloud-side name.
  disk_name = var.spec.disk_name != "" ? var.spec.disk_name : var.metadata.name

  # Empty optional strings become null so the provider applies its own
  # defaults (pd-standard family default type, Google-managed encryption,
  # READ_WRITE_SINGLE access) instead of receiving an empty string it
  # would reject.
  type         = var.spec.type != "" ? var.spec.type : null
  image        = var.spec.image != "" ? var.spec.image : null
  snapshot     = var.spec.source_snapshot != "" ? var.spec.source_snapshot : null
  source_disk  = var.spec.source_disk != "" ? var.spec.source_disk : null
  access_mode  = var.spec.access_mode != "" ? var.spec.access_mode : null
  architecture = var.spec.architecture != "" ? var.spec.architecture : null
  storage_pool = var.spec.storage_pool != "" ? var.spec.storage_pool : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.disk_name
    "planton-ai_kind"     = "gcpcomputedisk"
  }

  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "planton-ai_organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "planton-ai_environment" = var.metadata.env } : {}

  id_label = (
    var.metadata.id != null && var.metadata.id != ""
  ) ? { "planton-ai_id" = var.metadata.id } : {}

  # User labels first so platform attribution labels win on key conflicts —
  # identical merge order to the Pulumi module.
  final_labels = merge(
    var.spec.labels,
    local.base_labels,
    local.org_label,
    local.env_label,
    local.id_label,
  )
}
