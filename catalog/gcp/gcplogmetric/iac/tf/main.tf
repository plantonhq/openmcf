# Enable the Cloud Logging API so a fresh project can host the metric.
# disable_on_destroy is false: tearing down one metric must never disable
# logging for everything else in the project.
resource "google_project_service" "logging_api" {
  project = local.project_id
  service = "logging.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Cloud Logging log-based metric — the bridge from log entries matching
# `filter` to a chartable Cloud Monitoring metric (a counter, or a
# distribution extracted from entry fields).
#
# `disabled` is sent EXPLICITLY on every apply: it is Optional in the
# provider, and a spec transition true -> false must reach the API rather
# than being omitted (the send-true-or-omit class silently no-ops such
# transitions — a metric that silently stays paused hides the incident it
# was built to reveal).
resource "google_logging_metric" "this" {
  name    = local.metric_name
  filter  = var.spec.filter
  project = local.project_id

  bucket_name = var.spec.bucket_name != "" ? var.spec.bucket_name : null
  description = var.spec.description != "" ? var.spec.description : null

  # Explicit send — see the resource comment.
  disabled = var.spec.disabled

  value_extractor  = var.spec.value_extractor != "" ? var.spec.value_extractor : null
  label_extractors = length(var.spec.label_extractors) > 0 ? var.spec.label_extractors : null

  dynamic "metric_descriptor" {
    for_each = var.spec.metric_descriptor != null ? [var.spec.metric_descriptor] : []
    content {
      metric_kind  = metric_descriptor.value.metric_kind
      value_type   = metric_descriptor.value.value_type
      unit         = metric_descriptor.value.unit != "" ? metric_descriptor.value.unit : null
      display_name = metric_descriptor.value.display_name != "" ? metric_descriptor.value.display_name : null

      dynamic "labels" {
        for_each = metric_descriptor.value.labels
        content {
          key         = labels.value.key
          description = labels.value.description != "" ? labels.value.description : null
          value_type  = labels.value.value_type != "" ? labels.value.value_type : null
        }
      }
    }
  }

  dynamic "bucket_options" {
    for_each = var.spec.bucket_options != null ? [var.spec.bucket_options] : []
    content {
      dynamic "explicit_buckets" {
        for_each = bucket_options.value.explicit_buckets != null ? [bucket_options.value.explicit_buckets] : []
        content {
          bounds = explicit_buckets.value.bounds
        }
      }

      dynamic "exponential_buckets" {
        for_each = bucket_options.value.exponential_buckets != null ? [bucket_options.value.exponential_buckets] : []
        content {
          num_finite_buckets = exponential_buckets.value.num_finite_buckets
          growth_factor      = exponential_buckets.value.growth_factor
          scale              = exponential_buckets.value.scale
        }
      }

      dynamic "linear_buckets" {
        for_each = bucket_options.value.linear_buckets != null ? [bucket_options.value.linear_buckets] : []
        content {
          num_finite_buckets = linear_buckets.value.num_finite_buckets
          offset             = linear_buckets.value.offset
          width              = linear_buckets.value.width
        }
      }
    }
  }

  # Empty defers to the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.logging_api]
}
