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
  description = "AwsSesAccountSettings specification"
  type = object({
    region = string
    suppression = optional(object({
      reasons = optional(list(string), [])
    }))
    vdm = optional(object({
      enabled                   = optional(bool, false)
      engagement_metrics        = optional(bool)
      optimized_shared_delivery = optional(bool)
    }))
  })
}
