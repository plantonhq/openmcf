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
  description = "CloudflareZoneTlsSettings specification"
  type = object({
    zone_id               = string
    universal_ssl_enabled = optional(bool)
    total_tls = optional(object({
      enabled               = optional(bool, false)
      certificate_authority = optional(string)
    }))
    auto_origin_tls_kex         = optional(bool)
    origin_tls_compliance_modes = optional(list(string), [])
    hostname_settings = optional(list(object({
      hostname        = string
      min_tls_version = optional(string)
      http2           = optional(bool)
      ciphers         = optional(list(string), [])
    })), [])
    ca_hostname_associations = optional(list(object({
      hostnames           = list(string)
      mtls_certificate_id = optional(string, "")
    })), [])
  })
}