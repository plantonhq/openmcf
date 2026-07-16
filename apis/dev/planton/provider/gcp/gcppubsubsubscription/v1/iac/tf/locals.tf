locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  subscription_name = var.spec.subscription_name
  topic             = var.spec.topic

  ack_deadline_seconds       = var.spec.ack_deadline_seconds > 0 ? var.spec.ack_deadline_seconds : null
  message_retention_duration = var.spec.message_retention_duration != "" ? var.spec.message_retention_duration : null
  filter                     = var.spec.filter != "" ? var.spec.filter : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.spec.subscription_name
    "planton-ai_kind"     = "gcppubsubsubscription"
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
