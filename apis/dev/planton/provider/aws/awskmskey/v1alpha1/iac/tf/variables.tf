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
  })
}
