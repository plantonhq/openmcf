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
  description = "AwsBatchSchedulingPolicy specification"
  type = object({
    region = string
    compute_reservation = optional(number)
    share_decay_seconds = optional(number)
    share_distributions = optional(list(object({
      share_identifier = string
      weight_factor = optional(number, 0)
    })), [])
  })
}
