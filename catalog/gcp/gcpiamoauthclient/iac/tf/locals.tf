locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The cloud-side client ID defaults to metadata.name when the spec leaves
  # oauth_client_id empty — the same naming basis every kind uses.
  oauth_client_id = (
    var.spec.oauth_client_id != null && var.spec.oauth_client_id != ""
    ? var.spec.oauth_client_id
    : var.metadata.name
  )

  # The client's location defaults to "global" — the documented home for
  # workforce OAuth clients (the spec comment records the contract).
  location = var.spec.location != "" ? var.spec.location : "global"

  # Credentials keyed by their immutable resource IDs so plans stay stable
  # as list order changes.
  credentials = {
    for credential in var.spec.credentials : credential.credential_id => credential
  }
}
