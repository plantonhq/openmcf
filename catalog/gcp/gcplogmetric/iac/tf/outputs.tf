# The metric name — address it from Cloud Monitoring (dashboards, alert
# policy filters) as metric.type = "logging.googleapis.com/user/{name}".
output "metric_name" {
  description = "The log-based metric name"
  value       = google_logging_metric.this.name
}
