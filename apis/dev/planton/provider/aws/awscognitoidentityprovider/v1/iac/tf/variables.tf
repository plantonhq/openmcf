variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsCognitoIdentityProvider specification"
  type = object({
    region = string
    user_pool_id = string
    provider_name = string
    provider_type = string
    google = optional(object({
      client_id = string
      client_secret = string
      authorize_scopes = string
    }))
    facebook = optional(object({
      client_id = string
      client_secret = string
      authorize_scopes = string
      api_version = optional(string, "")
    }))
    login_with_amazon = optional(object({
      client_id = string
      client_secret = string
      authorize_scopes = string
    }))
    sign_in_with_apple = optional(object({
      client_id = string
      team_id = string
      key_id = string
      private_key = string
      authorize_scopes = string
    }))
    oidc = optional(object({
      client_id = string
      oidc_issuer = string
      authorize_scopes = optional(string, "")
      client_secret = optional(string, "")
      attributes_request_method = optional(string, "")
      authorize_url = optional(string, "")
      token_url = optional(string, "")
      attributes_url = optional(string, "")
      jwks_uri = optional(string, "")
      attributes_url_add_attributes = optional(bool, false)
    }))
    saml = optional(object({
      metadata_file = optional(string, "")
      metadata_url = optional(string, "")
      idp_sign_out = optional(bool, false)
      idp_init = optional(bool, false)
      encrypted_responses = optional(bool, false)
      request_signing_algorithm = optional(string, "")
    }))
    attribute_mapping = optional(map(string), {})
    idp_identifiers = optional(list(string), [])
  })
}
