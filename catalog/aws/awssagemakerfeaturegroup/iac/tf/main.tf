# Amazon SageMaker Feature Store feature group: online and/or offline
# stores over a declared feature schema.
#
# Lifecycle facts the renders below depend on:
#   - almost everything is create-time only (schema, stores, role) - the
#     ONLY in-place updates are the online store's TTL and the
#     throughput settings;
#   - the throughput expander mirrors the provider's create behavior:
#     provisioned capacity units are sent only in Provisioned mode
#     (spec-validated pairing);
#   - a Vector collection pairs exactly with its dimension
#     (spec-validated; AWS requires InMemory online storage for
#     collection types - server-enforced).

resource "aws_sagemaker_feature_group" "this" {
  # The component's name IS the feature group name.
  feature_group_name = local.feature_group_name

  record_identifier_feature_name = var.spec.record_identifier_feature_name
  event_time_feature_name        = var.spec.event_time_feature_name

  description = var.spec.description != "" ? var.spec.description : null

  role_arn = var.spec.role_arn

  dynamic "feature_definition" {
    for_each = var.spec.feature_definitions
    content {
      feature_name = feature_definition.value.name
      feature_type = feature_definition.value.type

      collection_type = feature_definition.value.collection_type != "" ? feature_definition.value.collection_type : null

      # The provider's collection_config has exactly one member
      # (vector_config) - rendered exactly when the dimension is set
      # (spec pairs it with the Vector type).
      dynamic "collection_config" {
        for_each = feature_definition.value.vector_dimension != null ? [feature_definition.value.vector_dimension] : []
        content {
          vector_config {
            dimension = collection_config.value
          }
        }
      }
    }
  }

  dynamic "online_store_config" {
    for_each = local.has_online ? [var.spec.online_store] : []
    content {
      enable_online_store = online_store_config.value.enabled
      storage_type        = online_store_config.value.storage_type != "" ? online_store_config.value.storage_type : null

      dynamic "security_config" {
        for_each = online_store_config.value.kms_key_arn != "" ? [online_store_config.value.kms_key_arn] : []
        content {
          kms_key_id = security_config.value
        }
      }

      # The one online-store surface that updates in place.
      dynamic "ttl_duration" {
        for_each = online_store_config.value.ttl != null ? [online_store_config.value.ttl] : []
        content {
          unit  = ttl_duration.value.unit
          value = ttl_duration.value.value
        }
      }
    }
  }

  dynamic "offline_store_config" {
    for_each = local.has_offline ? [var.spec.offline_store] : []
    content {
      s3_storage_config {
        s3_uri     = offline_store_config.value.s3_uri
        kms_key_id = offline_store_config.value.kms_key_arn != "" ? offline_store_config.value.kms_key_arn : null
      }

      disable_glue_table_creation = offline_store_config.value.disable_glue_table_creation
      table_format                = offline_store_config.value.table_format != "" ? offline_store_config.value.table_format : null

      dynamic "data_catalog_config" {
        for_each = offline_store_config.value.data_catalog != null ? [offline_store_config.value.data_catalog] : []
        content {
          catalog    = data_catalog_config.value.catalog
          database   = data_catalog_config.value.database
          table_name = data_catalog_config.value.table_name
        }
      }
    }
  }

  dynamic "throughput_config" {
    for_each = var.spec.throughput != null ? [var.spec.throughput] : []
    content {
      throughput_mode                  = throughput_config.value.mode
      provisioned_read_capacity_units  = throughput_config.value.provisioned_read_capacity_units
      provisioned_write_capacity_units = throughput_config.value.provisioned_write_capacity_units
    }
  }

  tags = local.aws_tags
}
