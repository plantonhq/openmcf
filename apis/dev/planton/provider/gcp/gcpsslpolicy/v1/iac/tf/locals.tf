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

  # The cloud-side name defaults to metadata.name when the spec leaves
  # ssl_policy_name empty — the same naming basis every kind uses.
  ssl_policy_name = (
    var.spec.ssl_policy_name != null && var.spec.ssl_policy_name != ""
    ? var.spec.ssl_policy_name
    : var.metadata.name
  )

  # Empty region means a GLOBAL SSL policy; a set region selects the regional
  # resource. Both branches in main.tf key off this.
  is_regional = var.spec.region != null && var.spec.region != ""

  # Unset enum-like strings fall through to the API defaults (COMPATIBLE /
  # TLS_1_0) as null rather than being sent as empty strings, which the API
  # would reject. Never coalesce() here — HCL's coalesce skips empty strings
  # too and errors when every argument is empty.
  profile         = var.spec.profile != "" ? var.spec.profile : null
  min_tls_version = var.spec.min_tls_version != "" ? var.spec.min_tls_version : null
}
