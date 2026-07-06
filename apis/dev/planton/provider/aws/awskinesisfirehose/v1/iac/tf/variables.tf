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
  description = "AwsKinesisFirehose specification"
  type = object({
    region = string
    kinesis_stream_source = optional(object({
      stream_arn = string
      role_arn = string
    }))
    sse_enabled = optional(bool, false)
    sse_kms_key_arn = optional(string, "")
    extended_s3 = optional(object({
      bucket_arn = string
      role_arn = string
      prefix = optional(string, "")
      error_output_prefix = optional(string, "")
      compression_format = optional(string, "")
      kms_key_arn = optional(string, "")
      buffering = optional(object({
        interval_in_seconds = optional(number, 0)
        size_in_mbs = optional(number, 0)
      }))
      custom_time_zone = optional(string, "")
      file_extension = optional(string, "")
      s3_backup_mode = optional(string, "")
      s3_backup = optional(object({
        bucket_arn = string
        role_arn = string
        prefix = optional(string, "")
        error_output_prefix = optional(string, "")
        compression_format = optional(string, "")
        kms_key_arn = optional(string, "")
        buffering = optional(object({
          interval_in_seconds = optional(number, 0)
          size_in_mbs = optional(number, 0)
        }))
        logging = optional(object({
          enabled = optional(bool, false)
          log_group_name = optional(string, "")
          log_stream_name = optional(string, "")
        }))
      }))
      processing = optional(object({
        enabled = optional(bool, false)
        lambda_arn = optional(string, "")
        buffer_size_in_mbs = optional(number, 0)
        buffer_interval_in_seconds = optional(number, 0)
        number_of_retries = optional(number, 0)
      }))
      logging = optional(object({
        enabled = optional(bool, false)
        log_group_name = optional(string, "")
        log_stream_name = optional(string, "")
      }))
      dynamic_partitioning = optional(object({
        enabled = optional(bool, false)
        retry_duration_in_seconds = optional(number, 0)
      }))
      data_format_conversion = optional(object({
        enabled = optional(bool, false)
        input_format = optional(string, "")
        output_format = optional(string, "")
        parquet_compression = optional(string, "")
        orc_compression = optional(string, "")
        schema = optional(object({
          database_name = string
          table_name = string
          role_arn = string
          catalog_id = optional(string, "")
          region = optional(string, "")
          version_id = optional(string, "")
        }))
      }))
    }))
    opensearch = optional(object({
      domain_arn = optional(string, "")
      cluster_endpoint = optional(string, "")
      index_name = string
      role_arn = string
      index_rotation_period = optional(string, "")
      type_name = optional(string, "")
      buffering = optional(object({
        interval_in_seconds = optional(number, 0)
        size_in_mbs = optional(number, 0)
      }))
      retry_duration_in_seconds = optional(number, 0)
      s3_backup_mode = optional(string, "")
      s3_config = object({
        bucket_arn = string
        role_arn = string
        prefix = optional(string, "")
        error_output_prefix = optional(string, "")
        compression_format = optional(string, "")
        kms_key_arn = optional(string, "")
        buffering = optional(object({
          interval_in_seconds = optional(number, 0)
          size_in_mbs = optional(number, 0)
        }))
        logging = optional(object({
          enabled = optional(bool, false)
          log_group_name = optional(string, "")
          log_stream_name = optional(string, "")
        }))
      })
      processing = optional(object({
        enabled = optional(bool, false)
        lambda_arn = optional(string, "")
        buffer_size_in_mbs = optional(number, 0)
        buffer_interval_in_seconds = optional(number, 0)
        number_of_retries = optional(number, 0)
      }))
      logging = optional(object({
        enabled = optional(bool, false)
        log_group_name = optional(string, "")
        log_stream_name = optional(string, "")
      }))
      vpc_config = optional(object({
        subnet_ids = optional(list(string), [])
        security_group_ids = optional(list(string), [])
        role_arn = string
      }))
    }))
    http_endpoint = optional(object({
      url = string
      name = optional(string, "")
      access_key = optional(string, "")
      role_arn = optional(string, "")
      buffering = optional(object({
        interval_in_seconds = optional(number, 0)
        size_in_mbs = optional(number, 0)
      }))
      retry_duration_in_seconds = optional(number, 0)
      s3_backup_mode = optional(string, "")
      s3_config = object({
        bucket_arn = string
        role_arn = string
        prefix = optional(string, "")
        error_output_prefix = optional(string, "")
        compression_format = optional(string, "")
        kms_key_arn = optional(string, "")
        buffering = optional(object({
          interval_in_seconds = optional(number, 0)
          size_in_mbs = optional(number, 0)
        }))
        logging = optional(object({
          enabled = optional(bool, false)
          log_group_name = optional(string, "")
          log_stream_name = optional(string, "")
        }))
      })
      processing = optional(object({
        enabled = optional(bool, false)
        lambda_arn = optional(string, "")
        buffer_size_in_mbs = optional(number, 0)
        buffer_interval_in_seconds = optional(number, 0)
        number_of_retries = optional(number, 0)
      }))
      logging = optional(object({
        enabled = optional(bool, false)
        log_group_name = optional(string, "")
        log_stream_name = optional(string, "")
      }))
      request_config = optional(object({
        content_encoding = optional(string, "")
        common_attributes = optional(list(object({
          name = string
          value = string
        })), [])
      }))
    }))
    redshift = optional(object({
      cluster_jdbcurl = string
      role_arn = string
      data_table_name = string
      data_table_columns = optional(string, "")
      copy_options = optional(string, "")
      username = optional(string, "")
      password = optional(string, "")
      s3_config = object({
        bucket_arn = string
        role_arn = string
        prefix = optional(string, "")
        error_output_prefix = optional(string, "")
        compression_format = optional(string, "")
        kms_key_arn = optional(string, "")
        buffering = optional(object({
          interval_in_seconds = optional(number, 0)
          size_in_mbs = optional(number, 0)
        }))
        logging = optional(object({
          enabled = optional(bool, false)
          log_group_name = optional(string, "")
          log_stream_name = optional(string, "")
        }))
      })
      retry_duration_in_seconds = optional(number, 0)
      s3_backup_mode = optional(string, "")
      s3_backup = optional(object({
        bucket_arn = string
        role_arn = string
        prefix = optional(string, "")
        error_output_prefix = optional(string, "")
        compression_format = optional(string, "")
        kms_key_arn = optional(string, "")
        buffering = optional(object({
          interval_in_seconds = optional(number, 0)
          size_in_mbs = optional(number, 0)
        }))
        logging = optional(object({
          enabled = optional(bool, false)
          log_group_name = optional(string, "")
          log_stream_name = optional(string, "")
        }))
      }))
      processing = optional(object({
        enabled = optional(bool, false)
        lambda_arn = optional(string, "")
        buffer_size_in_mbs = optional(number, 0)
        buffer_interval_in_seconds = optional(number, 0)
        number_of_retries = optional(number, 0)
      }))
      logging = optional(object({
        enabled = optional(bool, false)
        log_group_name = optional(string, "")
        log_stream_name = optional(string, "")
      }))
    }))
  })
}
