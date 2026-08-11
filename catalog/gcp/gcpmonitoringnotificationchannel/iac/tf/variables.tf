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
  description = "GcpMonitoringNotificationChannel specification"
  type = object({
    project_id     = optional(string, "")
    type           = string
    display_name   = optional(string, "")
    description    = optional(string, "")
    channel_labels = optional(map(string), {})
    sensitive_labels = optional(object({
      auth_token  = optional(string, "")
      password    = optional(string, "")
      service_key = optional(string, "")
    }))
    enabled         = optional(bool)
    force_delete    = optional(bool, false)
    labels          = optional(map(string), {})
    deletion_policy = optional(string, "")
  })
}