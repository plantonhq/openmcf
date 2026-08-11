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
  description = "GcpMonitoringUptimeCheck specification"
  type = object({
    project_id         = optional(string, "")
    display_name       = optional(string, "")
    timeout            = string
    period             = optional(string, "")
    checker_type       = optional(string, "")
    selected_regions   = optional(list(string), [])
    log_check_failures = optional(bool, false)
    monitored_resource = optional(object({
      type   = string
      labels = optional(map(string), {})
    }))
    resource_group = optional(object({
      group_id      = optional(string, "")
      resource_type = optional(string, "")
    }))
    synthetic_monitor = optional(object({
      cloud_function = string
    }))
    http_check = optional(object({
      path                = optional(string, "")
      port                = optional(number, 0)
      request_method      = optional(string, "")
      use_ssl             = optional(bool, false)
      validate_ssl        = optional(bool, false)
      body                = optional(string, "")
      content_type        = optional(string, "")
      custom_content_type = optional(string, "")
      headers             = optional(map(string), {})
      mask_headers        = optional(bool, false)
      auth_info = optional(object({
        username = string
        password = optional(string, "")
      }))
      service_agent_authentication = optional(object({
        type = optional(string, "")
      }))
      accepted_response_status_codes = optional(list(object({
        status_class = optional(string, "")
        status_value = optional(number, 0)
      })), [])
      ping_config = optional(object({
        pings_count = optional(number, 0)
      }))
    }))
    tcp_check = optional(object({
      port = optional(number, 0)
      ping_config = optional(object({
        pings_count = optional(number, 0)
      }))
    }))
    content_matchers = optional(list(object({
      content = string
      matcher = optional(string, "")
      json_path_matcher = optional(object({
        json_path    = string
        json_matcher = optional(string, "")
      }))
    })), [])
    labels          = optional(map(string), {})
    deletion_policy = optional(string, "")
  })
}