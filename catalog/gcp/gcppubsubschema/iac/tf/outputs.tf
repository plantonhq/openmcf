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

# A new revision is committed every time the definition changes. Topics
# consume this in schema_settings.first_revision_id / last_revision_id
# to pin validation to the revision this deploy produced.
output "revision_id" {
  description = "The current revision ID of the schema"
  value       = google_pubsub_schema.schema.revision_id
}
