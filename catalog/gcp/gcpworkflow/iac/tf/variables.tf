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
  description = "GcpWorkflow specification"
  type = object({
    project_id              = optional(string, "")
    region                  = optional(string, "")
    workflow_name           = optional(string, "")
    description             = optional(string, "")
    labels                  = optional(map(string), {})
    source_contents         = string
    service_account         = optional(string, "")
    crypto_key              = optional(string, "")
    call_log_level          = optional(string, "")
    execution_history_level = optional(string, "")
    user_env_vars           = optional(map(string), {})
    resource_manager_tags   = optional(map(string), {})
    deletion_protection     = optional(bool, true)
    deletion_policy         = optional(string, "")
  })
}