# Enable the Cloud Monitoring API so a fresh project can host the
# dashboard. disable_on_destroy is false: tearing down one dashboard must
# never disable monitoring for everything else in the project.
resource "google_project_service" "monitoring_api" {
  project = local.project_id
  service = "monitoring.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Cloud Monitoring dashboard from the spec's one JSON document (the
# Monitoring API's own Dashboard format — the provider deliberately models
# the fast-moving widget schema as a JSON string, and this module honors
# that judgment). The provider validates the document is JSON at plan time
# and suppresses diffs on server-added keys (etag, name), so a dashboard
# exported from the GCP console round-trips cleanly.
resource "google_monitoring_dashboard" "this" {
  dashboard_json = var.spec.dashboard_json
  project        = local.project_id

  # Empty defers to the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.monitoring_api]
}
