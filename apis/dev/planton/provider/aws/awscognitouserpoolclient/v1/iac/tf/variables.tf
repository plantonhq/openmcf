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
  description = "AwsCognitoUserPoolClient specification"
  type = object({
    region = string
    user_pool_id = string
    generate_secret = optional(bool, false)
    allowed_oauth_flows_user_pool_client = optional(bool, false)
    allowed_oauth_flows = optional(list(string), [])
    allowed_oauth_scopes = optional(list(string), [])
    callback_urls = optional(list(string), [])
    logout_urls = optional(list(string), [])
    default_redirect_uri = optional(string, "")
    supported_identity_providers = optional(list(string), [])
    explicit_auth_flows = optional(list(string), [])
    auth_session_validity = optional(number)
    access_token_validity = optional(number)
    id_token_validity = optional(number)
    refresh_token_validity = optional(number)
    token_validity_units = optional(object({
      access_token = optional(string, "")
      id_token = optional(string, "")
      refresh_token = optional(string, "")
    }))
    refresh_token_rotation = optional(object({
      feature = string
      retry_grace_period_seconds = optional(number, 0)
    }))
    enable_token_revocation = optional(bool)
    enable_propagate_additional_user_context_data = optional(bool, false)
    prevent_user_existence_errors = optional(string, "")
    read_attributes = optional(list(string), [])
    write_attributes = optional(list(string), [])
    analytics_configuration = optional(object({
      application_arn = optional(string, "")
      application_id = optional(string, "")
      external_id = optional(string, "")
      role_arn = optional(string, "")
      user_data_shared = optional(bool, false)
    }))
  })
}
