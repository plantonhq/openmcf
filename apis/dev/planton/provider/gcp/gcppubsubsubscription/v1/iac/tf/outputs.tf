# The resource ID is the fully qualified subscription path
# (projects/{project}/subscriptions/{name}) — the handle consumers and
# monitoring reference.
output "subscription_id" {
  description = "The fully qualified subscription ID (projects/{project}/subscriptions/{name})"
  value       = google_pubsub_subscription.this.id
}

output "subscription_name" {
  description = "The short subscription name"
  value       = google_pubsub_subscription.this.name
}
