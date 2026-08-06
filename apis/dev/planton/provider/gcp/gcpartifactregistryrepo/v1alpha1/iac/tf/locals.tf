locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The repository ID falls back to metadata.name — one honest identity,
  # never a format-derived suffix.
  repository_id = var.spec.repository_id != "" ? var.spec.repository_id : var.metadata.name

  location = var.spec.location

  # Empty optional strings become null so the provider applies its own
  # defaults (STANDARD_REPOSITORY mode, project-inherited scanning) instead
  # of receiving an empty string it would reject.
  mode                              = var.spec.mode != "" ? var.spec.mode : null
  description                       = var.spec.description != "" ? var.spec.description : null
  kms_key_name                      = var.spec.kms_key_name != "" ? var.spec.kms_key_name : null
  vulnerability_scanning_enablement = var.spec.vulnerability_scanning_enablement != "" ? var.spec.vulnerability_scanning_enablement : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.repository_id
    "planton-ai_kind"     = "gcpartifactregistryrepo"
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

  # Additive IAM grants keyed by (role, member) — the grant's identity.
  # A condition never changes the key: the same (role, member) pair with a
  # different condition is a replacement of that one grant.
  iam_members = {
    for m in var.spec.iam_members :
    "${m.role}/${m.member}" => m
  }
}
