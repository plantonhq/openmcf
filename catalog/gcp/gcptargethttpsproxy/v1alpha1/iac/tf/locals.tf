locals {
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
  # proxy_name empty — the same naming basis every kind uses.
  proxy_name = (
    var.spec.proxy_name != null && var.spec.proxy_name != ""
    ? var.spec.proxy_name
    : var.metadata.name
  )

  description = var.spec.description != "" ? var.spec.description : null

  # Empty lists become null so the API never sees an explicit empty list —
  # exactly one certificate mechanism may be present (spec CEL enforces it).
  ssl_certificates                 = length(var.spec.ssl_certificates) > 0 ? var.spec.ssl_certificates : null
  certificate_manager_certificates = length(var.spec.certificate_manager_certificates) > 0 ? var.spec.certificate_manager_certificates : null
  certificate_map                  = var.spec.certificate_map != "" ? var.spec.certificate_map : null

  ssl_policy        = var.spec.ssl_policy != "" ? var.spec.ssl_policy : null
  server_tls_policy = var.spec.server_tls_policy != "" ? var.spec.server_tls_policy : null

  # The middleware default (NONE) matches GCP's own default; null lets the
  # API compute NONE when unset.
  quic_override = var.spec.quic_override != "" ? var.spec.quic_override : null

  # Empty keeps GCP's default (DISABLED); the field is immutable, so a value
  # is only sent when the user chose a mode.
  tls_early_data = var.spec.tls_early_data != "" ? var.spec.tls_early_data : null

  # 0 means "let GCP apply its default" (610s on EXTERNAL_MANAGED); sending
  # null keeps the API in charge of the default.
  http_keep_alive_timeout_sec = (
    var.spec.http_keep_alive_timeout_sec != 0
    ? var.spec.http_keep_alive_timeout_sec
    : null
  )

  # proxy_bind is a Traffic Director lever; the API default is false, so only
  # an explicit true is worth sending (null lets the API compute it).
  proxy_bind = var.spec.proxy_bind ? true : null
}
