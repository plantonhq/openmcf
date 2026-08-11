# The server-assigned resource name of the dashboard — the provider stores
# it as the resource ID (projects/{project}/dashboards/{dashboard_id}).
output "dashboard_name" {
  description = "Resource name of the dashboard (projects/{project}/dashboards/{id})"
  value       = google_monitoring_dashboard.this.id
}
