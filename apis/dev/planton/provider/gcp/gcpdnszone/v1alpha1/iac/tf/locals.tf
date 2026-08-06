locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project (ambient credentials decide).
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.metadata.name
    "planton-ai_kind"     = "gcpdnszone"
  }

  id_label = (
    var.metadata.id != null && var.metadata.id != ""
  ) ? { "planton-ai_id" = var.metadata.id } : {}

  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "planton-ai_organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "planton-ai_environment" = var.metadata.env } : {}

  platform_labels = merge(local.base_labels, local.id_label, local.org_label, local.env_label)
  labels          = merge(local.platform_labels, var.spec.labels)

  managed_zone_name = replace(var.metadata.name, ".", "-")

  # When dns_name is omitted, derive from metadata.name (legacy behavior).
  zone_dns_name = var.spec.dns_name != "" ? var.spec.dns_name : "${var.metadata.name}."

  description = var.spec.description != "" ? var.spec.description : "managed-zone for ${var.metadata.name}"

  visibility = var.spec.visibility != "" ? var.spec.visibility : "public"
}
