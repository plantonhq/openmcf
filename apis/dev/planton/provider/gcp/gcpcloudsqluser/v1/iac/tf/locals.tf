locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Empty optional strings become null so the provider omits them from the
  # API payload (IAM users have no password; non-MySQL engines have no host).
  password = var.spec.password != "" ? var.spec.password : null
  host     = var.spec.host != "" ? var.spec.host : null

  # BUILT_IN is the API's implicit default; the provider diff-suppresses the
  # explicit form, so passing it verbatim is drift-free.
  user_type = var.spec.type

  password_policy = var.spec.password_policy
}
