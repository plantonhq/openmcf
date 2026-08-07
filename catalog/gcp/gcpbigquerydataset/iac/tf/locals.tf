locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # 0-/empty-means-unset scalars translate to null so the provider applies
  # its own server-side defaults (e.g. max_time_travel_hours defaults to 168
  # and storage_billing_model to LOGICAL when omitted).
  friendly_name                   = var.spec.friendly_name != "" ? var.spec.friendly_name : null
  description                     = var.spec.description != "" ? var.spec.description : null
  default_table_expiration_ms    = var.spec.default_table_expiration_ms > 0 ? var.spec.default_table_expiration_ms : null
  default_partition_expiration_ms = var.spec.default_partition_expiration_ms > 0 ? var.spec.default_partition_expiration_ms : null
  # The provider types max_time_travel_hours as a string.
  max_time_travel_hours = var.spec.max_time_travel_hours > 0 ? tostring(var.spec.max_time_travel_hours) : null
  default_collation     = var.spec.default_collation != "" ? var.spec.default_collation : null
  storage_billing_model = var.spec.storage_billing_model != "" ? var.spec.storage_billing_model : null
  kms_key_name          = var.spec.kms_key_name != "" ? var.spec.kms_key_name : null

  resource_tags = length(var.spec.resource_tags) > 0 ? var.spec.resource_tags : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.spec.dataset_id
    "planton-ai_kind"     = "gcpbigquerydataset"
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
