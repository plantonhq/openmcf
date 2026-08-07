# Fully qualified database ID — the canonical identifier for API calls
# and client library connections. The resource's own id attribute
# already carries the projects/{project}/databases/{name} form, with the
# project resolved even when the spec rode the provider default.
output "database_id" {
  description = "Fully qualified database ID (projects/{project}/databases/{name})"
  value       = google_firestore_database.this.id
}

output "database_name" {
  description = "Database name"
  value       = google_firestore_database.this.name
}

output "uid" {
  description = "Server-generated UUID4 for this database"
  value       = google_firestore_database.this.uid
}

output "create_time" {
  description = "Timestamp at which the database was created"
  value       = google_firestore_database.this.create_time
}

output "earliest_version_time" {
  description = "Earliest timestamp for point-in-time recovery reads"
  value       = google_firestore_database.this.earliest_version_time
}

# The retention window PITR restores and earliest-version reads can
# target (3600s without PITR, 604800s with it).
output "version_retention_period" {
  description = "How long past versions of data are retained"
  value       = google_firestore_database.this.version_retention_period
}

output "key_prefix" {
  description = "Key prefix for Datastore Mode app identifiers"
  value       = google_firestore_database.this.key_prefix
}

output "update_time" {
  description = "Timestamp of the database's last configuration update"
  value       = google_firestore_database.this.update_time
}
