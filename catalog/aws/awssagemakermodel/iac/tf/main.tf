# Amazon SageMaker AI model: the immutable serving definition endpoints
# deploy - a single primary container or an inference pipeline of
# containers.
#
# Lifecycle facts the renders below depend on:
#   - every argument is create-time only (the provider's update is
#     tags-only) - any spec change replaces the model, which is AWS's own
#     contract (roll a new model, repoint the endpoint);
#   - `primary_container` and `container` share one schema upstream; the
#     spec's exactly-one rule decides which form renders;
#   - the s3_data_source wrapper is single-valued upstream (the expander
#     reads index 0 only) - the spec flattens it to one message.

resource "aws_sagemaker_model" "this" {
  # The component's name IS the model name.
  name = local.model_name

  execution_role_arn = var.spec.execution_role_arn

  # Isolate the model container: no network calls in or out.
  enable_network_isolation = var.spec.enable_network_isolation

  # How pipeline containers are invoked - only meaningful with the
  # pipeline form (spec-validated).
  dynamic "inference_execution_config" {
    for_each = var.spec.inference_execution_mode != "" ? [var.spec.inference_execution_mode] : []
    content {
      mode = inference_execution_config.value
    }
  }

  # The single-container form.
  dynamic "primary_container" {
    for_each = var.spec.primary_container != null ? [var.spec.primary_container] : []
    content {
      image                        = primary_container.value.image != "" ? primary_container.value.image : null
      model_package_name           = primary_container.value.model_package_arn != "" ? primary_container.value.model_package_arn : null
      container_hostname           = primary_container.value.container_hostname != "" ? primary_container.value.container_hostname : null
      environment                  = length(primary_container.value.environment) > 0 ? primary_container.value.environment : null
      mode                         = primary_container.value.mode != "" ? primary_container.value.mode : null
      model_data_url               = primary_container.value.model_data_url != "" ? primary_container.value.model_data_url : null
      inference_specification_name = primary_container.value.inference_specification_name != "" ? primary_container.value.inference_specification_name : null

      dynamic "model_data_source" {
        for_each = primary_container.value.model_data_source != null ? [primary_container.value.model_data_source] : []
        content {
          s3_data_source {
            s3_uri           = model_data_source.value.s3_uri
            s3_data_type     = model_data_source.value.s3_data_type
            compression_type = model_data_source.value.compression_type
            dynamic "model_access_config" {
              for_each = model_data_source.value.accept_eula ? [true] : []
              content {
                accept_eula = true
              }
            }
          }
        }
      }

      dynamic "additional_model_data_source" {
        for_each = primary_container.value.additional_model_data_sources
        content {
          channel_name = additional_model_data_source.value.channel_name
          s3_data_source {
            s3_uri           = additional_model_data_source.value.source.s3_uri
            s3_data_type     = additional_model_data_source.value.source.s3_data_type
            compression_type = additional_model_data_source.value.source.compression_type
            dynamic "model_access_config" {
              for_each = additional_model_data_source.value.source.accept_eula ? [true] : []
              content {
                accept_eula = true
              }
            }
          }
        }
      }

      # MultiModel caching - only meaningful in MultiModel mode
      # (spec-validated).
      dynamic "multi_model_config" {
        for_each = primary_container.value.multi_model_cache != "" ? [primary_container.value.multi_model_cache] : []
        content {
          model_cache_setting = multi_model_config.value
        }
      }

      # Pull from a private VPC-reachable registry instead of ECR.
      dynamic "image_config" {
        for_each = primary_container.value.image_config != null ? [primary_container.value.image_config] : []
        content {
          repository_access_mode = image_config.value.repository_access_mode
          dynamic "repository_auth_config" {
            for_each = image_config.value.repository_credentials_provider_arn != "" ? [image_config.value.repository_credentials_provider_arn] : []
            content {
              repository_credentials_provider_arn = repository_auth_config.value
            }
          }
        }
      }
    }
  }

  # The inference-pipeline form (2-15 containers, same schema as the
  # primary container upstream).
  dynamic "container" {
    for_each = local.pipeline_containers
    content {
      image                        = container.value.image != "" ? container.value.image : null
      model_package_name           = container.value.model_package_arn != "" ? container.value.model_package_arn : null
      container_hostname           = container.value.container_hostname != "" ? container.value.container_hostname : null
      environment                  = length(container.value.environment) > 0 ? container.value.environment : null
      mode                         = container.value.mode != "" ? container.value.mode : null
      model_data_url               = container.value.model_data_url != "" ? container.value.model_data_url : null
      inference_specification_name = container.value.inference_specification_name != "" ? container.value.inference_specification_name : null

      dynamic "model_data_source" {
        for_each = container.value.model_data_source != null ? [container.value.model_data_source] : []
        content {
          s3_data_source {
            s3_uri           = model_data_source.value.s3_uri
            s3_data_type     = model_data_source.value.s3_data_type
            compression_type = model_data_source.value.compression_type
            dynamic "model_access_config" {
              for_each = model_data_source.value.accept_eula ? [true] : []
              content {
                accept_eula = true
              }
            }
          }
        }
      }

      dynamic "additional_model_data_source" {
        for_each = container.value.additional_model_data_sources
        content {
          channel_name = additional_model_data_source.value.channel_name
          s3_data_source {
            s3_uri           = additional_model_data_source.value.source.s3_uri
            s3_data_type     = additional_model_data_source.value.source.s3_data_type
            compression_type = additional_model_data_source.value.source.compression_type
            dynamic "model_access_config" {
              for_each = additional_model_data_source.value.source.accept_eula ? [true] : []
              content {
                accept_eula = true
              }
            }
          }
        }
      }

      dynamic "multi_model_config" {
        for_each = container.value.multi_model_cache != "" ? [container.value.multi_model_cache] : []
        content {
          model_cache_setting = multi_model_config.value
        }
      }

      dynamic "image_config" {
        for_each = container.value.image_config != null ? [container.value.image_config] : []
        content {
          repository_access_mode = image_config.value.repository_access_mode
          dynamic "repository_auth_config" {
            for_each = image_config.value.repository_credentials_provider_arn != "" ? [image_config.value.repository_credentials_provider_arn] : []
            content {
              repository_credentials_provider_arn = repository_auth_config.value
            }
          }
        }
      }
    }
  }

  # Attach the containers to your VPC (private serving).
  dynamic "vpc_config" {
    for_each = var.spec.vpc_config != null ? [var.spec.vpc_config] : []
    content {
      subnets            = vpc_config.value.subnet_ids
      security_group_ids = vpc_config.value.security_group_ids
    }
  }

  tags = local.aws_tags
}
