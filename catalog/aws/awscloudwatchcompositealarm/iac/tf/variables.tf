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
  description = "AwsCloudwatchCompositeAlarm specification"
  type = object({
    region = string
    alarm_rule = string
    alarm_description = optional(string, "")
    actions_enabled = optional(bool)
    alarm_actions = optional(list(string), [])
    ok_actions = optional(list(string), [])
    insufficient_data_actions = optional(list(string), [])
    actions_suppressor = optional(object({
      alarm = string
      wait_period = optional(number, 0)
      extension_period = optional(number, 0)
    }))
  })
}
