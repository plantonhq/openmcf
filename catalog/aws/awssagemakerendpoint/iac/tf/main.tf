# Amazon SageMaker AI endpoint + its endpoint configuration.
#
# Lifecycle facts the renders below depend on:
#   - the endpoint CONFIGURATION is immutable upstream (every argument
#     ForceNew; the provider's update is tags-only), while the ENDPOINT's
#     pointer to it updates in place (UpdateEndpoint, optionally shaped
#     by deployment_config). The declarative fold therefore rolls
#     configurations: name_prefix + create_before_destroy mint a NEW
#     suffixed configuration on any capacity change, UpdateEndpoint
#     repoints, and only then is the old configuration destroyed -- the
#     endpoint never references a deleted configuration (AWS's own
#     documented pattern);
#   - variant names default deterministically per position (locals) so a
#     re-plan never regenerates them (the provider would mint a random
#     name per plan otherwise, forcing a config roll every apply);
#   - the capacity-reservation preference has ONE legal value
#     ("capacity-reservations-only") -- the module owns the constant and
#     sends it exactly when an ML reservation ARN is configured.

resource "aws_sagemaker_endpoint_configuration" "this" {
  # Suffixed name so a changed configuration can coexist with its
  # predecessor during the endpoint repoint (create_before_destroy).
  name_prefix = "${local.endpoint_name}-cfg-"

  # Required only for inference-component endpoints (variants without a
  # model - spec-validated).
  execution_role_arn = var.spec.execution_role_arn != "" ? var.spec.execution_role_arn : null

  kms_key_arn = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  dynamic "production_variants" {
    for_each = local.production_variants
    content {
      variant_name           = production_variants.value.resolved_name
      model_name             = production_variants.value.model != "" ? production_variants.value.model : null
      instance_type          = production_variants.value.instance_type != "" ? production_variants.value.instance_type : null
      initial_instance_count = production_variants.value.initial_instance_count
      initial_variant_weight = production_variants.value.initial_variant_weight
      accelerator_type       = production_variants.value.accelerator_type != "" ? production_variants.value.accelerator_type : null
      inference_ami_version  = production_variants.value.inference_ami_version != "" ? production_variants.value.inference_ami_version : null
      enable_ssm_access      = production_variants.value.enable_ssm_access ? true : null
      volume_size_in_gb      = production_variants.value.volume_size_gb

      container_startup_health_check_timeout_in_seconds = production_variants.value.container_startup_health_check_timeout_seconds
      model_data_download_timeout_in_seconds            = production_variants.value.model_data_download_timeout_seconds

      dynamic "serverless_config" {
        for_each = production_variants.value.serverless != null ? [production_variants.value.serverless] : []
        content {
          max_concurrency         = serverless_config.value.max_concurrency
          memory_size_in_mb       = serverless_config.value.memory_size_mb
          provisioned_concurrency = serverless_config.value.provisioned_concurrency
        }
      }

      dynamic "managed_instance_scaling" {
        for_each = production_variants.value.managed_instance_scaling != null ? [production_variants.value.managed_instance_scaling] : []
        content {
          status             = managed_instance_scaling.value.status != "" ? managed_instance_scaling.value.status : null
          min_instance_count = managed_instance_scaling.value.min_instance_count
          max_instance_count = managed_instance_scaling.value.max_instance_count
        }
      }

      dynamic "routing_config" {
        for_each = production_variants.value.routing_strategy != "" ? [production_variants.value.routing_strategy] : []
        content {
          routing_strategy = routing_config.value
        }
      }

      dynamic "core_dump_config" {
        for_each = production_variants.value.core_dump != null ? [production_variants.value.core_dump] : []
        content {
          destination_s3_uri = core_dump_config.value.destination_s3_uri
          kms_key_id         = core_dump_config.value.kms_key_arn != "" ? core_dump_config.value.kms_key_arn : null
        }
      }

      dynamic "capacity_reservation_config" {
        for_each = production_variants.value.ml_capacity_reservation_arn != "" ? [production_variants.value.ml_capacity_reservation_arn] : []
        content {
          # The single legal preference value - the module's constant.
          capacity_reservation_preference = "capacity-reservations-only"
          ml_reservation_arn              = capacity_reservation_config.value
        }
      }
    }
  }

  dynamic "shadow_production_variants" {
    for_each = local.shadow_variants
    content {
      variant_name           = shadow_production_variants.value.resolved_name
      model_name             = shadow_production_variants.value.model != "" ? shadow_production_variants.value.model : null
      instance_type          = shadow_production_variants.value.instance_type != "" ? shadow_production_variants.value.instance_type : null
      initial_instance_count = shadow_production_variants.value.initial_instance_count
      initial_variant_weight = shadow_production_variants.value.initial_variant_weight
      accelerator_type       = shadow_production_variants.value.accelerator_type != "" ? shadow_production_variants.value.accelerator_type : null
      inference_ami_version  = shadow_production_variants.value.inference_ami_version != "" ? shadow_production_variants.value.inference_ami_version : null
      enable_ssm_access      = shadow_production_variants.value.enable_ssm_access ? true : null
      volume_size_in_gb      = shadow_production_variants.value.volume_size_gb

      container_startup_health_check_timeout_in_seconds = shadow_production_variants.value.container_startup_health_check_timeout_seconds
      model_data_download_timeout_in_seconds            = shadow_production_variants.value.model_data_download_timeout_seconds

      dynamic "serverless_config" {
        for_each = shadow_production_variants.value.serverless != null ? [shadow_production_variants.value.serverless] : []
        content {
          max_concurrency         = serverless_config.value.max_concurrency
          memory_size_in_mb       = serverless_config.value.memory_size_mb
          provisioned_concurrency = serverless_config.value.provisioned_concurrency
        }
      }

      dynamic "managed_instance_scaling" {
        for_each = shadow_production_variants.value.managed_instance_scaling != null ? [shadow_production_variants.value.managed_instance_scaling] : []
        content {
          status             = managed_instance_scaling.value.status != "" ? managed_instance_scaling.value.status : null
          min_instance_count = managed_instance_scaling.value.min_instance_count
          max_instance_count = managed_instance_scaling.value.max_instance_count
        }
      }

      dynamic "routing_config" {
        for_each = shadow_production_variants.value.routing_strategy != "" ? [shadow_production_variants.value.routing_strategy] : []
        content {
          routing_strategy = routing_config.value
        }
      }

      dynamic "core_dump_config" {
        for_each = shadow_production_variants.value.core_dump != null ? [shadow_production_variants.value.core_dump] : []
        content {
          destination_s3_uri = core_dump_config.value.destination_s3_uri
          kms_key_id         = core_dump_config.value.kms_key_arn != "" ? core_dump_config.value.kms_key_arn : null
        }
      }

      dynamic "capacity_reservation_config" {
        for_each = shadow_production_variants.value.ml_capacity_reservation_arn != "" ? [shadow_production_variants.value.ml_capacity_reservation_arn] : []
        content {
          capacity_reservation_preference = "capacity-reservations-only"
          ml_reservation_arn              = capacity_reservation_config.value
        }
      }
    }
  }

  # Queue requests and deliver responses to S3.
  dynamic "async_inference_config" {
    for_each = var.spec.async_inference != null ? [var.spec.async_inference] : []
    content {
      output_config {
        s3_output_path  = async_inference_config.value.output_s3_path
        s3_failure_path = async_inference_config.value.failure_s3_path != "" ? async_inference_config.value.failure_s3_path : null
        kms_key_id      = async_inference_config.value.kms_key_arn != "" ? async_inference_config.value.kms_key_arn : null

        dynamic "notification_config" {
          for_each = (async_inference_config.value.success_topic_arn != "" || async_inference_config.value.error_topic_arn != "" || length(async_inference_config.value.include_inference_response_in) > 0) ? [async_inference_config.value] : []
          content {
            success_topic                 = notification_config.value.success_topic_arn != "" ? notification_config.value.success_topic_arn : null
            error_topic                   = notification_config.value.error_topic_arn != "" ? notification_config.value.error_topic_arn : null
            include_inference_response_in = length(notification_config.value.include_inference_response_in) > 0 ? notification_config.value.include_inference_response_in : null
          }
        }
      }

      dynamic "client_config" {
        for_each = async_inference_config.value.max_concurrent_invocations_per_instance != null ? [async_inference_config.value.max_concurrent_invocations_per_instance] : []
        content {
          max_concurrent_invocations_per_instance = client_config.value
        }
      }
    }
  }

  # Capture request/response payloads to S3 (the Model Monitor feed).
  dynamic "data_capture_config" {
    for_each = var.spec.data_capture != null ? [var.spec.data_capture] : []
    content {
      destination_s3_uri          = data_capture_config.value.destination_s3_uri
      initial_sampling_percentage = data_capture_config.value.initial_sampling_percentage
      enable_capture              = data_capture_config.value.enable_capture
      kms_key_id                  = data_capture_config.value.kms_key_arn != "" ? data_capture_config.value.kms_key_arn : null

      dynamic "capture_options" {
        for_each = data_capture_config.value.capture_modes
        content {
          capture_mode = capture_options.value
        }
      }

      dynamic "capture_content_type_header" {
        for_each = (length(data_capture_config.value.csv_content_types) > 0 || length(data_capture_config.value.json_content_types) > 0) ? [data_capture_config.value] : []
        content {
          csv_content_types  = length(capture_content_type_header.value.csv_content_types) > 0 ? capture_content_type_header.value.csv_content_types : null
          json_content_types = length(capture_content_type_header.value.json_content_types) > 0 ? capture_content_type_header.value.json_content_types : null
        }
      }
    }
  }

  tags = local.aws_tags

  # A changed configuration must exist BEFORE the endpoint repoints at
  # it; the old one is destroyed only after the update.
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_sagemaker_endpoint" "this" {
  # The component's name IS the endpoint name.
  name = local.endpoint_name

  # Updates in place - the config roll's repoint edge.
  endpoint_config_name = aws_sagemaker_endpoint_configuration.this.name

  # How UpdateEndpoint rolls new capacity (exactly-one strategy,
  # spec-validated).
  dynamic "deployment_config" {
    for_each = var.spec.deployment != null ? [var.spec.deployment] : []
    content {
      dynamic "blue_green_update_policy" {
        for_each = deployment_config.value.blue_green != null ? [deployment_config.value.blue_green] : []
        content {
          traffic_routing_configuration {
            type                     = blue_green_update_policy.value.traffic_routing_type
            wait_interval_in_seconds = blue_green_update_policy.value.wait_interval_seconds

            dynamic "canary_size" {
              for_each = blue_green_update_policy.value.canary_size != null ? [blue_green_update_policy.value.canary_size] : []
              content {
                type  = canary_size.value.type
                value = canary_size.value.value
              }
            }

            dynamic "linear_step_size" {
              for_each = blue_green_update_policy.value.linear_step_size != null ? [blue_green_update_policy.value.linear_step_size] : []
              content {
                type  = linear_step_size.value.type
                value = linear_step_size.value.value
              }
            }
          }

          termination_wait_in_seconds          = blue_green_update_policy.value.termination_wait_seconds
          maximum_execution_timeout_in_seconds = blue_green_update_policy.value.maximum_execution_timeout_seconds
        }
      }

      dynamic "rolling_update_policy" {
        for_each = deployment_config.value.rolling != null ? [deployment_config.value.rolling] : []
        content {
          maximum_batch_size {
            type  = rolling_update_policy.value.maximum_batch_size.type
            value = rolling_update_policy.value.maximum_batch_size.value
          }

          wait_interval_in_seconds             = rolling_update_policy.value.wait_interval_seconds
          maximum_execution_timeout_in_seconds = rolling_update_policy.value.maximum_execution_timeout_seconds

          dynamic "rollback_maximum_batch_size" {
            for_each = rolling_update_policy.value.rollback_maximum_batch_size != null ? [rolling_update_policy.value.rollback_maximum_batch_size] : []
            content {
              type  = rollback_maximum_batch_size.value.type
              value = rollback_maximum_batch_size.value.value
            }
          }
        }
      }

      dynamic "auto_rollback_configuration" {
        for_each = length(deployment_config.value.auto_rollback_alarm_names) > 0 ? [deployment_config.value.auto_rollback_alarm_names] : []
        content {
          dynamic "alarms" {
            for_each = auto_rollback_configuration.value
            content {
              alarm_name = alarms.value
            }
          }
        }
      }
    }
  }

  tags = local.aws_tags
}
