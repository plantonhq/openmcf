# Enable the Cloud Monitoring API so a fresh project can host the channel.
# disable_on_destroy is false: tearing down one channel must never disable
# monitoring for everything else in the project.
resource "google_project_service" "monitoring_api" {
  project = local.project_id
  service = "monitoring.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Cloud Monitoring notification channel — the delivery endpoint alert
# policies notify when incidents open or close. Creating a channel sends
# nothing on its own; alert policies reference it by resource name (the
# channel_name output).
#
# Two label surfaces exist on this resource and they must never be
# conflated: `labels` is the TYPE-SPECIFIC channel configuration (an email
# address, a Slack channel name) fed from spec.channel_labels, while
# `user_labels` is freeform user metadata fed from spec.labels merged with
# the platform attribution labels. Credentials ride the separate
# sensitive_labels block, which GCP stores API-side and redacts on read.
#
# `enabled` is sent EXPLICITLY on every apply: it is Optional in the
# provider with a server default of true, and a spec transition
# true -> false must reach the API rather than being omitted (the
# send-true-or-omit class silently no-ops such transitions).
resource "google_monitoring_notification_channel" "this" {
  type         = var.spec.type
  display_name = local.display_name
  project      = local.project_id

  description = var.spec.description != "" ? var.spec.description : null

  # Type-specific NON-SECRET configuration; credentials are refused here by
  # the spec's validation and belong in sensitive_labels below.
  labels = length(var.spec.channel_labels) > 0 ? var.spec.channel_labels : null

  user_labels = local.final_labels

  enabled      = var.spec.enabled == null ? true : var.spec.enabled
  force_delete = var.spec.force_delete

  dynamic "sensitive_labels" {
    for_each = var.spec.sensitive_labels != null ? [var.spec.sensitive_labels] : []
    content {
      auth_token  = sensitive_labels.value.auth_token != "" ? sensitive_labels.value.auth_token : null
      password    = sensitive_labels.value.password != "" ? sensitive_labels.value.password : null
      service_key = sensitive_labels.value.service_key != "" ? sensitive_labels.value.service_key : null
    }
  }

  # Empty defers to the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.monitoring_api]
}
