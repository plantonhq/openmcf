variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name        = string
    id          = optional(string, "")
    org         = optional(string, "")
    env         = optional(string, "")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    tags        = optional(list(string), [])
  })
}

variable "spec" {
  description = "GcpIdentityPlatformTenant specification"
  type = object({
    project_id               = optional(string, "")
    display_name             = string
    allow_password_signup    = optional(bool, false)
    enable_email_link_signin = optional(bool, false)
    disable_auth             = optional(bool, false)
    client_permissions = optional(object({
      disabled_user_signup   = optional(bool, false)
      disabled_user_deletion = optional(bool, false)
    }))
    default_supported_idps = optional(list(object({
      idp_id        = string
      client_id     = string
      client_secret = string
      enabled       = optional(bool)
    })), [])
    oauth_idp_configs = optional(list(object({
      name          = string
      display_name  = string
      issuer        = string
      client_id     = string
      client_secret = optional(string, "")
      enabled       = optional(bool)
    })), [])
    inbound_saml_configs = optional(list(object({
      name         = string
      display_name = string
      enabled      = optional(bool)
      idp_config = object({
        idp_entity_id = string
        sso_url       = string
        sign_request  = optional(bool, false)
        idp_certificates = optional(list(object({
          x509_certificate = optional(string, "")
        })), [])
      })
      sp_config = object({
        callback_uri = string
        sp_entity_id = string
      })
    })), [])
    deletion_policy = optional(string, "")
  })
}
