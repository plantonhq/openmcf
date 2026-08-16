# One EventBridge Pipe: source -> (filter) -> (enrichment) -> target.
#
# Lifecycle facts the render below depends on:
#   - the pipe's name (metadata.name) and its SOURCE are fixed for
#     life (replace-on-change) - so are the per-family stream/topic
#     positions (starting_position, topic_name, queue_name,
#     consumer_group_id on self-managed Kafka); the TARGET swaps in
#     place;
#   - desired_state RUNNING/STOPPED flips consumption without
#     deleting; creates and state flips can take minutes (the
#     provider waits up to 30);
#   - Kafka/MQ credentials are Secrets Manager secret ARNs -
#     references, never credential values;
#   - removing target_parameters.input_template genuinely clears it at
#     AWS (the provider pre-seeds the empty value on update);
#   - the spec's CELs guarantee at most one source family block and at
#     most one target family block.

resource "aws_pipes_pipe" "this" {
  name     = var.metadata.name
  source   = var.spec.source
  target   = var.spec.target
  role_arn = var.spec.role_arn

  description        = var.spec.description != "" ? var.spec.description : null
  desired_state      = var.spec.desired_state != "" ? var.spec.desired_state : null
  enrichment         = var.spec.enrichment != "" ? var.spec.enrichment : null
  kms_key_identifier = var.spec.kms_key_identifier != "" ? var.spec.kms_key_identifier : null

  dynamic "source_parameters" {
    for_each = var.spec.source_parameters != null ? [var.spec.source_parameters] : []
    content {
      dynamic "filter_criteria" {
        for_each = source_parameters.value.filter_criteria != null ? [source_parameters.value.filter_criteria] : []
        content {
          dynamic "filter" {
            for_each = filter_criteria.value.filters
            content {
              pattern = filter.value.pattern
            }
          }
        }
      }

      dynamic "sqs_queue_parameters" {
        for_each = source_parameters.value.sqs != null ? [source_parameters.value.sqs] : []
        content {
          batch_size                         = sqs_queue_parameters.value.batch_size != null ? sqs_queue_parameters.value.batch_size : null
          maximum_batching_window_in_seconds = sqs_queue_parameters.value.maximum_batching_window_in_seconds != null ? sqs_queue_parameters.value.maximum_batching_window_in_seconds : null
        }
      }

      dynamic "kinesis_stream_parameters" {
        for_each = source_parameters.value.kinesis != null ? [source_parameters.value.kinesis] : []
        content {
          starting_position                  = kinesis_stream_parameters.value.starting_position
          starting_position_timestamp        = kinesis_stream_parameters.value.starting_position_timestamp != "" ? kinesis_stream_parameters.value.starting_position_timestamp : null
          batch_size                         = kinesis_stream_parameters.value.batch_size != null ? kinesis_stream_parameters.value.batch_size : null
          maximum_batching_window_in_seconds = kinesis_stream_parameters.value.maximum_batching_window_in_seconds != null ? kinesis_stream_parameters.value.maximum_batching_window_in_seconds : null
          maximum_record_age_in_seconds      = kinesis_stream_parameters.value.maximum_record_age_in_seconds != null ? kinesis_stream_parameters.value.maximum_record_age_in_seconds : null
          maximum_retry_attempts             = kinesis_stream_parameters.value.maximum_retry_attempts != null ? kinesis_stream_parameters.value.maximum_retry_attempts : null
          on_partial_batch_item_failure      = kinesis_stream_parameters.value.on_partial_batch_item_failure != "" ? kinesis_stream_parameters.value.on_partial_batch_item_failure : null
          parallelization_factor             = kinesis_stream_parameters.value.parallelization_factor != null ? kinesis_stream_parameters.value.parallelization_factor : null

          dynamic "dead_letter_config" {
            for_each = kinesis_stream_parameters.value.dead_letter_queue_arn != "" ? [1] : []
            content {
              arn = kinesis_stream_parameters.value.dead_letter_queue_arn
            }
          }
        }
      }

      dynamic "dynamodb_stream_parameters" {
        for_each = source_parameters.value.dynamodb != null ? [source_parameters.value.dynamodb] : []
        content {
          starting_position                  = dynamodb_stream_parameters.value.starting_position
          batch_size                         = dynamodb_stream_parameters.value.batch_size != null ? dynamodb_stream_parameters.value.batch_size : null
          maximum_batching_window_in_seconds = dynamodb_stream_parameters.value.maximum_batching_window_in_seconds != null ? dynamodb_stream_parameters.value.maximum_batching_window_in_seconds : null
          maximum_record_age_in_seconds      = dynamodb_stream_parameters.value.maximum_record_age_in_seconds != null ? dynamodb_stream_parameters.value.maximum_record_age_in_seconds : null
          maximum_retry_attempts             = dynamodb_stream_parameters.value.maximum_retry_attempts != null ? dynamodb_stream_parameters.value.maximum_retry_attempts : null
          on_partial_batch_item_failure      = dynamodb_stream_parameters.value.on_partial_batch_item_failure != "" ? dynamodb_stream_parameters.value.on_partial_batch_item_failure : null
          parallelization_factor             = dynamodb_stream_parameters.value.parallelization_factor != null ? dynamodb_stream_parameters.value.parallelization_factor : null

          dynamic "dead_letter_config" {
            for_each = dynamodb_stream_parameters.value.dead_letter_queue_arn != "" ? [1] : []
            content {
              arn = dynamodb_stream_parameters.value.dead_letter_queue_arn
            }
          }
        }
      }

      dynamic "managed_streaming_kafka_parameters" {
        for_each = source_parameters.value.msk != null ? [source_parameters.value.msk] : []
        content {
          topic_name                         = managed_streaming_kafka_parameters.value.topic_name
          consumer_group_id                  = managed_streaming_kafka_parameters.value.consumer_group_id != "" ? managed_streaming_kafka_parameters.value.consumer_group_id : null
          starting_position                  = managed_streaming_kafka_parameters.value.starting_position != "" ? managed_streaming_kafka_parameters.value.starting_position : null
          batch_size                         = managed_streaming_kafka_parameters.value.batch_size != null ? managed_streaming_kafka_parameters.value.batch_size : null
          maximum_batching_window_in_seconds = managed_streaming_kafka_parameters.value.maximum_batching_window_in_seconds != null ? managed_streaming_kafka_parameters.value.maximum_batching_window_in_seconds : null

          dynamic "credentials" {
            for_each = managed_streaming_kafka_parameters.value.credentials != null ? [managed_streaming_kafka_parameters.value.credentials] : []
            content {
              client_certificate_tls_auth = credentials.value.client_certificate_tls_auth != "" ? credentials.value.client_certificate_tls_auth : null
              sasl_scram_512_auth         = credentials.value.sasl_scram_512_auth != "" ? credentials.value.sasl_scram_512_auth : null
            }
          }
        }
      }

      dynamic "self_managed_kafka_parameters" {
        for_each = source_parameters.value.self_managed_kafka != null ? [source_parameters.value.self_managed_kafka] : []
        content {
          topic_name                         = self_managed_kafka_parameters.value.topic_name
          additional_bootstrap_servers       = self_managed_kafka_parameters.value.additional_bootstrap_servers != null && length(self_managed_kafka_parameters.value.additional_bootstrap_servers) > 0 ? self_managed_kafka_parameters.value.additional_bootstrap_servers : null
          consumer_group_id                  = self_managed_kafka_parameters.value.consumer_group_id != "" ? self_managed_kafka_parameters.value.consumer_group_id : null
          starting_position                  = self_managed_kafka_parameters.value.starting_position != "" ? self_managed_kafka_parameters.value.starting_position : null
          batch_size                         = self_managed_kafka_parameters.value.batch_size != null ? self_managed_kafka_parameters.value.batch_size : null
          maximum_batching_window_in_seconds = self_managed_kafka_parameters.value.maximum_batching_window_in_seconds != null ? self_managed_kafka_parameters.value.maximum_batching_window_in_seconds : null
          server_root_ca_certificate         = self_managed_kafka_parameters.value.server_root_ca_certificate != "" ? self_managed_kafka_parameters.value.server_root_ca_certificate : null

          dynamic "credentials" {
            for_each = self_managed_kafka_parameters.value.credentials != null ? [self_managed_kafka_parameters.value.credentials] : []
            content {
              basic_auth                  = credentials.value.basic_auth != "" ? credentials.value.basic_auth : null
              client_certificate_tls_auth = credentials.value.client_certificate_tls_auth != "" ? credentials.value.client_certificate_tls_auth : null
              sasl_scram_256_auth         = credentials.value.sasl_scram_256_auth != "" ? credentials.value.sasl_scram_256_auth : null
              sasl_scram_512_auth         = credentials.value.sasl_scram_512_auth != "" ? credentials.value.sasl_scram_512_auth : null
            }
          }

          dynamic "vpc" {
            for_each = self_managed_kafka_parameters.value.vpc != null ? [self_managed_kafka_parameters.value.vpc] : []
            content {
              subnets         = vpc.value.subnets
              security_groups = vpc.value.security_groups != null && length(vpc.value.security_groups) > 0 ? vpc.value.security_groups : null
            }
          }
        }
      }

      dynamic "activemq_broker_parameters" {
        for_each = source_parameters.value.activemq != null ? [source_parameters.value.activemq] : []
        content {
          queue_name                         = activemq_broker_parameters.value.queue_name
          batch_size                         = activemq_broker_parameters.value.batch_size != null ? activemq_broker_parameters.value.batch_size : null
          maximum_batching_window_in_seconds = activemq_broker_parameters.value.maximum_batching_window_in_seconds != null ? activemq_broker_parameters.value.maximum_batching_window_in_seconds : null

          credentials {
            basic_auth = activemq_broker_parameters.value.basic_auth_credentials
          }
        }
      }

      dynamic "rabbitmq_broker_parameters" {
        for_each = source_parameters.value.rabbitmq != null ? [source_parameters.value.rabbitmq] : []
        content {
          queue_name                         = rabbitmq_broker_parameters.value.queue_name
          virtual_host                       = rabbitmq_broker_parameters.value.virtual_host != "" ? rabbitmq_broker_parameters.value.virtual_host : null
          batch_size                         = rabbitmq_broker_parameters.value.batch_size != null ? rabbitmq_broker_parameters.value.batch_size : null
          maximum_batching_window_in_seconds = rabbitmq_broker_parameters.value.maximum_batching_window_in_seconds != null ? rabbitmq_broker_parameters.value.maximum_batching_window_in_seconds : null

          credentials {
            basic_auth = rabbitmq_broker_parameters.value.basic_auth_credentials
          }
        }
      }
    }
  }

  dynamic "enrichment_parameters" {
    for_each = var.spec.enrichment_parameters != null ? [var.spec.enrichment_parameters] : []
    content {
      input_template = enrichment_parameters.value.input_template != "" ? enrichment_parameters.value.input_template : null

      dynamic "http_parameters" {
        for_each = enrichment_parameters.value.http_parameters != null ? [enrichment_parameters.value.http_parameters] : []
        content {
          header_parameters       = length(http_parameters.value.header_parameters) > 0 ? http_parameters.value.header_parameters : null
          path_parameter_values   = http_parameters.value.path_parameter_value != "" ? [http_parameters.value.path_parameter_value] : null
          query_string_parameters = length(http_parameters.value.query_string_parameters) > 0 ? http_parameters.value.query_string_parameters : null
        }
      }
    }
  }

  dynamic "target_parameters" {
    for_each = var.spec.target_parameters != null ? [var.spec.target_parameters] : []
    content {
      input_template = target_parameters.value.input_template != "" ? target_parameters.value.input_template : null

      dynamic "sqs_queue_parameters" {
        for_each = target_parameters.value.sqs != null ? [target_parameters.value.sqs] : []
        content {
          message_group_id         = sqs_queue_parameters.value.message_group_id != "" ? sqs_queue_parameters.value.message_group_id : null
          message_deduplication_id = sqs_queue_parameters.value.message_deduplication_id != "" ? sqs_queue_parameters.value.message_deduplication_id : null
        }
      }

      dynamic "kinesis_stream_parameters" {
        for_each = target_parameters.value.kinesis != null ? [target_parameters.value.kinesis] : []
        content {
          partition_key = kinesis_stream_parameters.value.partition_key
        }
      }

      dynamic "lambda_function_parameters" {
        for_each = target_parameters.value.lambda != null ? [target_parameters.value.lambda] : []
        content {
          invocation_type = lambda_function_parameters.value.invocation_type
        }
      }

      dynamic "step_function_state_machine_parameters" {
        for_each = target_parameters.value.step_function != null ? [target_parameters.value.step_function] : []
        content {
          invocation_type = step_function_state_machine_parameters.value.invocation_type
        }
      }

      dynamic "ecs_task_parameters" {
        for_each = target_parameters.value.ecs_task != null ? [target_parameters.value.ecs_task] : []
        content {
          task_definition_arn     = ecs_task_parameters.value.task_definition_arn
          task_count              = ecs_task_parameters.value.task_count != null ? ecs_task_parameters.value.task_count : null
          launch_type             = ecs_task_parameters.value.launch_type != "" ? ecs_task_parameters.value.launch_type : null
          group                   = ecs_task_parameters.value.group != "" ? ecs_task_parameters.value.group : null
          platform_version        = ecs_task_parameters.value.platform_version != "" ? ecs_task_parameters.value.platform_version : null
          propagate_tags          = ecs_task_parameters.value.propagate_tags != "" ? ecs_task_parameters.value.propagate_tags : null
          reference_id            = ecs_task_parameters.value.reference_id != "" ? ecs_task_parameters.value.reference_id : null
          enable_ecs_managed_tags = ecs_task_parameters.value.enable_ecs_managed_tags
          enable_execute_command  = ecs_task_parameters.value.enable_execute_command
          tags                    = length(ecs_task_parameters.value.tags) > 0 ? ecs_task_parameters.value.tags : null

          dynamic "capacity_provider_strategy" {
            for_each = ecs_task_parameters.value.capacity_provider_strategy != null ? ecs_task_parameters.value.capacity_provider_strategy : []
            content {
              capacity_provider = capacity_provider_strategy.value.capacity_provider
              base              = capacity_provider_strategy.value.base
              weight            = capacity_provider_strategy.value.weight
            }
          }

          dynamic "network_configuration" {
            for_each = ecs_task_parameters.value.network_configuration != null ? [ecs_task_parameters.value.network_configuration] : []
            content {
              aws_vpc_configuration {
                subnets          = network_configuration.value.subnets
                security_groups  = network_configuration.value.security_groups != null && length(network_configuration.value.security_groups) > 0 ? network_configuration.value.security_groups : null
                assign_public_ip = network_configuration.value.assign_public_ip ? "ENABLED" : "DISABLED"
              }
            }
          }

          dynamic "overrides" {
            for_each = ecs_task_parameters.value.overrides != null ? [ecs_task_parameters.value.overrides] : []
            content {
              cpu                = overrides.value.cpu != "" ? overrides.value.cpu : null
              memory             = overrides.value.memory != "" ? overrides.value.memory : null
              execution_role_arn = overrides.value.execution_role_arn != "" ? overrides.value.execution_role_arn : null
              task_role_arn      = overrides.value.task_role_arn != "" ? overrides.value.task_role_arn : null

              dynamic "ephemeral_storage" {
                for_each = overrides.value.ephemeral_storage_size_in_gib != null ? [overrides.value.ephemeral_storage_size_in_gib] : []
                content {
                  size_in_gib = ephemeral_storage.value
                }
              }

              dynamic "container_override" {
                for_each = overrides.value.container_overrides != null ? overrides.value.container_overrides : []
                content {
                  name               = container_override.value.name != "" ? container_override.value.name : null
                  command            = container_override.value.command != null && length(container_override.value.command) > 0 ? container_override.value.command : null
                  cpu                = container_override.value.cpu != null ? container_override.value.cpu : null
                  memory             = container_override.value.memory != null ? container_override.value.memory : null
                  memory_reservation = container_override.value.memory_reservation != null ? container_override.value.memory_reservation : null

                  dynamic "environment" {
                    for_each = container_override.value.environment != null ? container_override.value.environment : []
                    content {
                      name  = environment.value.name
                      value = environment.value.value
                    }
                  }

                  dynamic "environment_file" {
                    for_each = container_override.value.environment_files != null ? container_override.value.environment_files : []
                    content {
                      type  = environment_file.value.type
                      value = environment_file.value.value
                    }
                  }

                  dynamic "resource_requirement" {
                    for_each = container_override.value.resource_requirements != null ? container_override.value.resource_requirements : []
                    content {
                      type  = resource_requirement.value.type
                      value = resource_requirement.value.value
                    }
                  }
                }
              }

              dynamic "inference_accelerator_override" {
                for_each = overrides.value.inference_accelerator_overrides != null ? overrides.value.inference_accelerator_overrides : []
                content {
                  device_name = inference_accelerator_override.value.device_name != "" ? inference_accelerator_override.value.device_name : null
                  device_type = inference_accelerator_override.value.device_type != "" ? inference_accelerator_override.value.device_type : null
                }
              }
            }
          }

          dynamic "placement_constraint" {
            for_each = ecs_task_parameters.value.placement_constraints != null ? ecs_task_parameters.value.placement_constraints : []
            content {
              type       = placement_constraint.value.type
              expression = placement_constraint.value.expression != "" ? placement_constraint.value.expression : null
            }
          }

          dynamic "placement_strategy" {
            for_each = ecs_task_parameters.value.placement_strategy != null ? ecs_task_parameters.value.placement_strategy : []
            content {
              type  = placement_strategy.value.type
              field = placement_strategy.value.field != "" ? placement_strategy.value.field : null
            }
          }
        }
      }

      dynamic "batch_job_parameters" {
        for_each = target_parameters.value.batch_job != null ? [target_parameters.value.batch_job] : []
        content {
          job_definition = batch_job_parameters.value.job_definition
          job_name       = batch_job_parameters.value.job_name
          parameters     = length(batch_job_parameters.value.parameters) > 0 ? batch_job_parameters.value.parameters : null

          dynamic "array_properties" {
            for_each = batch_job_parameters.value.array_size != null ? [batch_job_parameters.value.array_size] : []
            content {
              size = array_properties.value
            }
          }

          dynamic "retry_strategy" {
            for_each = batch_job_parameters.value.retry_attempts != null ? [batch_job_parameters.value.retry_attempts] : []
            content {
              attempts = retry_strategy.value
            }
          }

          dynamic "depends_on" {
            for_each = batch_job_parameters.value.depends_on != null ? batch_job_parameters.value.depends_on : []
            content {
              job_id = depends_on.value.job_id != "" ? depends_on.value.job_id : null
              type   = depends_on.value.type != "" ? depends_on.value.type : null
            }
          }

          dynamic "container_overrides" {
            for_each = batch_job_parameters.value.container_overrides != null ? [batch_job_parameters.value.container_overrides] : []
            content {
              command       = container_overrides.value.command != null && length(container_overrides.value.command) > 0 ? container_overrides.value.command : null
              instance_type = container_overrides.value.instance_type != "" ? container_overrides.value.instance_type : null

              dynamic "environment" {
                for_each = container_overrides.value.environment != null ? container_overrides.value.environment : []
                content {
                  name  = environment.value.name
                  value = environment.value.value
                }
              }

              dynamic "resource_requirement" {
                for_each = container_overrides.value.resource_requirements != null ? container_overrides.value.resource_requirements : []
                content {
                  type  = resource_requirement.value.type
                  value = resource_requirement.value.value
                }
              }
            }
          }
        }
      }

      dynamic "redshift_data_parameters" {
        for_each = target_parameters.value.redshift_data != null ? [target_parameters.value.redshift_data] : []
        content {
          database           = redshift_data_parameters.value.database
          sqls               = redshift_data_parameters.value.sqls
          db_user            = redshift_data_parameters.value.db_user != "" ? redshift_data_parameters.value.db_user : null
          secret_manager_arn = redshift_data_parameters.value.secret_manager_arn != "" ? redshift_data_parameters.value.secret_manager_arn : null
          statement_name     = redshift_data_parameters.value.statement_name != "" ? redshift_data_parameters.value.statement_name : null
          with_event         = redshift_data_parameters.value.with_event
        }
      }

      dynamic "sagemaker_pipeline_parameters" {
        for_each = target_parameters.value.sagemaker_pipeline != null ? [target_parameters.value.sagemaker_pipeline] : []
        content {
          dynamic "pipeline_parameter" {
            for_each = sagemaker_pipeline_parameters.value.pipeline_parameters != null ? sagemaker_pipeline_parameters.value.pipeline_parameters : []
            content {
              name  = pipeline_parameter.value.name
              value = pipeline_parameter.value.value
            }
          }
        }
      }

      dynamic "eventbridge_event_bus_parameters" {
        for_each = target_parameters.value.eventbridge_event_bus != null ? [target_parameters.value.eventbridge_event_bus] : []
        content {
          detail_type = eventbridge_event_bus_parameters.value.detail_type != "" ? eventbridge_event_bus_parameters.value.detail_type : null
          source      = eventbridge_event_bus_parameters.value.source != "" ? eventbridge_event_bus_parameters.value.source : null
          endpoint_id = eventbridge_event_bus_parameters.value.endpoint_id != "" ? eventbridge_event_bus_parameters.value.endpoint_id : null
          resources   = eventbridge_event_bus_parameters.value.resources != null && length(eventbridge_event_bus_parameters.value.resources) > 0 ? eventbridge_event_bus_parameters.value.resources : null
          time        = eventbridge_event_bus_parameters.value.time != "" ? eventbridge_event_bus_parameters.value.time : null
        }
      }

      dynamic "cloudwatch_logs_parameters" {
        for_each = target_parameters.value.cloudwatch_logs != null ? [target_parameters.value.cloudwatch_logs] : []
        content {
          log_stream_name = cloudwatch_logs_parameters.value.log_stream_name != "" ? cloudwatch_logs_parameters.value.log_stream_name : null
          timestamp       = cloudwatch_logs_parameters.value.timestamp != "" ? cloudwatch_logs_parameters.value.timestamp : null
        }
      }

      dynamic "http_parameters" {
        for_each = target_parameters.value.http != null ? [target_parameters.value.http] : []
        content {
          header_parameters       = length(http_parameters.value.header_parameters) > 0 ? http_parameters.value.header_parameters : null
          path_parameter_values   = http_parameters.value.path_parameter_value != "" ? [http_parameters.value.path_parameter_value] : null
          query_string_parameters = length(http_parameters.value.query_string_parameters) > 0 ? http_parameters.value.query_string_parameters : null
        }
      }
    }
  }

  dynamic "log_configuration" {
    for_each = var.spec.log_configuration != null ? [var.spec.log_configuration] : []
    content {
      level                  = log_configuration.value.level
      include_execution_data = log_configuration.value.include_execution_data ? ["ALL"] : null

      dynamic "cloudwatch_logs_log_destination" {
        for_each = log_configuration.value.cloudwatch_logs != null ? [log_configuration.value.cloudwatch_logs] : []
        content {
          log_group_arn = cloudwatch_logs_log_destination.value.log_group_arn
        }
      }

      dynamic "firehose_log_destination" {
        for_each = log_configuration.value.firehose != null ? [log_configuration.value.firehose] : []
        content {
          delivery_stream_arn = firehose_log_destination.value.delivery_stream_arn
        }
      }

      dynamic "s3_log_destination" {
        for_each = log_configuration.value.s3 != null ? [log_configuration.value.s3] : []
        content {
          bucket_name   = s3_log_destination.value.bucket_name
          bucket_owner  = s3_log_destination.value.bucket_owner
          output_format = s3_log_destination.value.output_format != "" ? s3_log_destination.value.output_format : null
          prefix        = s3_log_destination.value.prefix != "" ? s3_log_destination.value.prefix : null
        }
      }
    }
  }

  tags = local.aws_tags
}
