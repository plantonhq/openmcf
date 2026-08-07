locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Instance name defaults to metadata.name when instance_name is omitted.
  instance_name = var.spec.instance_name != "" ? var.spec.instance_name : var.metadata.name

  # 0-means-unset scalars translate to null so the provider applies its own
  # defaults (an unset capacity on a PROVISIONED instance defaults to 1 node
  # server-side).
  num_nodes        = var.spec.num_nodes > 0 ? var.spec.num_nodes : null
  processing_units = var.spec.processing_units > 0 ? var.spec.processing_units : null

  instance_type                = var.spec.instance_type != "" ? var.spec.instance_type : null
  edition                      = var.spec.edition != "" ? var.spec.edition : null
  default_backup_schedule_type = var.spec.default_backup_schedule_type != "" ? var.spec.default_backup_schedule_type : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.instance_name
    "planton-ai_kind"     = "gcpspannerinstance"
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
