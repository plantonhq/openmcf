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
  description = "AwsHttpApiGateway specification"
  type = object({
    region = string
    description = optional(string, "")
    api_version = optional(string, "")
    cors_configuration = optional(object({
      allow_origins = optional(list(string), [])
      allow_methods = optional(list(string), [])
      allow_headers = optional(list(string), [])
      expose_headers = optional(list(string), [])
      max_age_seconds = optional(number, 0)
      allow_credentials = optional(bool, false)
    }))
    disable_execute_api_endpoint = optional(bool, false)
    ip_address_type = optional(string, "")
    stage = optional(object({
      name = optional(string, "")
      auto_deploy = optional(bool)
      description = optional(string, "")
      access_log = optional(object({
        destination_arn = string
        format = string
      }))
      default_throttle = optional(object({
        burst_limit = optional(number, 0)
        rate_limit = optional(number, 0)
      }))
      detailed_metrics_enabled = optional(bool, false)
      route_settings = optional(list(object({
        route_key = string
        throttling_burst_limit = optional(number, 0)
        throttling_rate_limit = optional(number, 0)
        detailed_metrics_enabled = optional(bool, false)
      })), [])
      stage_variables = optional(map(string), {})
    }))
    routes = list(object({
      route_key = string
      integration = object({
        integration_type = string
        integration_uri = optional(string, "")
        integration_subtype = optional(string, "")
        payload_format_version = optional(string, "")
        integration_method = optional(string, "")
        timeout_milliseconds = optional(number, 0)
        connection_type = optional(string, "")
        connection_id = optional(string, "")
        credentials_arn = optional(string, "")
        request_parameters = optional(map(string), {})
        response_parameters = optional(list(object({
          status_code = optional(string, "")
          mappings = optional(map(string), {})
        })), [])
        tls_server_name_to_verify = optional(string, "")
        description = optional(string, "")
      })
      authorization_type = optional(string, "")
      authorizer_name = optional(string, "")
      authorization_scopes = optional(list(string), [])
      operation_name = optional(string, "")
    }))
    authorizers = optional(list(object({
      name = string
      authorizer_type = string
      jwt_configuration = optional(object({
        issuer = string
        audiences = optional(list(string), [])
      }))
      authorizer_uri = optional(string, "")
      authorizer_credentials_arn = optional(string, "")
      identity_sources = optional(list(string), [])
      result_ttl_seconds = optional(number, 0)
      enable_simple_responses = optional(bool, false)
      authorizer_payload_format_version = optional(string, "")
    })), [])
  })
}