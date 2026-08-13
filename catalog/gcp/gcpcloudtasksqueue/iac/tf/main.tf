# Enable the Cloud Tasks API — the control plane that owns queues.
# disable_on_destroy is false: tearing down one queue must never disable
# the API for everything else in the project (other queues keep dispatching).
resource "google_project_service" "cloudtasks_api" {
  project = local.project_id
  service = "cloudtasks.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Cloud Tasks queue — the dispatch-rate, retry, and routing contract for
# every task enqueued into it. Name and location are ForceNew, and a deleted
# queue's ID stays reserved by the API for up to 7 days, so renames both
# replace the queue and burn the old identifier for that window.
resource "google_cloud_tasks_queue" "this" {
  name     = local.queue_name
  location = local.location
  project  = local.project_id

  # Declarative pause/resume: the provider issues the pause or resume
  # call whenever this differs from the queue's live state. Sent
  # explicitly so a PAUSED -> RUNNING spec edit resumes dispatch.
  desired_state = var.spec.desired_state != "" ? var.spec.desired_state : "RUNNING"

  # DELETE (provider default) removes the queue and its backlog on
  # destroy; PREVENT fails the destroy; ABANDON leaves the queue running.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  # Queue-level HTTP task settings. These OVERRIDE task-level configuration
  # at dispatch time — the pattern that lets producers enqueue bare payloads
  # while the queue owns auth and routing.
  dynamic "http_target" {
    for_each = var.spec.http_target != null ? [var.spec.http_target] : []
    content {
      http_method = http_target.value.http_method != "" ? http_target.value.http_method : null

      dynamic "header_overrides" {
        for_each = http_target.value.header_overrides
        content {
          header {
            key   = header_overrides.value.key
            value = header_overrides.value.value
          }
        }
      }

      # OAuth (for *.googleapis.com targets) and OIDC (for Cloud Run /
      # Cloud Functions / custom endpoints) are mutually exclusive —
      # enforced pre-deploy by the spec's CEL rule.
      dynamic "oauth_token" {
        for_each = http_target.value.oauth_token != null ? [http_target.value.oauth_token] : []
        content {
          service_account_email = oauth_token.value.service_account_email
          scope                 = oauth_token.value.scope != "" ? oauth_token.value.scope : null
        }
      }

      dynamic "oidc_token" {
        for_each = http_target.value.oidc_token != null ? [http_target.value.oidc_token] : []
        content {
          service_account_email = oidc_token.value.service_account_email
          audience              = oidc_token.value.audience != "" ? oidc_token.value.audience : null
        }
      }

      # The spec flattens the provider's nested path_override/query_override
      # single-field blocks into plain path/query_params strings; the blocks
      # are only sent when set, which also avoids the provider's
      # query_override perpetual-diff on the 6.x line when the block would
      # otherwise be sent empty.
      dynamic "uri_override" {
        for_each = http_target.value.uri_override != null ? [http_target.value.uri_override] : []
        content {
          scheme                    = uri_override.value.scheme != "" ? uri_override.value.scheme : null
          host                      = uri_override.value.host != "" ? uri_override.value.host : null
          port                      = uri_override.value.port != "" ? uri_override.value.port : null
          uri_override_enforce_mode = uri_override.value.enforce_mode != "" ? uri_override.value.enforce_mode : null

          dynamic "path_override" {
            for_each = uri_override.value.path != "" ? [uri_override.value.path] : []
            content {
              path = path_override.value
            }
          }

          dynamic "query_override" {
            for_each = uri_override.value.query_params != "" ? [uri_override.value.query_params] : []
            content {
              query_params = query_override.value
            }
          }
        }
      }
    }
  }

  # Routing override for App Engine tasks: pins the whole queue's App Engine
  # tasks to one service/version/instance instead of per-task routing.
  # Ignored for HTTP tasks.
  dynamic "app_engine_routing_override" {
    for_each = var.spec.app_engine_routing_override != null ? [var.spec.app_engine_routing_override] : []
    content {
      service  = app_engine_routing_override.value.service != "" ? app_engine_routing_override.value.service : null
      version  = app_engine_routing_override.value.version != "" ? app_engine_routing_override.value.version : null
      instance = app_engine_routing_override.value.instance != "" ? app_engine_routing_override.value.instance : null
    }
  }

  # Zero means "not set" for these dials — Cloud Tasks then applies its own
  # defaults (which it also does for a wholly omitted block).
  dynamic "rate_limits" {
    for_each = var.spec.rate_limits != null ? [var.spec.rate_limits] : []
    content {
      max_dispatches_per_second = rate_limits.value.max_dispatches_per_second > 0 ? rate_limits.value.max_dispatches_per_second : null
      max_concurrent_dispatches = rate_limits.value.max_concurrent_dispatches > 0 ? rate_limits.value.max_concurrent_dispatches : null
    }
  }

  # max_attempts genuinely distinguishes -1 (unlimited) from unset, so the
  # not-set sentinel is 0, never -1.
  dynamic "retry_config" {
    for_each = var.spec.retry_config != null ? [var.spec.retry_config] : []
    content {
      max_attempts       = retry_config.value.max_attempts != 0 ? retry_config.value.max_attempts : null
      max_retry_duration = retry_config.value.max_retry_duration != "" ? retry_config.value.max_retry_duration : null
      min_backoff        = retry_config.value.min_backoff != "" ? retry_config.value.min_backoff : null
      max_backoff        = retry_config.value.max_backoff != "" ? retry_config.value.max_backoff : null
      max_doublings      = retry_config.value.max_doublings != 0 ? retry_config.value.max_doublings : null
    }
  }

  # sampling_ratio 0.0 is a meaningful value (log nothing), so the block is
  # driven purely by presence.
  dynamic "stackdriver_logging_config" {
    for_each = var.spec.stackdriver_logging_config != null ? [var.spec.stackdriver_logging_config] : []
    content {
      sampling_ratio = stackdriver_logging_config.value.sampling_ratio
    }
  }

  depends_on = [
    google_project_service.cloudtasks_api,
  ]
}
