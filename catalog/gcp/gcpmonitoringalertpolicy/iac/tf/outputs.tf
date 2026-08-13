# The server-assigned resource name of the policy — the handle for
# cross-referencing it in dashboards, snooze configurations, and the
# Monitoring API.
output "policy_name" {
  description = "Resource name of the alert policy (projects/{project}/alertPolicies/{id})"
  value       = google_monitoring_alert_policy.this.name
}
