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
  description = "GcpIamDenyPolicy specification"
  type = object({
    parent = optional(object({
      project_id      = optional(string, "")
      folder_id       = optional(string, "")
      organization_id = optional(string, "")
    }))
    policy_name  = optional(string, "")
    display_name = optional(string, "")
    rules = list(object({
      description = optional(string, "")
      deny_rule = object({
        denied_principals     = optional(list(string), [])
        exception_principals  = optional(list(string), [])
        denied_permissions    = optional(list(string), [])
        exception_permissions = optional(list(string), [])
        denial_condition = optional(object({
          expression  = string
          title       = optional(string, "")
          description = optional(string, "")
          location    = optional(string, "")
        }))
      })
    }))
    deletion_policy = optional(string, "")
  })
}
