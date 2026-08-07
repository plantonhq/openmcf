locals {
  # The graph identity; the provider-visible name is spec.provider_name (it is
  # the IdP's AWS identity within its pool and is ForceNew).
  resource_name = var.metadata.name

  # AWS takes the provider configuration as one flat string map whose keys
  # depend on the provider type. Exactly one typed config block is present
  # (spec CEL); each is translated to the key convention its provider type
  # expects: snake_case for the OAuth family, PascalCase for SAML.
  google_details = var.spec.google != null ? {
    client_id        = var.spec.google.client_id
    client_secret    = var.spec.google.client_secret
    authorize_scopes = var.spec.google.authorize_scopes
  } : null

  facebook_details = var.spec.facebook != null ? merge({
    client_id        = var.spec.facebook.client_id
    client_secret    = var.spec.facebook.client_secret
    authorize_scopes = var.spec.facebook.authorize_scopes
    }, var.spec.facebook.api_version != "" ? {
    api_version = var.spec.facebook.api_version
  } : {}) : null

  login_with_amazon_details = var.spec.login_with_amazon != null ? {
    client_id        = var.spec.login_with_amazon.client_id
    client_secret    = var.spec.login_with_amazon.client_secret
    authorize_scopes = var.spec.login_with_amazon.authorize_scopes
  } : null

  sign_in_with_apple_details = var.spec.sign_in_with_apple != null ? {
    client_id        = var.spec.sign_in_with_apple.client_id
    team_id          = var.spec.sign_in_with_apple.team_id
    key_id           = var.spec.sign_in_with_apple.key_id
    private_key      = var.spec.sign_in_with_apple.private_key
    authorize_scopes = var.spec.sign_in_with_apple.authorize_scopes
  } : null

  # OIDC optionals are only sent when set -- Cognito auto-discovers endpoint
  # URLs from the issuer's .well-known document when overrides are absent.
  oidc_details = var.spec.oidc != null ? merge(
    {
      client_id   = var.spec.oidc.client_id
      oidc_issuer = var.spec.oidc.oidc_issuer
    },
    var.spec.oidc.authorize_scopes != "" ? { authorize_scopes = var.spec.oidc.authorize_scopes } : {},
    var.spec.oidc.client_secret != "" ? { client_secret = var.spec.oidc.client_secret } : {},
    var.spec.oidc.attributes_request_method != "" ? { attributes_request_method = var.spec.oidc.attributes_request_method } : {},
    var.spec.oidc.authorize_url != "" ? { authorize_url = var.spec.oidc.authorize_url } : {},
    var.spec.oidc.token_url != "" ? { token_url = var.spec.oidc.token_url } : {},
    var.spec.oidc.attributes_url != "" ? { attributes_url = var.spec.oidc.attributes_url } : {},
    var.spec.oidc.jwks_uri != "" ? { jwks_uri = var.spec.oidc.jwks_uri } : {},
    var.spec.oidc.attributes_url_add_attributes ? { attributes_url_add_attributes = "true" } : {},
  ) : null

  saml_details = var.spec.saml != null ? merge(
    var.spec.saml.metadata_file != "" ? { MetadataFile = var.spec.saml.metadata_file } : {},
    var.spec.saml.metadata_url != "" ? { MetadataURL = var.spec.saml.metadata_url } : {},
    var.spec.saml.idp_sign_out ? { IDPSignout = "true" } : {},
    var.spec.saml.idp_init ? { IDPInit = "true" } : {},
    var.spec.saml.encrypted_responses ? { EncryptedResponses = "true" } : {},
    var.spec.saml.request_signing_algorithm != "" ? { RequestSigningAlgorithm = var.spec.saml.request_signing_algorithm } : {},
  ) : null

  # Select whichever typed config is present (exactly one, per the spec CEL).
  provider_details = coalesce(
    local.google_details,
    local.facebook_details,
    local.login_with_amazon_details,
    local.sign_in_with_apple_details,
    local.oidc_details,
    local.saml_details,
    {}
  )
}
