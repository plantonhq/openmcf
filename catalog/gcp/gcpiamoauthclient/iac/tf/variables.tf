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
  description = "GcpIamOauthClient specification"
  type = object({
    project_id            = optional(string, "")
    location              = optional(string, "")
    oauth_client_id       = optional(string, "")
    display_name          = optional(string, "")
    description           = optional(string, "")
    disabled              = optional(bool, false)
    client_type           = optional(string, "")
    allowed_grant_types   = list(string)
    allowed_scopes        = list(string)
    allowed_redirect_uris = list(string)
    credentials = optional(list(object({
      credential_id = string
      display_name  = optional(string, "")
      disabled      = optional(bool, false)
    })), [])
    deletion_policy = optional(string, "")
  })
}
