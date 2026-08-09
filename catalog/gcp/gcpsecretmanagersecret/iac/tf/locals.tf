locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The cloud-side secret ID defaults to metadata.name when the spec leaves
  # secret_id empty — the same naming basis every kind uses.
  secret_id = (
    var.spec.secret_id != null && var.spec.secret_id != ""
    ? var.spec.secret_id
    : var.metadata.name
  )

  # Empty region means a GLOBAL secret (replication control); a set region
  # selects the regional resources. All three resource pairs in main.tf key
  # off this.
  is_regional = var.spec.region != null && var.spec.region != ""

  # The same planton-ai_* label set the Pulumi module applies, so a resource
  # is attributable to its Planton object regardless of the engine that
  # created it.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.metadata.name
    "planton-ai_kind"     = "gcpsecretmanagersecret"
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

  # User labels first: the platform labels win on key conflicts.
  final_labels = merge(var.spec.labels, local.base_labels, local.org_label, local.env_label, local.id_label)

  # IAM grants keyed by "role|member" — additive iam_member semantics; the
  # composite key lets one member hold several roles (and vice versa)
  # without collisions.
  iam_members = {
    for member in var.spec.iam_members :
    "${member.role}|${member.member}" => member
  }
}
