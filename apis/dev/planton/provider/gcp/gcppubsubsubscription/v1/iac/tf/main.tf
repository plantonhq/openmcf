# Enable the Pub/Sub API — the control plane that owns the subscription.
# disable_on_destroy is false: tearing down one subscription must never
# disable the API for everything else in the project.
resource "google_project_service" "pubsub_api" {
  project = local.project_id
  service = "pubsub.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Pub/Sub subscription — the named message stream from one topic to
# one consumer. name, topic, filter, and enable_message_ordering are
# ForceNew: changing any of them replaces the subscription (and its
# backlog), so treat them as permanent.
resource "google_pubsub_subscription" "this" {
  name    = local.subscription_name
  topic   = local.topic
  project = local.project_id
  labels  = local.final_labels

  # 0/empty accept the API defaults (10s deadline, 7-day retention).
  ack_deadline_seconds         = local.ack_deadline_seconds
  message_retention_duration   = local.message_retention_duration
  retain_acked_messages        = var.spec.retain_acked_messages
  filter                       = local.filter
  enable_message_ordering      = var.spec.enable_message_ordering
  enable_exactly_once_delivery = var.spec.enable_exactly_once_delivery

  # An empty ttl string means "never expires" — sent as-is; the API
  # treats empty as the never-expire sentinel.
  dynamic "expiration_policy" {
    for_each = var.spec.expiration_policy != null ? [var.spec.expiration_policy] : []
    content {
      ttl = expiration_policy.value.ttl
    }
  }

  # Dead-lettering: the Pub/Sub service agent needs Subscriber on this
  # subscription and Publisher on the dead-letter topic.
  dynamic "dead_letter_policy" {
    for_each = var.spec.dead_letter_policy != null ? [var.spec.dead_letter_policy] : []
    content {
      dead_letter_topic     = dead_letter_policy.value.dead_letter_topic != "" ? dead_letter_policy.value.dead_letter_topic : null
      max_delivery_attempts = dead_letter_policy.value.max_delivery_attempts > 0 ? dead_letter_policy.value.max_delivery_attempts : null
    }
  }

  dynamic "retry_policy" {
    for_each = var.spec.retry_policy != null ? [var.spec.retry_policy] : []
    content {
      minimum_backoff = retry_policy.value.minimum_backoff != "" ? retry_policy.value.minimum_backoff : null
      maximum_backoff = retry_policy.value.maximum_backoff != "" ? retry_policy.value.maximum_backoff : null
    }
  }

  # Delivery mode: at most one of push/bigquery/cloud-storage (spec CEL
  # enforces exclusivity); none set = pull.
  dynamic "push_config" {
    for_each = var.spec.push_config != null ? [var.spec.push_config] : []
    content {
      push_endpoint = push_config.value.push_endpoint
      attributes    = length(push_config.value.attributes) > 0 ? push_config.value.attributes : null

      dynamic "oidc_token" {
        for_each = push_config.value.oidc_token != null ? [push_config.value.oidc_token] : []
        content {
          service_account_email = oidc_token.value.service_account_email
          audience              = oidc_token.value.audience != "" ? oidc_token.value.audience : null
        }
      }

      dynamic "no_wrapper" {
        for_each = push_config.value.no_wrapper != null ? [push_config.value.no_wrapper] : []
        content {
          write_metadata = no_wrapper.value.write_metadata
        }
      }
    }
  }

  dynamic "bigquery_config" {
    for_each = var.spec.bigquery_config != null ? [var.spec.bigquery_config] : []
    content {
      # The table arrives resolved to {project}.{dataset}.{table}.
      table                 = bigquery_config.value.table
      use_topic_schema      = bigquery_config.value.use_topic_schema
      use_table_schema      = bigquery_config.value.use_table_schema
      drop_unknown_fields   = bigquery_config.value.drop_unknown_fields
      write_metadata        = bigquery_config.value.write_metadata
      service_account_email = bigquery_config.value.service_account_email != "" ? bigquery_config.value.service_account_email : null
    }
  }

  dynamic "cloud_storage_config" {
    for_each = var.spec.cloud_storage_config != null ? [var.spec.cloud_storage_config] : []
    content {
      bucket                   = cloud_storage_config.value.bucket
      filename_prefix          = cloud_storage_config.value.filename_prefix != "" ? cloud_storage_config.value.filename_prefix : null
      filename_suffix          = cloud_storage_config.value.filename_suffix != "" ? cloud_storage_config.value.filename_suffix : null
      filename_datetime_format = cloud_storage_config.value.filename_datetime_format != "" ? cloud_storage_config.value.filename_datetime_format : null
      max_bytes                = cloud_storage_config.value.max_bytes > 0 ? cloud_storage_config.value.max_bytes : null
      max_duration             = cloud_storage_config.value.max_duration != "" ? cloud_storage_config.value.max_duration : null
      max_messages             = cloud_storage_config.value.max_messages > 0 ? cloud_storage_config.value.max_messages : null
      service_account_email    = cloud_storage_config.value.service_account_email != "" ? cloud_storage_config.value.service_account_email : null

      dynamic "avro_config" {
        for_each = cloud_storage_config.value.avro_config != null ? [cloud_storage_config.value.avro_config] : []
        content {
          use_topic_schema = avro_config.value.use_topic_schema
          write_metadata   = avro_config.value.write_metadata
        }
      }
    }
  }

  # Ordered transform pipeline: transforms run in list order on every
  # message before delivery to THIS subscription only; a disabled
  # transform keeps its position (the staging lever) without being
  # applied.
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
