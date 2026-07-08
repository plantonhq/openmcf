# Enable the Cloud Scheduler API — the control plane that owns jobs.
# disable_on_destroy is false: tearing down one job must never disable
# the API for everything else in the project (other jobs keep firing).
resource "google_project_service" "cloudscheduler_api" {
  project = local.project_id
  service = "cloudscheduler.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Cloud Scheduler job — a managed cron entry that fires on schedule and
# dispatches to exactly one target (HTTP endpoint, Pub/Sub topic, or App
# Engine handler). Name and region are ForceNew.
resource "google_cloud_scheduler_job" "this" {
  name     = local.job_name
  region   = local.location
  project  = local.project_id
  schedule = var.spec.schedule

  # time_zone defaults to Etc/UTC and attempt_deadline to 180s on the
  # provider — both are only sent when the manifest sets them, so the
  # provider defaults apply identically on both engines.
  time_zone        = var.spec.time_zone != "" ? var.spec.time_zone : null
  description      = var.spec.description != "" ? var.spec.description : null
  attempt_deadline = var.spec.attempt_deadline != "" ? var.spec.attempt_deadline : null

  # paused=true creates the job without firing it; false is the provider
  # default, so it is only sent when set (null otherwise) to avoid a
  # meaningless diff on the Optional+Computed attribute.
  paused = var.spec.paused ? true : null

  # HTTP target: the most common shape — trigger a Cloud Run service,
  # Cloud Function, or any HTTP endpoint, with OAuth (for
  # *.googleapis.com) XOR OIDC (for Cloud Run/Functions/custom) auth
  # enforced pre-deploy by the spec's CEL rule.
  dynamic "http_target" {
    for_each = var.spec.http_target != null ? [var.spec.http_target] : []
    content {
      uri         = http_target.value.uri
      http_method = http_target.value.http_method != "" ? http_target.value.http_method : null
      body        = http_target.value.body != "" ? http_target.value.body : null
      headers     = length(http_target.value.headers) > 0 ? http_target.value.headers : null

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
    }
  }

  # Pub/Sub target: publishes a message on schedule. topic_name arrives
  # resolved to the fully qualified projects/{project}/topics/{name} path
  # (the GcpPubSubTopic reference's topic_id output).
  dynamic "pubsub_target" {
    for_each = var.spec.pubsub_target != null ? [var.spec.pubsub_target] : []
    content {
      topic_name = pubsub_target.value.topic_name
      data       = pubsub_target.value.data != "" ? pubsub_target.value.data : null
      attributes = length(pubsub_target.value.attributes) > 0 ? pubsub_target.value.attributes : null
    }
  }

  # App Engine target: dispatches to an App Engine handler in the same
  # project; routing pins a specific service/version/instance.
  dynamic "app_engine_http_target" {
    for_each = var.spec.app_engine_http_target != null ? [var.spec.app_engine_http_target] : []
    content {
      relative_uri = app_engine_http_target.value.relative_uri
      http_method  = app_engine_http_target.value.http_method != "" ? app_engine_http_target.value.http_method : null
      body         = app_engine_http_target.value.body != "" ? app_engine_http_target.value.body : null
      headers      = length(app_engine_http_target.value.headers) > 0 ? app_engine_http_target.value.headers : null

      dynamic "app_engine_routing" {
        for_each = app_engine_http_target.value.app_engine_routing != null ? [app_engine_http_target.value.app_engine_routing] : []
        content {
          service  = app_engine_routing.value.service != "" ? app_engine_routing.value.service : null
          version  = app_engine_routing.value.version != "" ? app_engine_routing.value.version : null
          instance = app_engine_routing.value.instance != "" ? app_engine_routing.value.instance : null
        }
      }
    }
  }

  # Zero/empty means "not set" for these dials — the API then applies its
  # defaults. retry_count 0 genuinely means "fail after the first attempt"
  # in the API, but the provider treats an omitted field the same way, so
  # the sentinel is safe.
  dynamic "retry_config" {
    for_each = var.spec.retry_config != null ? [var.spec.retry_config] : []
    content {
      retry_count          = retry_config.value.retry_count != 0 ? retry_config.value.retry_count : null
      max_retry_duration   = retry_config.value.max_retry_duration != "" ? retry_config.value.max_retry_duration : null
      min_backoff_duration = retry_config.value.min_backoff_duration != "" ? retry_config.value.min_backoff_duration : null
      max_backoff_duration = retry_config.value.max_backoff_duration != "" ? retry_config.value.max_backoff_duration : null
      max_doublings        = retry_config.value.max_doublings != 0 ? retry_config.value.max_doublings : null
    }
  }

  depends_on = [
    google_project_service.cloudscheduler_api,
  ]
}
