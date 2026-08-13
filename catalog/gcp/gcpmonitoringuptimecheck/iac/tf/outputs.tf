# The full resource name of the uptime check.
output "uptime_check_name" {
  description = "Resource name (projects/{project}/uptimeCheckConfigs/{id})"
  value       = google_monitoring_uptime_check_config.this.name
}

# The bare check ID — the value an alert policy's threshold filter
# references as metric.label.check_id to page on this check's failures.
output "uptime_check_id" {
  description = "The uptime check ID (last segment of the resource name)"
  value       = google_monitoring_uptime_check_config.this.uptime_check_id
}
