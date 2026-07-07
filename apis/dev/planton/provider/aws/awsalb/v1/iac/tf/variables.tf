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
  description = "AwsAlb specification"
  type = object({
    region = string
    subnets = list(string)
    security_groups = optional(list(string), [])
    internal = optional(bool, false)
    ip_address_type = optional(string, "")
    delete_protection_enabled = optional(bool, false)
    idle_timeout_seconds = optional(number, 0)
    client_keep_alive_seconds = optional(number, 0)
    http2_enabled = optional(bool)
    waf_fail_open_enabled = optional(bool, false)
    web_acl_arn = optional(string, "")
    zonal_shift_enabled = optional(bool, false)
    drop_invalid_header_fields = optional(bool, false)
    preserve_host_header = optional(bool, false)
    xff_client_port_enabled = optional(bool, false)
    xff_header_processing_mode = optional(string, "")
    desync_mitigation_mode = optional(string, "")
    tls_version_and_cipher_suite_headers_enabled = optional(bool, false)
    access_logs = optional(object({
      bucket = string
      prefix = optional(string, "")
    }))
    connection_logs = optional(object({
      bucket = string
      prefix = optional(string, "")
    }))
    health_check_logs = optional(object({
      bucket = string
      prefix = optional(string, "")
    }))
    dns = optional(object({
      enabled = optional(bool, false)
      route53_zone_id = optional(string, "")
      hostnames = optional(list(string), [])
    }))
  })
}
