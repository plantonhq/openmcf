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
    msk_source = optional(object({
      msk_cluster_arn = string
      topic_name = string
      connectivity = string
      role_arn = string
      read_from_timestamp = optional(string, "")
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
        processors = optional(list(object({
          lambda = optional(object({
            lambda_arn = string
            buffer_size_in_mbs = optional(number, 0)
            buffer_interval_in_seconds = optional(number, 0)
            number_of_retries = optional(number, 0)
            role_arn = optional(string, "")
          }))
          metadata_extraction = optional(object({
            query = string
            json_parsing_engine = optional(string, "")
          }))
          decompression = optional(object({
            compression_format = string
          }))
          cloudwatch_log_processing = optional(object({
            data_message_extraction = optional(bool, false)
          }))
          append_delimiter = optional(object({
            delimiter = string
          }))
          record_deaggregation = optional(object({
            sub_record_type = string
            delimiter = optional(string, "")
          }))
        })), [])
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
        open_x_json = optional(object({
          case_insensitive = optional(bool)
          column_to_json_key_mappings = optional(map(string), {})
          convert_dots_in_json_keys_to_underscores = optional(bool, false)
        }))
        hive_json = optional(object({
          timestamp_formats = optional(list(string), [])
        }))
        parquet = optional(object({
          compression = optional(string, "")
          block_size_bytes = optional(number, 0)
          page_size_bytes = optional(number, 0)
          max_padding_bytes = optional(number, 0)
          enable_dictionary_compression = optional(bool, false)
          writer_version = optional(string, "")
        }))
        orc = optional(object({
          compression = optional(string, "")
          block_size_bytes = optional(number, 0)
          stripe_size_bytes = optional(number, 0)
          bloom_filter_columns = optional(list(string), [])
          bloom_filter_false_positive_probability = optional(number)
          dictionary_key_threshold = optional(number, 0)
          enable_padding = optional(bool, false)
          padding_tolerance = optional(number)
          format_version = optional(string, "")
          row_index_stride = optional(number, 0)
        }))
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
      default_document_id_format = optional(string, "")
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
        processors = optional(list(object({
          lambda = optional(object({
            lambda_arn = string
            buffer_size_in_mbs = optional(number, 0)
            buffer_interval_in_seconds = optional(number, 0)
            number_of_retries = optional(number, 0)
            role_arn = optional(string, "")
          }))
          metadata_extraction = optional(object({
            query = string
            json_parsing_engine = optional(string, "")
          }))
          decompression = optional(object({
            compression_format = string
          }))
          cloudwatch_log_processing = optional(object({
            data_message_extraction = optional(bool, false)
          }))
          append_delimiter = optional(object({
            delimiter = string
          }))
          record_deaggregation = optional(object({
            sub_record_type = string
            delimiter = optional(string, "")
          }))
        })), [])
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
    opensearch_serverless = optional(object({
      collection_endpoint = string
      index_name = string
      role_arn = string
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
        processors = optional(list(object({
          lambda = optional(object({
            lambda_arn = string
            buffer_size_in_mbs = optional(number, 0)
            buffer_interval_in_seconds = optional(number, 0)
            number_of_retries = optional(number, 0)
            role_arn = optional(string, "")
          }))
          metadata_extraction = optional(object({
            query = string
            json_parsing_engine = optional(string, "")
          }))
          decompression = optional(object({
            compression_format = string
          }))
          cloudwatch_log_processing = optional(object({
            data_message_extraction = optional(bool, false)
          }))
          append_delimiter = optional(object({
            delimiter = string
          }))
          record_deaggregation = optional(object({
            sub_record_type = string
            delimiter = optional(string, "")
          }))
        })), [])
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
      secrets_manager = optional(object({
        secret_arn = string
        role_arn = optional(string, "")
      }))
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
        processors = optional(list(object({
          lambda = optional(object({
            lambda_arn = string
            buffer_size_in_mbs = optional(number, 0)
            buffer_interval_in_seconds = optional(number, 0)
            number_of_retries = optional(number, 0)
            role_arn = optional(string, "")
          }))
          metadata_extraction = optional(object({
            query = string
            json_parsing_engine = optional(string, "")
          }))
          decompression = optional(object({
            compression_format = string
          }))
          cloudwatch_log_processing = optional(object({
            data_message_extraction = optional(bool, false)
          }))
          append_delimiter = optional(object({
            delimiter = string
          }))
          record_deaggregation = optional(object({
            sub_record_type = string
            delimiter = optional(string, "")
          }))
        })), [])
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
      secrets_manager = optional(object({
        secret_arn = string
        role_arn = optional(string, "")
      }))
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
        processors = optional(list(object({
          lambda = optional(object({
            lambda_arn = string
            buffer_size_in_mbs = optional(number, 0)
            buffer_interval_in_seconds = optional(number, 0)
            number_of_retries = optional(number, 0)
            role_arn = optional(string, "")
          }))
          metadata_extraction = optional(object({
            query = string
            json_parsing_engine = optional(string, "")
          }))
          decompression = optional(object({
            compression_format = string
          }))
          cloudwatch_log_processing = optional(object({
            data_message_extraction = optional(bool, false)
          }))
          append_delimiter = optional(object({
            delimiter = string
          }))
          record_deaggregation = optional(object({
            sub_record_type = string
            delimiter = optional(string, "")
          }))
        })), [])
      }))
      logging = optional(object({
        enabled = optional(bool, false)
        log_group_name = optional(string, "")
        log_stream_name = optional(string, "")
      }))
    }))
    splunk = optional(object({
      hec_endpoint = string
      hec_endpoint_type = optional(string, "")
      hec_token = optional(string, "")
      secrets_manager = optional(object({
        secret_arn = string
        role_arn = optional(string, "")
      }))
      hec_acknowledgment_timeout_in_seconds = optional(number, 0)
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
        processors = optional(list(object({
          lambda = optional(object({
            lambda_arn = string
            buffer_size_in_mbs = optional(number, 0)
            buffer_interval_in_seconds = optional(number, 0)
            number_of_retries = optional(number, 0)
            role_arn = optional(string, "")
          }))
          metadata_extraction = optional(object({
            query = string
            json_parsing_engine = optional(string, "")
          }))
          decompression = optional(object({
            compression_format = string
          }))
          cloudwatch_log_processing = optional(object({
            data_message_extraction = optional(bool, false)
          }))
          append_delimiter = optional(object({
            delimiter = string
          }))
          record_deaggregation = optional(object({
            sub_record_type = string
            delimiter = optional(string, "")
          }))
        })), [])
      }))
      logging = optional(object({
        enabled = optional(bool, false)
        log_group_name = optional(string, "")
        log_stream_name = optional(string, "")
      }))
    }))
    snowflake = optional(object({
      account_url = string
      database = string
      schema = string
      table = string
      role_arn = string
      user = optional(string, "")
      private_key = optional(string, "")
      key_passphrase = optional(string, "")
      secrets_manager = optional(object({
        secret_arn = string
        role_arn = optional(string, "")
      }))
      data_loading_option = optional(string, "")
      content_column_name = optional(string, "")
      metadata_column_name = optional(string, "")
      snowflake_role = optional(string, "")
      private_link_vpce_id = optional(string, "")
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
        processors = optional(list(object({
          lambda = optional(object({
            lambda_arn = string
            buffer_size_in_mbs = optional(number, 0)
            buffer_interval_in_seconds = optional(number, 0)
            number_of_retries = optional(number, 0)
            role_arn = optional(string, "")
          }))
          metadata_extraction = optional(object({
            query = string
            json_parsing_engine = optional(string, "")
          }))
          decompression = optional(object({
            compression_format = string
          }))
          cloudwatch_log_processing = optional(object({
            data_message_extraction = optional(bool, false)
          }))
          append_delimiter = optional(object({
            delimiter = string
          }))
          record_deaggregation = optional(object({
            sub_record_type = string
            delimiter = optional(string, "")
          }))
        })), [])
      }))
      logging = optional(object({
        enabled = optional(bool, false)
        log_group_name = optional(string, "")
        log_stream_name = optional(string, "")
      }))
    }))
    iceberg = optional(object({
      catalog_arn = string
      role_arn = string
      destination_tables = optional(list(object({
        database_name = string
        table_name = string
        s3_error_output_prefix = optional(string, "")
        unique_keys = optional(list(string), [])
      })), [])
      append_only = optional(bool, false)
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
        processors = optional(list(object({
          lambda = optional(object({
            lambda_arn = string
            buffer_size_in_mbs = optional(number, 0)
            buffer_interval_in_seconds = optional(number, 0)
            number_of_retries = optional(number, 0)
            role_arn = optional(string, "")
          }))
          metadata_extraction = optional(object({
            query = string
            json_parsing_engine = optional(string, "")
          }))
          decompression = optional(object({
            compression_format = string
          }))
          cloudwatch_log_processing = optional(object({
            data_message_extraction = optional(bool, false)
          }))
          append_delimiter = optional(object({
            delimiter = string
          }))
          record_deaggregation = optional(object({
            sub_record_type = string
            delimiter = optional(string, "")
          }))
        })), [])
      }))
      logging = optional(object({
        enabled = optional(bool, false)
        log_group_name = optional(string, "")
        log_stream_name = optional(string, "")
      }))
    }))
  })
}
