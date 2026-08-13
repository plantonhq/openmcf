# ---------------------------------------------------------------------------
# Kinesis Data Firehose delivery stream
# ---------------------------------------------------------------------------
# Exactly one destination block renders (the proto oneof guarantees one arm;
# local.destination_type selects it). The typed processor pipeline from the
# spec is normalized ONCE in locals.tf (local.processors_normalized) and every
# destination renders the same list — keeping the spec's intent-level model
# while sending AWS its raw {type, parameters[]} shape.

resource "aws_kinesis_firehose_delivery_stream" "this" {
  name        = local.delivery_stream_name
  destination = local.destination_type
  tags        = local.aws_tags

  # ---------------------------------------------------------------------------
  # Source configuration (optional — Direct PUT when absent)
  # ---------------------------------------------------------------------------

  dynamic "kinesis_source_configuration" {
    for_each = local.has_kinesis_source ? [var.spec.kinesis_stream_source] : []
    iterator = src
    content {
      kinesis_stream_arn = src.value.stream_arn
      role_arn           = src.value.role_arn
    }
  }

  # The MSK source (whole block ForceNew). Connectivity + role live in the
  # nested authentication_configuration; read_from_timestamp rewinds the
  # topic to a point in time at creation.
  dynamic "msk_source_configuration" {
    for_each = local.has_msk_source ? [var.spec.msk_source] : []
    iterator = src
    content {
      msk_cluster_arn     = src.value.msk_cluster_arn
      topic_name          = src.value.topic_name
      read_from_timestamp = try(src.value.read_from_timestamp, null) != "" ? src.value.read_from_timestamp : null

      authentication_configuration {
        connectivity = src.value.connectivity
        role_arn     = src.value.role_arn
      }
    }
  }

  # ---------------------------------------------------------------------------
  # Server-side encryption (Direct PUT only)
  # ---------------------------------------------------------------------------

  dynamic "server_side_encryption" {
    for_each = local.sse_enabled ? [1] : []
    content {
      enabled  = true
      key_type = local.sse_key_type
      key_arn  = local.sse_kms_key_arn
    }
  }

  # ===========================================================================
  # Extended S3 destination
  # ===========================================================================

  dynamic "extended_s3_configuration" {
    for_each = local.destination_type == "extended_s3" ? [var.spec.extended_s3] : []
    iterator = dest
    content {
      # Required fields
      bucket_arn = dest.value.bucket_arn
      role_arn   = dest.value.role_arn

      # S3 delivery options
      prefix              = try(dest.value.prefix, null) != "" ? dest.value.prefix : null
      error_output_prefix = try(dest.value.error_output_prefix, null) != "" ? dest.value.error_output_prefix : null
      compression_format  = try(dest.value.compression_format, null) != "" ? dest.value.compression_format : null
      kms_key_arn         = dest.value.kms_key_arn != "" ? dest.value.kms_key_arn : null
      buffering_interval  = try(dest.value.buffering.interval_in_seconds, 0) > 0 ? dest.value.buffering.interval_in_seconds : null
      buffering_size      = try(dest.value.buffering.size_in_mbs, 0) > 0 ? dest.value.buffering.size_in_mbs : null
      custom_time_zone    = try(dest.value.custom_time_zone, null) != "" ? dest.value.custom_time_zone : null
      file_extension      = try(dest.value.file_extension, null) != "" ? dest.value.file_extension : null

      # S3 backup mode
      s3_backup_mode = try(dest.value.s3_backup_mode, null) != "" ? dest.value.s3_backup_mode : null

      # --- S3 backup configuration (source record backup) ---

      dynamic "s3_backup_configuration" {
        for_each = try(dest.value.s3_backup, null) != null ? [dest.value.s3_backup] : []
        iterator = bkp
        content {
          bucket_arn          = bkp.value.bucket_arn
          role_arn            = bkp.value.role_arn
          prefix              = try(bkp.value.prefix, null) != "" ? bkp.value.prefix : null
          error_output_prefix = try(bkp.value.error_output_prefix, null) != "" ? bkp.value.error_output_prefix : null
          compression_format  = try(bkp.value.compression_format, null) != "" ? bkp.value.compression_format : null
          kms_key_arn         = bkp.value.kms_key_arn != "" ? bkp.value.kms_key_arn : null
          buffering_interval  = try(bkp.value.buffering.interval_in_seconds, 0) > 0 ? bkp.value.buffering.interval_in_seconds : null
          buffering_size      = try(bkp.value.buffering.size_in_mbs, 0) > 0 ? bkp.value.buffering.size_in_mbs : null

          # CloudWatch logging for the backup S3 leg — distinct from the
          # destination-level logging block.
          dynamic "cloudwatch_logging_options" {
            for_each = try(bkp.value.logging.enabled, false) ? [bkp.value.logging] : []
            iterator = s3log
            content {
              enabled         = true
              log_group_name  = s3log.value.log_group_name
              log_stream_name = s3log.value.log_stream_name
            }
          }
        }
      }

      # --- Processing pipeline (normalized typed processors) ---

      dynamic "processing_configuration" {
        for_each = local.processing_enabled ? [1] : []
        content {
          enabled = true

          dynamic "processors" {
            for_each = local.processors_normalized
            content {
              type = processors.value.type

              dynamic "parameters" {
                for_each = processors.value.parameters
                content {
                  parameter_name  = parameters.value.name
                  parameter_value = parameters.value.value
                }
              }
            }
          }
        }
      }

      # --- CloudWatch logging ---

      dynamic "cloudwatch_logging_options" {
        for_each = try(dest.value.logging.enabled, false) ? [dest.value.logging] : []
        iterator = log
        content {
          enabled         = true
          log_group_name  = log.value.log_group_name
          log_stream_name = log.value.log_stream_name
        }
      }

      # --- Dynamic partitioning (ForceNew) ---

      dynamic "dynamic_partitioning_configuration" {
        for_each = try(dest.value.dynamic_partitioning.enabled, false) ? [dest.value.dynamic_partitioning] : []
        iterator = dp
        content {
          enabled        = true
          retry_duration = try(dp.value.retry_duration_in_seconds, 0) > 0 ? dp.value.retry_duration_in_seconds : null
        }
      }

      # --- Data format conversion (Parquet/ORC via Glue catalog) ---

      dynamic "data_format_conversion_configuration" {
        for_each = try(dest.value.data_format_conversion.enabled, false) ? [dest.value.data_format_conversion] : []
        iterator = dfc
        content {
          enabled = true

          # Input format (deserializer) — the spec enforces exactly one arm.
          # Unset optional leaves are omitted so the AWS defaults win.
          input_format_configuration {
            deserializer {
              dynamic "open_x_json_ser_de" {
                for_each = try(dfc.value.open_x_json, null) != null ? [dfc.value.open_x_json] : []
                iterator = oxj
                content {
                  # AWS defaults case_insensitive to true; a null passes
                  # through as "omitted" so only an explicit false is sent.
                  case_insensitive                         = try(oxj.value.case_insensitive, null)
                  column_to_json_key_mappings              = length(try(oxj.value.column_to_json_key_mappings, {})) > 0 ? oxj.value.column_to_json_key_mappings : null
                  convert_dots_in_json_keys_to_underscores = try(oxj.value.convert_dots_in_json_keys_to_underscores, false) ? true : null
                }
              }
              dynamic "hive_json_ser_de" {
                for_each = try(dfc.value.hive_json, null) != null ? [dfc.value.hive_json] : []
                iterator = hj
                content {
                  timestamp_formats = length(try(hj.value.timestamp_formats, [])) > 0 ? hj.value.timestamp_formats : null
                }
              }
            }
          }

          # Output format (serializer) — the spec enforces exactly one arm.
          output_format_configuration {
            serializer {
              dynamic "parquet_ser_de" {
                for_each = try(dfc.value.parquet, null) != null ? [dfc.value.parquet] : []
                iterator = pq
                content {
                  compression                   = try(pq.value.compression, null) != "" ? pq.value.compression : null
                  block_size_bytes              = try(pq.value.block_size_bytes, 0) > 0 ? pq.value.block_size_bytes : null
                  page_size_bytes               = try(pq.value.page_size_bytes, 0) > 0 ? pq.value.page_size_bytes : null
                  max_padding_bytes             = try(pq.value.max_padding_bytes, 0) > 0 ? pq.value.max_padding_bytes : null
                  enable_dictionary_compression = try(pq.value.enable_dictionary_compression, false) ? true : null
                  writer_version                = try(pq.value.writer_version, null) != "" ? pq.value.writer_version : null
                }
              }
              dynamic "orc_ser_de" {
                for_each = try(dfc.value.orc, null) != null ? [dfc.value.orc] : []
                iterator = orc
                content {
                  compression          = try(orc.value.compression, null) != "" ? orc.value.compression : null
                  block_size_bytes     = try(orc.value.block_size_bytes, 0) > 0 ? orc.value.block_size_bytes : null
                  stripe_size_bytes    = try(orc.value.stripe_size_bytes, 0) > 0 ? orc.value.stripe_size_bytes : null
                  bloom_filter_columns = length(try(orc.value.bloom_filter_columns, [])) > 0 ? orc.value.bloom_filter_columns : null
                  # Optional-with-presence in the spec: explicit 0 is
                  # AWS-legal here, distinct from "unset -> AWS default".
                  bloom_filter_false_positive_probability = try(orc.value.bloom_filter_false_positive_probability, null)
                  dictionary_key_threshold                = try(orc.value.dictionary_key_threshold, 0) > 0 ? orc.value.dictionary_key_threshold : null
                  enable_padding                          = try(orc.value.enable_padding, false) ? true : null
                  padding_tolerance                       = try(orc.value.padding_tolerance, null)
                  format_version                          = try(orc.value.format_version, null) != "" ? orc.value.format_version : null
                  row_index_stride                        = try(orc.value.row_index_stride, 0) > 0 ? orc.value.row_index_stride : null
                }
              }
            }
          }

          # Glue Data Catalog schema
          schema_configuration {
            database_name = dfc.value.schema.database_name
            table_name    = dfc.value.schema.table_name
            role_arn      = dfc.value.schema.role_arn
            catalog_id    = try(dfc.value.schema.catalog_id, null) != "" ? dfc.value.schema.catalog_id : null
            region        = try(dfc.value.schema.region, null) != "" ? dfc.value.schema.region : null
            version_id    = try(dfc.value.schema.version_id, null) != "" ? dfc.value.schema.version_id : null
          }
        }
      }
    }
  }

  # ===========================================================================
  # OpenSearch destination
  # ===========================================================================

  dynamic "opensearch_configuration" {
    for_each = local.destination_type == "opensearch" ? [var.spec.opensearch] : []
    iterator = dest
    content {
      # Target — exactly one of domain_arn or cluster_endpoint
      domain_arn       = dest.value.domain_arn != "" ? dest.value.domain_arn : null
      cluster_endpoint = try(dest.value.cluster_endpoint, null) != "" ? dest.value.cluster_endpoint : null

      # Indexing configuration
      index_name            = dest.value.index_name
      role_arn              = dest.value.role_arn
      index_rotation_period = try(dest.value.index_rotation_period, null) != "" ? dest.value.index_rotation_period : null
      type_name             = try(dest.value.type_name, null) != "" ? dest.value.type_name : null

      # Delivery configuration
      buffering_interval = try(dest.value.buffering.interval_in_seconds, 0) > 0 ? dest.value.buffering.interval_in_seconds : null
      buffering_size     = try(dest.value.buffering.size_in_mbs, 0) > 0 ? dest.value.buffering.size_in_mbs : null
      retry_duration     = try(dest.value.retry_duration_in_seconds, 0) > 0 ? dest.value.retry_duration_in_seconds : null

      # S3 backup mode (ForceNew on OpenSearch destinations)
      s3_backup_mode = try(dest.value.s3_backup_mode, null) != "" ? dest.value.s3_backup_mode : null

      # --- Document ID assignment ---

      dynamic "document_id_options" {
        for_each = try(dest.value.default_document_id_format, null) != "" ? [dest.value.default_document_id_format] : []
        content {
          default_document_id_format = document_id_options.value
        }
      }

      # --- S3 config (required — backs up failed/all documents) ---

      s3_configuration {
        bucket_arn          = dest.value.s3_config.bucket_arn
        role_arn            = dest.value.s3_config.role_arn
        prefix              = try(dest.value.s3_config.prefix, null) != "" ? dest.value.s3_config.prefix : null
        error_output_prefix = try(dest.value.s3_config.error_output_prefix, null) != "" ? dest.value.s3_config.error_output_prefix : null
        compression_format  = try(dest.value.s3_config.compression_format, null) != "" ? dest.value.s3_config.compression_format : null
        kms_key_arn         = dest.value.s3_config.kms_key_arn != "" ? dest.value.s3_config.kms_key_arn : null
        buffering_interval  = try(dest.value.s3_config.buffering.interval_in_seconds, 0) > 0 ? dest.value.s3_config.buffering.interval_in_seconds : null
        buffering_size      = try(dest.value.s3_config.buffering.size_in_mbs, 0) > 0 ? dest.value.s3_config.buffering.size_in_mbs : null

        # CloudWatch logging for THIS S3 leg (backup/staging errors) —
        # distinct from the destination-level logging block below.
        dynamic "cloudwatch_logging_options" {
          for_each = try(dest.value.s3_config.logging.enabled, false) ? [dest.value.s3_config.logging] : []
          iterator = s3log
          content {
            enabled         = true
            log_group_name  = s3log.value.log_group_name
            log_stream_name = s3log.value.log_stream_name
          }
        }
      }

      # --- Processing pipeline (normalized typed processors) ---

      dynamic "processing_configuration" {
        for_each = local.processing_enabled ? [1] : []
        content {
          enabled = true

          dynamic "processors" {
            for_each = local.processors_normalized
            content {
              type = processors.value.type

              dynamic "parameters" {
                for_each = processors.value.parameters
                content {
                  parameter_name  = parameters.value.name
                  parameter_value = parameters.value.value
                }
              }
            }
          }
        }
      }

      # --- CloudWatch logging ---

      dynamic "cloudwatch_logging_options" {
        for_each = try(dest.value.logging.enabled, false) ? [dest.value.logging] : []
        iterator = log
        content {
          enabled         = true
          log_group_name  = log.value.log_group_name
          log_stream_name = log.value.log_stream_name
        }
      }

      # --- VPC config (for VPC-deployed OpenSearch domains, ForceNew) ---

      dynamic "vpc_config" {
        for_each = try(dest.value.vpc_config, null) != null ? [dest.value.vpc_config] : []
        iterator = vpc
        content {
          role_arn           = vpc.value.role_arn
          subnet_ids         = vpc.value.subnet_ids
          security_group_ids = vpc.value.security_group_ids
        }
      }
    }
  }

  # ===========================================================================
  # OpenSearch Serverless destination
  # ===========================================================================

  dynamic "opensearchserverless_configuration" {
    for_each = local.destination_type == "opensearchserverless" ? [var.spec.opensearch_serverless] : []
    iterator = dest
    content {
      collection_endpoint = dest.value.collection_endpoint
      index_name          = dest.value.index_name
      role_arn            = dest.value.role_arn

      # Delivery configuration
      buffering_interval = try(dest.value.buffering.interval_in_seconds, 0) > 0 ? dest.value.buffering.interval_in_seconds : null
      buffering_size     = try(dest.value.buffering.size_in_mbs, 0) > 0 ? dest.value.buffering.size_in_mbs : null
      retry_duration     = try(dest.value.retry_duration_in_seconds, 0) > 0 ? dest.value.retry_duration_in_seconds : null

      # S3 backup mode (ForceNew)
      s3_backup_mode = try(dest.value.s3_backup_mode, null) != "" ? dest.value.s3_backup_mode : null

      # --- S3 config (required — backs up failed/all documents) ---

      s3_configuration {
        bucket_arn          = dest.value.s3_config.bucket_arn
        role_arn            = dest.value.s3_config.role_arn
        prefix              = try(dest.value.s3_config.prefix, null) != "" ? dest.value.s3_config.prefix : null
        error_output_prefix = try(dest.value.s3_config.error_output_prefix, null) != "" ? dest.value.s3_config.error_output_prefix : null
        compression_format  = try(dest.value.s3_config.compression_format, null) != "" ? dest.value.s3_config.compression_format : null
        kms_key_arn         = dest.value.s3_config.kms_key_arn != "" ? dest.value.s3_config.kms_key_arn : null
        buffering_interval  = try(dest.value.s3_config.buffering.interval_in_seconds, 0) > 0 ? dest.value.s3_config.buffering.interval_in_seconds : null
        buffering_size      = try(dest.value.s3_config.buffering.size_in_mbs, 0) > 0 ? dest.value.s3_config.buffering.size_in_mbs : null

        # CloudWatch logging for THIS S3 leg (backup/staging errors) —
        # distinct from the destination-level logging block below.
        dynamic "cloudwatch_logging_options" {
          for_each = try(dest.value.s3_config.logging.enabled, false) ? [dest.value.s3_config.logging] : []
          iterator = s3log
          content {
            enabled         = true
            log_group_name  = s3log.value.log_group_name
            log_stream_name = s3log.value.log_stream_name
          }
        }
      }

      # --- Processing pipeline (normalized typed processors) ---

      dynamic "processing_configuration" {
        for_each = local.processing_enabled ? [1] : []
        content {
          enabled = true

          dynamic "processors" {
            for_each = local.processors_normalized
            content {
              type = processors.value.type

              dynamic "parameters" {
                for_each = processors.value.parameters
                content {
                  parameter_name  = parameters.value.name
                  parameter_value = parameters.value.value
                }
              }
            }
          }
        }
      }

      # --- CloudWatch logging ---

      dynamic "cloudwatch_logging_options" {
        for_each = try(dest.value.logging.enabled, false) ? [dest.value.logging] : []
        iterator = log
        content {
          enabled         = true
          log_group_name  = log.value.log_group_name
          log_stream_name = log.value.log_stream_name
        }
      }

      # --- VPC config (ForceNew) ---

      dynamic "vpc_config" {
        for_each = try(dest.value.vpc_config, null) != null ? [dest.value.vpc_config] : []
        iterator = vpc
        content {
          role_arn           = vpc.value.role_arn
          subnet_ids         = vpc.value.subnet_ids
          security_group_ids = vpc.value.security_group_ids
        }
      }
    }
  }

  # ===========================================================================
  # HTTP Endpoint destination
  # ===========================================================================

  dynamic "http_endpoint_configuration" {
    for_each = local.destination_type == "http_endpoint" ? [var.spec.http_endpoint] : []
    iterator = dest
    content {
      # Endpoint configuration
      url        = dest.value.url
      name       = try(dest.value.name, null) != "" ? dest.value.name : null
      access_key = try(dest.value.access_key, null) != "" ? dest.value.access_key : null
      role_arn   = dest.value.role_arn != "" ? dest.value.role_arn : null

      # Delivery configuration
      buffering_interval = try(dest.value.buffering.interval_in_seconds, 0) > 0 ? dest.value.buffering.interval_in_seconds : null
      buffering_size     = try(dest.value.buffering.size_in_mbs, 0) > 0 ? dest.value.buffering.size_in_mbs : null
      retry_duration     = try(dest.value.retry_duration_in_seconds, 0) > 0 ? dest.value.retry_duration_in_seconds : null

      # S3 backup mode
      s3_backup_mode = try(dest.value.s3_backup_mode, null) != "" ? dest.value.s3_backup_mode : null

      # --- Secrets Manager credentials (XOR with access_key, ForceNew switch) ---

      dynamic "secrets_manager_configuration" {
        for_each = try(dest.value.secrets_manager, null) != null ? [dest.value.secrets_manager] : []
        iterator = sm
        content {
          enabled    = true
          secret_arn = sm.value.secret_arn
          role_arn   = sm.value.role_arn != "" ? sm.value.role_arn : null
        }
      }

      # --- S3 config (required — backs up failed/all records) ---

      s3_configuration {
        bucket_arn          = dest.value.s3_config.bucket_arn
        role_arn            = dest.value.s3_config.role_arn
        prefix              = try(dest.value.s3_config.prefix, null) != "" ? dest.value.s3_config.prefix : null
        error_output_prefix = try(dest.value.s3_config.error_output_prefix, null) != "" ? dest.value.s3_config.error_output_prefix : null
        compression_format  = try(dest.value.s3_config.compression_format, null) != "" ? dest.value.s3_config.compression_format : null
        kms_key_arn         = dest.value.s3_config.kms_key_arn != "" ? dest.value.s3_config.kms_key_arn : null
        buffering_interval  = try(dest.value.s3_config.buffering.interval_in_seconds, 0) > 0 ? dest.value.s3_config.buffering.interval_in_seconds : null
        buffering_size      = try(dest.value.s3_config.buffering.size_in_mbs, 0) > 0 ? dest.value.s3_config.buffering.size_in_mbs : null

        # CloudWatch logging for THIS S3 leg (backup/staging errors) —
        # distinct from the destination-level logging block below.
        dynamic "cloudwatch_logging_options" {
          for_each = try(dest.value.s3_config.logging.enabled, false) ? [dest.value.s3_config.logging] : []
          iterator = s3log
          content {
            enabled         = true
            log_group_name  = s3log.value.log_group_name
            log_stream_name = s3log.value.log_stream_name
          }
        }
      }

      # --- Processing pipeline (normalized typed processors) ---

      dynamic "processing_configuration" {
        for_each = local.processing_enabled ? [1] : []
        content {
          enabled = true

          dynamic "processors" {
            for_each = local.processors_normalized
            content {
              type = processors.value.type

              dynamic "parameters" {
                for_each = processors.value.parameters
                content {
                  parameter_name  = parameters.value.name
                  parameter_value = parameters.value.value
                }
              }
            }
          }
        }
      }

      # --- CloudWatch logging ---

      dynamic "cloudwatch_logging_options" {
        for_each = try(dest.value.logging.enabled, false) ? [dest.value.logging] : []
        iterator = log
        content {
          enabled         = true
          log_group_name  = log.value.log_group_name
          log_stream_name = log.value.log_stream_name
        }
      }

      # --- Request configuration ---

      dynamic "request_configuration" {
        for_each = try(dest.value.request_config, null) != null ? [dest.value.request_config] : []
        iterator = rc
        content {
          content_encoding = try(rc.value.content_encoding, null) != "" ? rc.value.content_encoding : null

          dynamic "common_attributes" {
            for_each = try(rc.value.common_attributes, [])
            content {
              name  = common_attributes.value.name
              value = common_attributes.value.value
            }
          }
        }
      }
    }
  }

  # ===========================================================================
  # Redshift destination
  # ===========================================================================
  # Provider quirk (documented, not visible in code): on updates the provider
  # force-nulls error_output_prefix on both the staging and backup S3 blocks —
  # the Redshift COPY API rejects it. Avoid relying on error prefixes here.

  dynamic "redshift_configuration" {
    for_each = local.destination_type == "redshift" ? [var.spec.redshift] : []
    iterator = dest
    content {
      # Redshift target
      cluster_jdbcurl    = dest.value.cluster_jdbcurl
      role_arn           = dest.value.role_arn
      data_table_name    = dest.value.data_table_name
      data_table_columns = try(dest.value.data_table_columns, null) != "" ? dest.value.data_table_columns : null
      copy_options       = try(dest.value.copy_options, null) != "" ? dest.value.copy_options : null

      # Authentication — plaintext pair XOR Secrets Manager (CEL-enforced)
      username = try(dest.value.username, null) != "" ? dest.value.username : null
      password = try(dest.value.password, null) != "" ? dest.value.password : null

      # Delivery configuration
      retry_duration = try(dest.value.retry_duration_in_seconds, 0) > 0 ? dest.value.retry_duration_in_seconds : null

      # S3 backup mode for source records
      s3_backup_mode = try(dest.value.s3_backup_mode, null) != "" ? dest.value.s3_backup_mode : null

      # --- Secrets Manager credentials (ForceNew switch) ---

      dynamic "secrets_manager_configuration" {
        for_each = try(dest.value.secrets_manager, null) != null ? [dest.value.secrets_manager] : []
        iterator = sm
        content {
          enabled    = true
          secret_arn = sm.value.secret_arn
          role_arn   = sm.value.role_arn != "" ? sm.value.role_arn : null
        }
      }

      # --- S3 intermediate staging config (required for Redshift COPY) ---

      s3_configuration {
        bucket_arn          = dest.value.s3_config.bucket_arn
        role_arn            = dest.value.s3_config.role_arn
        prefix              = try(dest.value.s3_config.prefix, null) != "" ? dest.value.s3_config.prefix : null
        error_output_prefix = try(dest.value.s3_config.error_output_prefix, null) != "" ? dest.value.s3_config.error_output_prefix : null
        compression_format  = try(dest.value.s3_config.compression_format, null) != "" ? dest.value.s3_config.compression_format : null
        kms_key_arn         = dest.value.s3_config.kms_key_arn != "" ? dest.value.s3_config.kms_key_arn : null
        buffering_interval  = try(dest.value.s3_config.buffering.interval_in_seconds, 0) > 0 ? dest.value.s3_config.buffering.interval_in_seconds : null
        buffering_size      = try(dest.value.s3_config.buffering.size_in_mbs, 0) > 0 ? dest.value.s3_config.buffering.size_in_mbs : null

        # CloudWatch logging for THIS S3 leg (backup/staging errors) —
        # distinct from the destination-level logging block below.
        dynamic "cloudwatch_logging_options" {
          for_each = try(dest.value.s3_config.logging.enabled, false) ? [dest.value.s3_config.logging] : []
          iterator = s3log
          content {
            enabled         = true
            log_group_name  = s3log.value.log_group_name
            log_stream_name = s3log.value.log_stream_name
          }
        }
      }

      # --- S3 backup configuration (source record backup) ---

      dynamic "s3_backup_configuration" {
        for_each = try(dest.value.s3_backup, null) != null ? [dest.value.s3_backup] : []
        iterator = bkp
        content {
          bucket_arn          = bkp.value.bucket_arn
          role_arn            = bkp.value.role_arn
          prefix              = try(bkp.value.prefix, null) != "" ? bkp.value.prefix : null
          error_output_prefix = try(bkp.value.error_output_prefix, null) != "" ? bkp.value.error_output_prefix : null
          compression_format  = try(bkp.value.compression_format, null) != "" ? bkp.value.compression_format : null
          kms_key_arn         = bkp.value.kms_key_arn != "" ? bkp.value.kms_key_arn : null
          buffering_interval  = try(bkp.value.buffering.interval_in_seconds, 0) > 0 ? bkp.value.buffering.interval_in_seconds : null
          buffering_size      = try(bkp.value.buffering.size_in_mbs, 0) > 0 ? bkp.value.buffering.size_in_mbs : null

          # CloudWatch logging for the backup S3 leg — distinct from the
          # destination-level logging block.
          dynamic "cloudwatch_logging_options" {
            for_each = try(bkp.value.logging.enabled, false) ? [bkp.value.logging] : []
            iterator = s3log
            content {
              enabled         = true
              log_group_name  = s3log.value.log_group_name
              log_stream_name = s3log.value.log_stream_name
            }
          }
        }
      }

      # --- Processing pipeline (normalized typed processors) ---

      dynamic "processing_configuration" {
        for_each = local.processing_enabled ? [1] : []
        content {
          enabled = true

          dynamic "processors" {
            for_each = local.processors_normalized
            content {
              type = processors.value.type

              dynamic "parameters" {
                for_each = processors.value.parameters
                content {
                  parameter_name  = parameters.value.name
                  parameter_value = parameters.value.value
                }
              }
            }
          }
        }
      }

      # --- CloudWatch logging ---

      dynamic "cloudwatch_logging_options" {
        for_each = try(dest.value.logging.enabled, false) ? [dest.value.logging] : []
        iterator = log
        content {
          enabled         = true
          log_group_name  = log.value.log_group_name
          log_stream_name = log.value.log_stream_name
        }
      }
    }
  }

  # ===========================================================================
  # Splunk destination
  # ===========================================================================

  dynamic "splunk_configuration" {
    for_each = local.destination_type == "splunk" ? [var.spec.splunk] : []
    iterator = dest
    content {
      # HEC endpoint — note: Splunk has no destination-level role_arn; the
      # S3 configuration's role carries the backup permissions.
      hec_endpoint               = dest.value.hec_endpoint
      hec_endpoint_type          = try(dest.value.hec_endpoint_type, null) != "" ? dest.value.hec_endpoint_type : null
      hec_token                  = try(dest.value.hec_token, null) != "" ? dest.value.hec_token : null
      hec_acknowledgment_timeout = try(dest.value.hec_acknowledgment_timeout_in_seconds, 0) > 0 ? dest.value.hec_acknowledgment_timeout_in_seconds : null

      # Delivery configuration (Splunk caps: 0-60s interval, 1-5 MiB size)
      buffering_interval = try(dest.value.buffering.interval_in_seconds, 0) > 0 ? dest.value.buffering.interval_in_seconds : null
      buffering_size     = try(dest.value.buffering.size_in_mbs, 0) > 0 ? dest.value.buffering.size_in_mbs : null
      retry_duration     = try(dest.value.retry_duration_in_seconds, 0) > 0 ? dest.value.retry_duration_in_seconds : null

      # S3 backup mode
      s3_backup_mode = try(dest.value.s3_backup_mode, null) != "" ? dest.value.s3_backup_mode : null

      # --- Secrets Manager credentials (XOR with hec_token, ForceNew switch) ---

      dynamic "secrets_manager_configuration" {
        for_each = try(dest.value.secrets_manager, null) != null ? [dest.value.secrets_manager] : []
        iterator = sm
        content {
          enabled    = true
          secret_arn = sm.value.secret_arn
          role_arn   = sm.value.role_arn != "" ? sm.value.role_arn : null
        }
      }

      # --- S3 config (required — backs up failed/all events) ---

      s3_configuration {
        bucket_arn          = dest.value.s3_config.bucket_arn
        role_arn            = dest.value.s3_config.role_arn
        prefix              = try(dest.value.s3_config.prefix, null) != "" ? dest.value.s3_config.prefix : null
        error_output_prefix = try(dest.value.s3_config.error_output_prefix, null) != "" ? dest.value.s3_config.error_output_prefix : null
        compression_format  = try(dest.value.s3_config.compression_format, null) != "" ? dest.value.s3_config.compression_format : null
        kms_key_arn         = dest.value.s3_config.kms_key_arn != "" ? dest.value.s3_config.kms_key_arn : null
        buffering_interval  = try(dest.value.s3_config.buffering.interval_in_seconds, 0) > 0 ? dest.value.s3_config.buffering.interval_in_seconds : null
        buffering_size      = try(dest.value.s3_config.buffering.size_in_mbs, 0) > 0 ? dest.value.s3_config.buffering.size_in_mbs : null

        # CloudWatch logging for THIS S3 leg (backup/staging errors) —
        # distinct from the destination-level logging block below.
        dynamic "cloudwatch_logging_options" {
          for_each = try(dest.value.s3_config.logging.enabled, false) ? [dest.value.s3_config.logging] : []
          iterator = s3log
          content {
            enabled         = true
            log_group_name  = s3log.value.log_group_name
            log_stream_name = s3log.value.log_stream_name
          }
        }
      }

      # --- Processing pipeline (normalized typed processors) ---

      dynamic "processing_configuration" {
        for_each = local.processing_enabled ? [1] : []
        content {
          enabled = true

          dynamic "processors" {
            for_each = local.processors_normalized
            content {
              type = processors.value.type

              dynamic "parameters" {
                for_each = processors.value.parameters
                content {
                  parameter_name  = parameters.value.name
                  parameter_value = parameters.value.value
                }
              }
            }
          }
        }
      }

      # --- CloudWatch logging ---

      dynamic "cloudwatch_logging_options" {
        for_each = try(dest.value.logging.enabled, false) ? [dest.value.logging] : []
        iterator = log
        content {
          enabled         = true
          log_group_name  = log.value.log_group_name
          log_stream_name = log.value.log_stream_name
        }
      }
    }
  }

  # ===========================================================================
  # Snowflake destination
  # ===========================================================================

  dynamic "snowflake_configuration" {
    for_each = local.destination_type == "snowflake" ? [var.spec.snowflake] : []
    iterator = dest
    content {
      # Snowflake target
      account_url = dest.value.account_url
      database    = dest.value.database
      schema      = dest.value.schema
      table       = dest.value.table
      role_arn    = dest.value.role_arn

      # Authentication — inline key pair XOR Secrets Manager (CEL-enforced)
      user           = try(dest.value.user, null) != "" ? dest.value.user : null
      private_key    = try(dest.value.private_key, null) != "" ? dest.value.private_key : null
      key_passphrase = try(dest.value.key_passphrase, null) != "" ? dest.value.key_passphrase : null

      # Data loading
      data_loading_option  = try(dest.value.data_loading_option, null) != "" ? dest.value.data_loading_option : null
      content_column_name  = try(dest.value.content_column_name, null) != "" ? dest.value.content_column_name : null
      metadata_column_name = try(dest.value.metadata_column_name, null) != "" ? dest.value.metadata_column_name : null

      # Delivery configuration (Snowpipe Streaming defaults: 0s / 1 MiB)
      buffering_interval = try(dest.value.buffering.interval_in_seconds, 0) > 0 ? dest.value.buffering.interval_in_seconds : null
      buffering_size     = try(dest.value.buffering.size_in_mbs, 0) > 0 ? dest.value.buffering.size_in_mbs : null
      retry_duration     = try(dest.value.retry_duration_in_seconds, 0) > 0 ? dest.value.retry_duration_in_seconds : null

      # S3 backup mode
      s3_backup_mode = try(dest.value.s3_backup_mode, null) != "" ? dest.value.s3_backup_mode : null

      # --- Secrets Manager credentials (ForceNew switch) ---

      dynamic "secrets_manager_configuration" {
        for_each = try(dest.value.secrets_manager, null) != null ? [dest.value.secrets_manager] : []
        iterator = sm
        content {
          enabled    = true
          secret_arn = sm.value.secret_arn
          role_arn   = sm.value.role_arn != "" ? sm.value.role_arn : null
        }
      }

      # --- Snowflake role (least-privilege ingestion role) ---

      dynamic "snowflake_role_configuration" {
        for_each = try(dest.value.snowflake_role, null) != "" ? [dest.value.snowflake_role] : []
        content {
          enabled        = true
          snowflake_role = snowflake_role_configuration.value
        }
      }

      # --- PrivateLink connectivity ---

      dynamic "snowflake_vpc_configuration" {
        for_each = try(dest.value.private_link_vpce_id, null) != "" ? [dest.value.private_link_vpce_id] : []
        content {
          private_link_vpce_id = snowflake_vpc_configuration.value
        }
      }

      # --- S3 config (required — backs up failed/all records) ---

      s3_configuration {
        bucket_arn          = dest.value.s3_config.bucket_arn
        role_arn            = dest.value.s3_config.role_arn
        prefix              = try(dest.value.s3_config.prefix, null) != "" ? dest.value.s3_config.prefix : null
        error_output_prefix = try(dest.value.s3_config.error_output_prefix, null) != "" ? dest.value.s3_config.error_output_prefix : null
        compression_format  = try(dest.value.s3_config.compression_format, null) != "" ? dest.value.s3_config.compression_format : null
        kms_key_arn         = dest.value.s3_config.kms_key_arn != "" ? dest.value.s3_config.kms_key_arn : null
        buffering_interval  = try(dest.value.s3_config.buffering.interval_in_seconds, 0) > 0 ? dest.value.s3_config.buffering.interval_in_seconds : null
        buffering_size      = try(dest.value.s3_config.buffering.size_in_mbs, 0) > 0 ? dest.value.s3_config.buffering.size_in_mbs : null

        # CloudWatch logging for THIS S3 leg (backup/staging errors) —
        # distinct from the destination-level logging block below.
        dynamic "cloudwatch_logging_options" {
          for_each = try(dest.value.s3_config.logging.enabled, false) ? [dest.value.s3_config.logging] : []
          iterator = s3log
          content {
            enabled         = true
            log_group_name  = s3log.value.log_group_name
            log_stream_name = s3log.value.log_stream_name
          }
        }
      }

      # --- Processing pipeline (normalized typed processors) ---

      dynamic "processing_configuration" {
        for_each = local.processing_enabled ? [1] : []
        content {
          enabled = true

          dynamic "processors" {
            for_each = local.processors_normalized
            content {
              type = processors.value.type

              dynamic "parameters" {
                for_each = processors.value.parameters
                content {
                  parameter_name  = parameters.value.name
                  parameter_value = parameters.value.value
                }
              }
            }
          }
        }
      }

      # --- CloudWatch logging ---

      dynamic "cloudwatch_logging_options" {
        for_each = try(dest.value.logging.enabled, false) ? [dest.value.logging] : []
        iterator = log
        content {
          enabled         = true
          log_group_name  = log.value.log_group_name
          log_stream_name = log.value.log_stream_name
        }
      }
    }
  }

  # ===========================================================================
  # Iceberg destination
  # ===========================================================================

  dynamic "iceberg_configuration" {
    for_each = local.destination_type == "iceberg" ? [var.spec.iceberg] : []
    iterator = dest
    content {
      # Glue catalog target (ForceNew) and delivery role
      catalog_arn = dest.value.catalog_arn
      role_arn    = dest.value.role_arn

      # Append-only mode (ForceNew) — omit when false so AWS owns the default.
      append_only = dest.value.append_only ? true : null

      # Delivery configuration
      buffering_interval = try(dest.value.buffering.interval_in_seconds, 0) > 0 ? dest.value.buffering.interval_in_seconds : null
      buffering_size     = try(dest.value.buffering.size_in_mbs, 0) > 0 ? dest.value.buffering.size_in_mbs : null
      retry_duration     = try(dest.value.retry_duration_in_seconds, 0) > 0 ? dest.value.retry_duration_in_seconds : null

      # S3 backup mode
      s3_backup_mode = try(dest.value.s3_backup_mode, null) != "" ? dest.value.s3_backup_mode : null

      # --- Destination tables (ForceNew; unique keys enable upserts) ---

      dynamic "destination_table_configuration" {
        for_each = try(dest.value.destination_tables, [])
        iterator = tbl
        content {
          database_name          = tbl.value.database_name
          table_name             = tbl.value.table_name
          s3_error_output_prefix = try(tbl.value.s3_error_output_prefix, null) != "" ? tbl.value.s3_error_output_prefix : null
          unique_keys            = length(try(tbl.value.unique_keys, [])) > 0 ? tbl.value.unique_keys : null
        }
      }

      # --- S3 config (required — backs up failed/all records) ---

      s3_configuration {
        bucket_arn          = dest.value.s3_config.bucket_arn
        role_arn            = dest.value.s3_config.role_arn
        prefix              = try(dest.value.s3_config.prefix, null) != "" ? dest.value.s3_config.prefix : null
        error_output_prefix = try(dest.value.s3_config.error_output_prefix, null) != "" ? dest.value.s3_config.error_output_prefix : null
        compression_format  = try(dest.value.s3_config.compression_format, null) != "" ? dest.value.s3_config.compression_format : null
        kms_key_arn         = dest.value.s3_config.kms_key_arn != "" ? dest.value.s3_config.kms_key_arn : null
        buffering_interval  = try(dest.value.s3_config.buffering.interval_in_seconds, 0) > 0 ? dest.value.s3_config.buffering.interval_in_seconds : null
        buffering_size      = try(dest.value.s3_config.buffering.size_in_mbs, 0) > 0 ? dest.value.s3_config.buffering.size_in_mbs : null

        # CloudWatch logging for THIS S3 leg (backup/staging errors) —
        # distinct from the destination-level logging block below.
        dynamic "cloudwatch_logging_options" {
          for_each = try(dest.value.s3_config.logging.enabled, false) ? [dest.value.s3_config.logging] : []
          iterator = s3log
          content {
            enabled         = true
            log_group_name  = s3log.value.log_group_name
            log_stream_name = s3log.value.log_stream_name
          }
        }
      }

      # --- Processing pipeline (normalized typed processors) ---

      dynamic "processing_configuration" {
        for_each = local.processing_enabled ? [1] : []
        content {
          enabled = true

          dynamic "processors" {
            for_each = local.processors_normalized
            content {
              type = processors.value.type

              dynamic "parameters" {
                for_each = processors.value.parameters
                content {
                  parameter_name  = parameters.value.name
                  parameter_value = parameters.value.value
                }
              }
            }
          }
        }
      }

      # --- CloudWatch logging ---

      dynamic "cloudwatch_logging_options" {
        for_each = try(dest.value.logging.enabled, false) ? [dest.value.logging] : []
        iterator = log
        content {
          enabled         = true
          log_group_name  = log.value.log_group_name
          log_stream_name = log.value.log_stream_name
        }
      }
    }
  }
}
