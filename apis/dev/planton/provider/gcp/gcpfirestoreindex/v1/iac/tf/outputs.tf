# Server-defined index resource name — the canonical identifier for
# Admin API calls.
output "index_id" {
  description = "Fully qualified index resource name"
  value       = google_firestore_index.this.name
}

# The collection (group) the index applies to — confirms the target without
# chasing the reference chain.
output "collection" {
  description = "Collection (group) ID the index applies to"
  value       = var.spec.collection
}
