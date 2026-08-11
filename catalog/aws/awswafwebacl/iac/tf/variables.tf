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
  description = "AwsWafWebAcl specification"
  type = object({
    region = string
    scope = string
    default_action = object({
      type = string
      custom_response = optional(object({
        response_code = number
        custom_response_body_key = optional(string, "")
        response_headers = optional(list(object({
          name = string
          value = string
        })), [])
      }))
      custom_request_headers = optional(list(object({
        name = string
        value = string
      })), [])
    })
    description = optional(string, "")
    rules = optional(any, [])
    visibility_config = optional(object({
      cloudwatch_metrics_enabled = optional(bool, false)
      sampled_requests_enabled = optional(bool, false)
      metric_name = optional(string, "")
    }))
    custom_response_bodies = optional(list(object({
      key = string
      content = string
      content_type = string
    })), [])
    token_domains = optional(list(string), [])
    captcha_config = optional(object({
      immunity_time_sec = number
    }))
    challenge_config = optional(object({
      immunity_time_sec = number
    }))
    association_config = optional(object({
      cloudfront_request_body_limit = optional(string, "")
      api_gateway_request_body_limit = optional(string, "")
      cognito_user_pool_request_body_limit = optional(string, "")
      app_runner_service_request_body_limit = optional(string, "")
      verified_access_instance_request_body_limit = optional(string, "")
    }))
    data_protection_config = optional(object({
      data_protections = list(object({
        field_type = string
        field_keys = optional(list(string), [])
        action = string
        exclude_rule_match_details = optional(bool, false)
        exclude_rate_based_details = optional(bool, false)
      }))
    }))
    logging = optional(object({
      destination_arn = string
      redacted_header_names = optional(list(string), [])
      redact_uri_path = optional(bool, false)
      redact_query_string = optional(bool, false)
      redact_method = optional(bool, false)
      filter = optional(object({
        default_behavior = string
        filters = list(object({
          behavior = string
          requirement = string
          conditions = list(object({
            action = optional(string, "")
            label_name = optional(string, "")
          }))
        }))
      }))
    }))
  })
}
