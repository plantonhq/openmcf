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
  description = "AwsSnsSubscription specification"
  type = object({
    region = string
    topic_arn = string
    protocol = string
    endpoint = string
    filter_policy = optional(any)
    filter_policy_scope = optional(string, "")
    raw_message_delivery = optional(bool, false)
    dead_letter_config = optional(object({
      dead_letter_target_arn = string
    }))
    delivery_policy = optional(string, "")
    replay_policy = optional(any)
    subscription_role_arn = optional(string, "")
    endpoint_auto_confirms = optional(bool, false)
    confirmation_timeout_minutes = optional(number, 0)
  })
}
