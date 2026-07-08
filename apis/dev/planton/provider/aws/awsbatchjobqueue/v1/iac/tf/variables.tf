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
  description = "AwsBatchJobQueue specification"
  type = object({
    region = string
    state = optional(string)
    priority = optional(number, 0)
    compute_environment_order = list(object({
      order = optional(number, 0)
      compute_environment = string
    }))
    scheduling_policy = optional(string, "")
    job_state_time_limit_actions = optional(list(object({
      action = string
      max_time_seconds = number
      reason = string
      state = string
    })), [])
  })
}
