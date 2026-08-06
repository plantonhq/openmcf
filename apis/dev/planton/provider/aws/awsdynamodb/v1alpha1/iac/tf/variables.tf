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
  description = "AwsDynamodb specification"
  type = object({
    region = string
    billing_mode = optional(string, "")
    attribute_definitions = optional(list(object({
      name = string
      type = string
    })), [])
    key_schema = optional(list(object({
      attribute_name = string
      key_type = string
    })), [])
    provisioned_throughput = optional(object({
      read_capacity_units = optional(number, 0)
      write_capacity_units = optional(number, 0)
    }))
    on_demand_throughput = optional(object({
      max_read_request_units = optional(number, 0)
      max_write_request_units = optional(number, 0)
    }))
    warm_throughput = optional(object({
      read_units_per_second = optional(number, 0)
      write_units_per_second = optional(number, 0)
    }))
    global_secondary_indexes = optional(list(object({
      name = string
      key_schema = list(object({
        attribute_name = string
        key_type = string
      }))
      projection = object({
        type = string
        non_key_attributes = optional(list(string), [])
      })
      provisioned_throughput = optional(object({
        read_capacity_units = optional(number, 0)
        write_capacity_units = optional(number, 0)
      }))
      on_demand_throughput = optional(object({
        max_read_request_units = optional(number, 0)
        max_write_request_units = optional(number, 0)
      }))
      warm_throughput = optional(object({
        read_units_per_second = optional(number, 0)
        write_units_per_second = optional(number, 0)
      }))
    })), [])
    local_secondary_indexes = optional(list(object({
      name = string
      range_key = string
      projection = object({
        type = string
        non_key_attributes = optional(list(string), [])
      })
    })), [])
    ttl = optional(object({
      enabled = optional(bool, false)
      attribute_name = optional(string, "")
    }))
    stream_enabled = optional(bool, false)
    stream_view_type = optional(string, "")
    point_in_time_recovery = optional(object({
      enabled = optional(bool, false)
      recovery_period_in_days = optional(number, 0)
    }))
    server_side_encryption = optional(object({
      enabled = optional(bool, false)
      kms_key_arn = optional(string, "")
    }))
    table_class = optional(string, "")
    deletion_protection_enabled = optional(bool, false)
    contributor_insights = optional(object({
      enabled = optional(bool, false)
      mode = optional(string, "")
      gsi_index_names = optional(list(string), [])
    }))
    resource_policy = optional(string, "")
    kinesis_streaming_destination = optional(object({
      stream_arn = string
      approximate_creation_date_time_precision = optional(string, "")
    }))
    replicas = optional(list(object({
      region_name = string
      kms_key_arn = optional(string, "")
      point_in_time_recovery = optional(bool, false)
      deletion_protection_enabled = optional(bool, false)
      propagate_tags = optional(bool, false)
      consistency_mode = optional(string, "")
    })), [])
    global_table_witness = optional(object({
      region_name = string
    }))
    restore_source_name = optional(string, "")
    restore_source_table_arn = optional(string, "")
    restore_date_time = optional(string, "")
    restore_to_latest_time = optional(bool, false)
    restore_backup_arn = optional(string, "")
    import_table = optional(object({
      s3_bucket = string
      s3_bucket_owner = optional(string, "")
      s3_key_prefix = optional(string, "")
      input_format = string
      input_compression_type = optional(string, "")
      csv = optional(object({
        delimiter = optional(string, "")
        header_list = optional(list(string), [])
      }))
    }))
  })
}
