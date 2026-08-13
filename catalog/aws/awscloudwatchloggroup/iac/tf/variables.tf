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
  description = "AwsCloudwatchLogGroup specification"
  type = object({
    region = string
    retention_in_days = optional(number, 0)
    kms_key_id = optional(string, "")
    log_group_class = optional(string, "")
    deletion_protection_enabled = optional(bool)
    metric_filters = optional(list(object({
      name = string
      pattern = optional(string, "")
      apply_on_transformed_logs = optional(bool)
      transformation = object({
        metric_name = string
        metric_namespace = string
        metric_value = string
        default_value = optional(number)
        dimensions = optional(map(string), {})
        unit = optional(string, "")
      })
    })), [])
    subscription_filters = optional(list(object({
      name = string
      destination_arn = string
      filter_pattern = optional(string, "")
      role_arn = optional(string, "")
      distribution = optional(string, "")
      emit_system_fields = optional(list(string), [])
      apply_on_transformed_logs = optional(bool)
    })), [])
    data_protection_policy = optional(any)
    field_index_policy = optional(any)
    log_streams = optional(list(string), [])
    transformer = optional(object({
      processors = list(object({
        add_keys = optional(object({
          entries = list(object({
            key = string
            value = string
            overwrite_if_exists = optional(bool, false)
          }))
        }))
        copy_value = optional(object({
          entries = list(object({
            source = string
            target = string
            overwrite_if_exists = optional(bool, false)
          }))
        }))
        csv = optional(object({
          columns = optional(list(string), [])
          delimiter = optional(string, "")
          quote_character = optional(string, "")
          source = optional(string, "")
        }))
        date_time_converter = optional(object({
          source = string
          target = string
          match_patterns = list(string)
          locale = optional(string, "")
          source_timezone = optional(string, "")
          target_format = optional(string, "")
          target_timezone = optional(string, "")
        }))
        delete_keys = optional(object({
          with_keys = list(string)
        }))
        grok = optional(object({
          match = string
          source = optional(string, "")
        }))
        list_to_map = optional(object({
          source = string
          key = string
          value_key = optional(string, "")
          target = optional(string, "")
          flatten = optional(bool, false)
          flattened_element = optional(string, "")
        }))
        lower_case_string = optional(object({
          with_keys = list(string)
        }))
        move_keys = optional(object({
          entries = list(object({
            source = string
            target = string
            overwrite_if_exists = optional(bool, false)
          }))
        }))
        parse_cloudfront = optional(object({
          source = optional(string, "")
        }))
        parse_json = optional(object({
          source = optional(string, "")
          destination = optional(string, "")
        }))
        parse_key_value = optional(object({
          source = optional(string, "")
          destination = optional(string, "")
          field_delimiter = optional(string, "")
          key_value_delimiter = optional(string, "")
          key_prefix = optional(string, "")
          non_match_value = optional(string, "")
          overwrite_if_exists = optional(bool, false)
        }))
        parse_postgres = optional(object({
          source = optional(string, "")
        }))
        parse_route53 = optional(object({
          source = optional(string, "")
        }))
        parse_to_ocsf = optional(object({
          event_source = string
          ocsf_version = string
          source = optional(string, "")
        }))
        parse_vpc = optional(object({
          source = optional(string, "")
        }))
        parse_waf = optional(object({
          source = optional(string, "")
        }))
        rename_keys = optional(object({
          entries = list(object({
            key = string
            rename_to = string
            overwrite_if_exists = optional(bool, false)
          }))
        }))
        split_string = optional(object({
          entries = list(object({
            source = string
            delimiter = string
          }))
        }))
        substitute_string = optional(object({
          entries = list(object({
            source = string
            from = string
            to = string
          }))
        }))
        trim_string = optional(object({
          with_keys = list(string)
        }))
        type_converter = optional(object({
          entries = list(object({
            key = string
            type = string
          }))
        }))
        upper_case_string = optional(object({
          with_keys = list(string)
        }))
      }))
    }))
  })
}
