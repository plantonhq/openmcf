locals {
  key_ring_id = var.spec.key_ring_id
  key_name    = var.spec.key_name

  # Empty optional strings become null so the provider applies its own
  # defaults (ENCRYPT_DECRYPT purpose, 30-day destroy window) instead of
  # receiving an empty string it would reject.
  purpose                    = var.spec.purpose != "" ? var.spec.purpose : null
  rotation_period            = var.spec.rotation_period != "" ? var.spec.rotation_period : null
  destroy_scheduled_duration = var.spec.destroy_scheduled_duration != "" ? var.spec.destroy_scheduled_duration : null
  crypto_key_backend         = var.spec.crypto_key_backend != "" ? var.spec.crypto_key_backend : null

  version_template = var.spec.version_template

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.spec.key_name
    "planton-ai_kind"     = "gcpkmskey"
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
