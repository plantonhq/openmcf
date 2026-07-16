# The resource ID is the fully qualified topic path
# (projects/{project}/topics/{name}) — the exact string subscriptions,
# Cloud Functions event triggers, and Scheduler pubsub targets consume.
output "topic_id" {
  description = "The fully qualified topic ID (projects/{project}/topics/{name})"
  value       = google_pubsub_topic.this.id
}

output "topic_name" {
  description = "The short topic name"
  value       = google_pubsub_topic.this.name
}
