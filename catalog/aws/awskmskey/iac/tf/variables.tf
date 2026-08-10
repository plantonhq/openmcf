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
  description = "AwsKmsKey specification"
  type = object({
    region = string
    description = optional(string, "")
    key_spec = optional(string, "")
    key_usage = optional(string, "")
    policy = optional(string, "")
    bypass_policy_lockout_safety_check = optional(bool, false)
    disabled = optional(bool, false)
    enable_key_rotation = optional(bool, false)
    rotation_period_in_days = optional(number, 0)
    multi_region = optional(bool, false)
    deletion_window_days = optional(number, 0)
    aliases = optional(list(string), [])
    custom_key_store_id = optional(string, "")
    xks_key_id = optional(string, "")
    grants = optional(list(object({
      name = optional(string, "")
      grantee_principal = string
      operations = list(string)
      retiring_principal = optional(string, "")
      encryption_context_equals = optional(map(string), {})
      encryption_context_subset = optional(map(string), {})
      retire_on_delete = optional(bool, false)
    })), [])
  })
}