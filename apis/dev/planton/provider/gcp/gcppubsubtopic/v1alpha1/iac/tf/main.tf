# Enable the Pub/Sub API — the control plane that owns the topic.
# disable_on_destroy is false: tearing down one topic must never disable
# the API for everything else in the project.
resource "google_project_service" "pubsub_api" {
  project = local.project_id
  service = "pubsub.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Pub/Sub topic — the named channel publishers send to. The topic name
# is ForceNew: renaming replaces the topic (and orphans its subscriptions),
# so treat names as permanent.
resource "google_pubsub_topic" "this" {
  name    = local.topic_name
  project = local.project_id
  labels  = local.final_labels

  # CMEK: the Pub/Sub service agent must hold
  # cloudkms.cryptoKeyEncrypterDecrypter on the key or publishes fail.
  kms_key_name = local.kms_key_name

  # Topic-level retention is independent of subscription retention, and
  # the lever that lets any subscription seek back within the window.
  message_retention_duration = local.message_retention_duration

  # Region pinning is enforced at publish time; enforce_in_transit
  # additionally rejects publishes from non-allowed regions instead of
  # rerouting them.
  dynamic "message_storage_policy" {
    for_each = var.spec.message_storage_policy != null ? [var.spec.message_storage_policy] : []
    content {
      allowed_persistence_regions = message_storage_policy.value.allowed_persistence_regions
      enforce_in_transit          = message_storage_policy.value.enforce_in_transit
    }
  }

  # Schema validation: the schema arrives resolved to the fully qualified
  # projects/{project}/schemas/{name} path.
  dynamic "schema_settings" {
    for_each = var.spec.schema_settings != null ? [var.spec.schema_settings] : []
    content {
      schema   = schema_settings.value.schema
      encoding = schema_settings.value.encoding != "" ? schema_settings.value.encoding : null
    }
  }

  dynamic "ingestion_data_source_settings" {
    for_each = var.spec.ingestion_data_source_settings != null ? [var.spec.ingestion_data_source_settings] : []
    content {

      dynamic "aws_kinesis" {
        for_each = ingestion_data_source_settings.value.aws_kinesis != null ? [ingestion_data_source_settings.value.aws_kinesis] : []
        content {
          stream_arn          = aws_kinesis.value.stream_arn
          consumer_arn        = aws_kinesis.value.consumer_arn
          aws_role_arn        = aws_kinesis.value.aws_role_arn
          gcp_service_account = aws_kinesis.value.gcp_service_account
        }
      }

      dynamic "aws_msk" {
        for_each = ingestion_data_source_settings.value.aws_msk != null ? [ingestion_data_source_settings.value.aws_msk] : []
        content {
          cluster_arn         = aws_msk.value.cluster_arn
          topic               = aws_msk.value.topic
          aws_role_arn        = aws_msk.value.aws_role_arn
          gcp_service_account = aws_msk.value.gcp_service_account
        }
      }

      dynamic "azure_event_hubs" {
        for_each = ingestion_data_source_settings.value.azure_event_hubs != null ? [ingestion_data_source_settings.value.azure_event_hubs] : []
        content {
          resource_group      = azure_event_hubs.value.resource_group != "" ? azure_event_hubs.value.resource_group : null
          namespace           = azure_event_hubs.value.namespace != "" ? azure_event_hubs.value.namespace : null
          event_hub           = azure_event_hubs.value.event_hub != "" ? azure_event_hubs.value.event_hub : null
          client_id           = azure_event_hubs.value.client_id != "" ? azure_event_hubs.value.client_id : null
          tenant_id           = azure_event_hubs.value.tenant_id != "" ? azure_event_hubs.value.tenant_id : null
          subscription_id     = azure_event_hubs.value.subscription_id != "" ? azure_event_hubs.value.subscription_id : null
          gcp_service_account = azure_event_hubs.value.gcp_service_account != "" ? azure_event_hubs.value.gcp_service_account : null
        }
      }

      dynamic "cloud_storage" {
        for_each = ingestion_data_source_settings.value.cloud_storage != null ? [ingestion_data_source_settings.value.cloud_storage] : []
        content {
          bucket                     = cloud_storage.value.bucket
          match_glob                 = cloud_storage.value.match_glob != "" ? cloud_storage.value.match_glob : null
          minimum_object_create_time = cloud_storage.value.minimum_object_create_time != "" ? cloud_storage.value.minimum_object_create_time : null

          # Format selection: exactly one should be set (API-enforced).
          dynamic "avro_format" {
            for_each = cloud_storage.value.avro_format ? [1] : []
            content {}
          }

          dynamic "pubsub_avro_format" {
            for_each = cloud_storage.value.pubsub_avro_format ? [1] : []
            content {}
          }

          dynamic "text_format" {
            for_each = cloud_storage.value.text_format != null ? [cloud_storage.value.text_format] : []
            content {
              delimiter = text_format.value.delimiter != "" ? text_format.value.delimiter : null
            }
          }
        }
      }

      dynamic "confluent_cloud" {
        for_each = ingestion_data_source_settings.value.confluent_cloud != null ? [ingestion_data_source_settings.value.confluent_cloud] : []
        content {
          bootstrap_server    = confluent_cloud.value.bootstrap_server
          topic               = confluent_cloud.value.topic
          identity_pool_id    = confluent_cloud.value.identity_pool_id
          gcp_service_account = confluent_cloud.value.gcp_service_account
          cluster_id          = confluent_cloud.value.cluster_id != "" ? confluent_cloud.value.cluster_id : null
        }
      }

      dynamic "platform_logs_settings" {
        for_each = ingestion_data_source_settings.value.platform_logs_settings != null ? [ingestion_data_source_settings.value.platform_logs_settings] : []
        content {
          severity = platform_logs_settings.value.severity != "" ? platform_logs_settings.value.severity : null
        }
      }
    }
  }

  # Ordered transform pipeline: transforms run in list order on every
  # published message; a disabled transform keeps its position (the
  # staging lever) without being applied.
  dynamic "message_transforms" {
    for_each = var.spec.message_transforms
    content {
      disabled = message_transforms.value.disabled
      javascript_udf {
        function_name = message_transforms.value.javascript_udf.function_name
        code          = message_transforms.value.javascript_udf.code
      }
    }
  }

  depends_on = [
    google_project_service.pubsub_api,
  ]
}
