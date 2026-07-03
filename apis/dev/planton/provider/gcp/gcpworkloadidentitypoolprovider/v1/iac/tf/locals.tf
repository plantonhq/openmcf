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

  # AWS, SAML, and X.509 issuers have server-side default mappings, so an
  # empty map stays unset (null); OIDC requires an explicit mapping, which the
  # variable validation already guarantees carries google.subject.
  attribute_mapping = length(var.spec.attribute_mapping) > 0 ? var.spec.attribute_mapping : null
}
