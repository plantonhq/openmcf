# The event source mapping: the managed poller that reads an SQS queue, a
# Kinesis or DynamoDB stream, a Kafka topic (MSK or self-managed), an
# Amazon MQ queue, or a DocumentDB change stream and invokes a Lambda
# function with batched records. The event source (ARN or self-managed
# Kafka bootstrap servers) is create-time immutable; batching, filters,
# failure handling, and the target function edit in place.
resource "aws_lambda_event_source_mapping" "this" {
  function_name = var.spec.function_arn
  enabled       = !var.spec.disabled

  event_source_arn = !local.is_self_managed_kafka && var.spec.event_source_arn != "" ? var.spec.event_source_arn : null

  dynamic "self_managed_event_source" {
    for_each = local.is_self_managed_kafka ? [var.spec.self_managed_kafka] : []
    content {
      endpoints = {
        KAFKA_BOOTSTRAP_SERVERS = join(",", self_managed_event_source.value.bootstrap_servers)
      }
    }
  }

  batch_size                         = var.spec.batch_size != 0 ? var.spec.batch_size : null
  maximum_batching_window_in_seconds = var.spec.maximum_batching_window_seconds != 0 ? var.spec.maximum_batching_window_seconds : null

  dynamic "filter_criteria" {
    for_each = length(var.spec.filters) > 0 ? [1] : []
    content {
      dynamic "filter" {
        for_each = var.spec.filters
        content {
          pattern = filter.value.pattern != "" ? filter.value.pattern : null
        }
      }
    }
  }

  kms_key_arn = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  function_response_types = length(var.spec.function_response_types) > 0 ? var.spec.function_response_types : null

  dynamic "scaling_config" {
    for_each = var.spec.scaling_max_concurrency != 0 ? [var.spec.scaling_max_concurrency] : []
    content {
      maximum_concurrency = scaling_config.value
    }
  }

  dynamic "metrics_config" {
    for_each = length(var.spec.metrics) > 0 ? [var.spec.metrics] : []
    content {
      metrics = metrics_config.value
    }
  }

  starting_position           = var.spec.starting_position != "" ? var.spec.starting_position : null
  starting_position_timestamp = var.spec.starting_position_timestamp != "" ? var.spec.starting_position_timestamp : null

  parallelization_factor = var.spec.parallelization_factor != 0 ? var.spec.parallelization_factor : null

  maximum_record_age_in_seconds = var.spec.maximum_record_age_seconds != 0 ? var.spec.maximum_record_age_seconds : null

  maximum_retry_attempts = var.spec.maximum_retry_attempts

  bisect_batch_on_function_error = var.spec.bisect_batch_on_function_error ? true : null

  tumbling_window_in_seconds = var.spec.tumbling_window_seconds != 0 ? var.spec.tumbling_window_seconds : null

  dynamic "destination_config" {
    for_each = var.spec.on_failure_destination_arn != "" ? [var.spec.on_failure_destination_arn] : []
    content {
      on_failure {
        destination_arn = destination_config.value
      }
    }
  }

  topics = length(var.spec.topics) > 0 ? var.spec.topics : null

  dynamic "amazon_managed_kafka_event_source_config" {
    for_each = local.is_msk_source && local.has_kafka_source_config ? [1] : []
    content {
      consumer_group_id = var.spec.kafka_consumer_group_id != "" ? var.spec.kafka_consumer_group_id : null

      dynamic "schema_registry_config" {
        for_each = var.spec.schema_registry != null ? [var.spec.schema_registry] : []
        content {
          schema_registry_uri = schema_registry_config.value.uri
          event_record_format = schema_registry_config.value.event_record_format

          dynamic "schema_validation_config" {
            for_each = schema_registry_config.value.validation_attributes
            content {
              attribute = schema_validation_config.value
            }
          }

          dynamic "access_config" {
            for_each = schema_registry_config.value.access_configurations
            content {
              type = access_config.value.type
              uri  = access_config.value.uri
            }
          }
        }
      }
    }
  }

  dynamic "self_managed_kafka_event_source_config" {
    for_each = local.is_self_managed_kafka && local.has_kafka_source_config ? [1] : []
    content {
      consumer_group_id = var.spec.kafka_consumer_group_id != "" ? var.spec.kafka_consumer_group_id : null

      dynamic "schema_registry_config" {
        for_each = var.spec.schema_registry != null ? [var.spec.schema_registry] : []
        content {
          schema_registry_uri = schema_registry_config.value.uri
          event_record_format = schema_registry_config.value.event_record_format

          dynamic "schema_validation_config" {
            for_each = schema_registry_config.value.validation_attributes
            content {
              attribute = schema_validation_config.value
            }
          }

          dynamic "access_config" {
            for_each = schema_registry_config.value.access_configurations
            content {
              type = access_config.value.type
              uri  = access_config.value.uri
            }
          }
        }
      }
    }
  }

  dynamic "source_access_configuration" {
    for_each = var.spec.source_access_configurations
    content {
      type = source_access_configuration.value.type
      uri  = source_access_configuration.value.uri
    }
  }

  dynamic "provisioned_poller_config" {
    for_each = var.spec.provisioned_pollers != null ? [var.spec.provisioned_pollers] : []
    content {
      minimum_pollers = provisioned_poller_config.value.minimum_pollers != 0 ? provisioned_poller_config.value.minimum_pollers : null
      maximum_pollers = provisioned_poller_config.value.maximum_pollers != 0 ? provisioned_poller_config.value.maximum_pollers : null
    }
  }

  queues = var.spec.mq_queue != "" ? [var.spec.mq_queue] : null

  dynamic "document_db_event_source_config" {
    for_each = var.spec.document_db != null ? [var.spec.document_db] : []
    content {
      database_name   = document_db_event_source_config.value.database_name
      collection_name = document_db_event_source_config.value.collection_name != "" ? document_db_event_source_config.value.collection_name : null
      full_document   = document_db_event_source_config.value.full_document != "" ? document_db_event_source_config.value.full_document : null
    }
  }

  tags = local.aws_tags
}
