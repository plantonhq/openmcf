locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  location     = var.spec.location
  display_name = var.spec.display_name

  # Vertex AI requires the endpoint's `name` to be numeric (max 10 digits,
  # no leading zero) and the API will not generate one — so when the spec
  # omits endpoint_name, derive a stable ID from the resource's own
  # identity. The derivation (sha256 of "org/env/name", first 48 bits,
  # mapped into [1000000000, 9999999999]) is implemented IDENTICALLY in
  # the Pulumi module: the same manifest yields the same endpoint ID on
  # either engine, and re-applies never regenerate it.
  endpoint_identity = format(
    "%s/%s/%s",
    var.metadata.org != null ? var.metadata.org : "",
    var.metadata.env != null ? var.metadata.env : "",
    var.metadata.name,
  )
  derived_endpoint_name = tostring(
    1000000000 + parseint(substr(sha256(local.endpoint_identity), 0, 12), 16) % 9000000000
  )
  endpoint_name = var.spec.endpoint_name != "" ? var.spec.endpoint_name : local.derived_endpoint_name

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = lower(var.metadata.name)
    "planton-ai_kind"     = "gcpvertexaiendpoint"
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
