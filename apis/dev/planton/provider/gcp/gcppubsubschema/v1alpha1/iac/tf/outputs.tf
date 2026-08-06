# The resource ID is the fully qualified schema path
# (projects/{project}/schemas/{name}) — the exact string a topic's
# schema_settings.schema reference consumes.
output "schema_id" {
  description = "Fully qualified schema resource path"
  value       = google_pubsub_schema.schema.id
}

output "schema_name" {
  description = "The short name of the schema"
  value       = google_pubsub_schema.schema.name
}
