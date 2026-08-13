locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The composed tenant IdP configs keyed for for_each: default-supported by
  # the provider's canonical IdP ID, OIDC and SAML by their user-chosen
  # resource names — all immutable identity keys, so plans stay stable as
  # list order changes.
  default_supported_idps = {
    for idp in var.spec.default_supported_idps : idp.idp_id => idp
  }
  oauth_idp_configs = {
    for oidc in var.spec.oauth_idp_configs : oidc.name => oidc
  }
  inbound_saml_configs = {
    for saml in var.spec.inbound_saml_configs : saml.name => saml
  }
}
