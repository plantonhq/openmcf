output "index_id" {
  description = "Fully qualified index resource path (projects/{project}/locations/{location}/indexes/{indexId}) — the value a GcpVertexAiDeployedIndex passes as its index reference"
  value       = google_vertex_ai_index.this.id
}

output "index_name" {
  description = "The GCP-assigned numeric index ID (the last path segment of index_id)"
  value       = google_vertex_ai_index.this.name
}

output "metadata_schema_uri" {
  description = "Cloud Storage URI of the YAML schema describing additional index-specific information"
  value       = google_vertex_ai_index.this.metadata_schema_uri
}

output "create_time" {
  description = "RFC3339 timestamp of when the index was created"
  value       = google_vertex_ai_index.this.create_time
}

output "update_time" {
  description = "RFC3339 timestamp of when the index was last updated"
  value       = google_vertex_ai_index.this.update_time
}
