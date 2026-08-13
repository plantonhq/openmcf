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
  description = "GcpCertificateMap specification"
  type = object({
    project_id  = optional(string, "")
    map_name    = optional(string, "")
    description = optional(string, "")
    labels      = optional(map(string), {})
    entries = optional(list(object({
      entry_name   = string
      hostname     = optional(string, "")
      matcher      = optional(string, "")
      certificates = list(string)
      description  = optional(string, "")
      labels       = optional(map(string), {})
    })), [])
    deletion_policy = optional(string, "")
  })
}