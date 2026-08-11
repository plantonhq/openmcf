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
  description = "GcpEventarcMessageBus specification"
  type = object({
    project_id     = optional(string, "")
    location       = string
    message_bus_id = optional(string, "")
    display_name   = optional(string, "")
    labels         = optional(map(string), {})
    annotations    = optional(map(string), {})
    crypto_key     = optional(string, "")
    log_severity   = optional(string, "")
    google_api_sources = optional(list(object({
      source_id    = string
      display_name = optional(string, "")
      labels       = optional(map(string), {})
      annotations  = optional(map(string), {})
      crypto_key   = optional(string, "")
      log_severity = optional(string, "")
    })), [])
    pipelines = optional(list(object({
      pipeline_id = string
      destination = object({
        http_endpoint = optional(object({
          uri                      = string
          message_binding_template = optional(string, "")
          network_attachment       = optional(string, "")
        }))
        topic       = optional(string, "")
        workflow    = optional(string, "")
        message_bus = optional(string, "")
      })
      authentication = optional(object({
        google_oidc = optional(object({
          service_account = string
          audience        = optional(string, "")
        }))
        oauth_token = optional(object({
          service_account = string
          scope           = optional(string, "")
        }))
      }))
      input_payload_format = optional(object({
        avro = optional(object({
          schema_definition = optional(string, "")
        }))
        json = optional(bool, false)
        protobuf = optional(object({
          schema_definition = optional(string, "")
        }))
      }))
      output_payload_format = optional(object({
        avro = optional(object({
          schema_definition = optional(string, "")
        }))
        json = optional(bool, false)
        protobuf = optional(object({
          schema_definition = optional(string, "")
        }))
      }))
      mediation_transformation_template = optional(string, "")
      retry_policy = optional(object({
        max_attempts    = optional(number, 0)
        min_retry_delay = optional(string, "")
        max_retry_delay = optional(string, "")
      }))
      display_name = optional(string, "")
      labels       = optional(map(string), {})
      annotations  = optional(map(string), {})
      crypto_key   = optional(string, "")
      log_severity = optional(string, "")
    })), [])
    enrollments = optional(list(object({
      enrollment_id = string
      cel_match     = string
      pipeline      = string
      display_name  = optional(string, "")
      labels        = optional(map(string), {})
      annotations   = optional(map(string), {})
    })), [])
    deletion_policy = optional(string, "")
  })
}