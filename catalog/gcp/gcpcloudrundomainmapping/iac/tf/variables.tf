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
  description = "GcpCloudRunDomainMapping specification"
  type = object({
    project_id       = optional(string, "")
    region           = string
    domain           = string
    route            = string
    certificate_mode = optional(string, "")
    force_override   = optional(bool, false)
    namespace        = optional(string, "")
    labels           = optional(map(string), {})
    annotations      = optional(map(string), {})
    deletion_policy  = optional(string, "")
  })
}