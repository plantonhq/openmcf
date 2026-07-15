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
  description = "AwsCloudFront specification"
  type = object({
    region = string
    enabled = optional(bool)
    aliases = optional(list(string), [])
    comment = optional(string, "")
    default_root_object = optional(string, "")
    http_version = optional(string, "")
    is_ipv6_enabled = optional(bool, false)
    price_class = optional(string, "")
    web_acl_arn = optional(string, "")
    origins = list(object({
      origin_id = string
      domain_name = string
      origin_path = optional(string, "")
      connection_attempts = optional(number, 0)
      connection_timeout_seconds = optional(number, 0)
      custom_headers = optional(list(object({
        name = string
        value = string
      })), [])
      origin_shield = optional(object({
        origin_shield_region = string
      }))
      s3_origin = optional(object({
        create_origin_access_control = optional(bool, false)
        origin_access_control_id = optional(string, "")
        origin_access_identity = optional(string, "")
      }))
      custom_origin = optional(object({
        protocol_policy = string
        http_port = optional(number, 0)
        https_port = optional(number, 0)
        ssl_protocols = optional(list(string), [])
        keepalive_timeout_seconds = optional(number, 0)
        read_timeout_seconds = optional(number, 0)
      }))
      vpc_origin = optional(object({
        vpc_origin_id = string
        keepalive_timeout_seconds = optional(number, 0)
        read_timeout_seconds = optional(number, 0)
      }))
    }))
    origin_groups = optional(list(object({
      origin_group_id = string
      member_origin_ids = list(string)
      failover_status_codes = list(number)
    })), [])
    default_cache_behavior = object({
      target_origin_id = string
      viewer_protocol_policy = string
      allowed_methods = optional(list(string), [])
      cached_methods = optional(list(string), [])
      compress = optional(bool, false)
      cache_policy_id = optional(string, "")
      origin_request_policy_id = optional(string, "")
      response_headers_policy_id = optional(string, "")
      forwarded_values = optional(object({
        query_string = optional(bool, false)
        query_string_cache_keys = optional(list(string), [])
        headers = optional(list(string), [])
        cookies_forward = string
        whitelisted_cookie_names = optional(list(string), [])
      }))
      min_ttl_seconds = optional(number, 0)
      default_ttl_seconds = optional(number, 0)
      max_ttl_seconds = optional(number, 0)
      function_associations = optional(list(object({
        event_type = string
        function_arn = string
      })), [])
      lambda_function_associations = optional(list(object({
        event_type = string
        lambda_arn = string
        include_body = optional(bool, false)
      })), [])
      trusted_key_group_ids = optional(list(string), [])
      trusted_signers = optional(list(string), [])
      field_level_encryption_id = optional(string, "")
      realtime_log_config_arn = optional(string, "")
      smooth_streaming = optional(bool, false)
      grpc_enabled = optional(bool, false)
    })
    ordered_cache_behaviors = optional(list(object({
      path_pattern = string
      behavior = object({
        target_origin_id = string
        viewer_protocol_policy = string
        allowed_methods = optional(list(string), [])
        cached_methods = optional(list(string), [])
        compress = optional(bool, false)
        cache_policy_id = optional(string, "")
        origin_request_policy_id = optional(string, "")
        response_headers_policy_id = optional(string, "")
        forwarded_values = optional(object({
          query_string = optional(bool, false)
          query_string_cache_keys = optional(list(string), [])
          headers = optional(list(string), [])
          cookies_forward = string
          whitelisted_cookie_names = optional(list(string), [])
        }))
        min_ttl_seconds = optional(number, 0)
        default_ttl_seconds = optional(number, 0)
        max_ttl_seconds = optional(number, 0)
        function_associations = optional(list(object({
          event_type = string
          function_arn = string
        })), [])
        lambda_function_associations = optional(list(object({
          event_type = string
          lambda_arn = string
          include_body = optional(bool, false)
        })), [])
        trusted_key_group_ids = optional(list(string), [])
        trusted_signers = optional(list(string), [])
        field_level_encryption_id = optional(string, "")
        realtime_log_config_arn = optional(string, "")
        smooth_streaming = optional(bool, false)
        grpc_enabled = optional(bool, false)
      })
    })), [])
    viewer_certificate = optional(object({
      acm_certificate_arn = optional(string, "")
      iam_certificate_id = optional(string, "")
      ssl_support_method = optional(string, "")
      minimum_protocol_version = optional(string, "")
    }))
    custom_error_responses = optional(list(object({
      error_code = number
      response_code = optional(number, 0)
      response_page_path = optional(string, "")
      error_caching_min_ttl_seconds = optional(number, 0)
    })), [])
    geo_restriction = optional(object({
      restriction_type = string
      locations = list(string)
    }))
    logging = optional(object({
      bucket = string
      prefix = optional(string, "")
      include_cookies = optional(bool, false)
    }))
    wait_for_deployment = optional(bool)
    retain_on_delete = optional(bool, false)
    enable_additional_metrics = optional(bool, false)
  })
}
