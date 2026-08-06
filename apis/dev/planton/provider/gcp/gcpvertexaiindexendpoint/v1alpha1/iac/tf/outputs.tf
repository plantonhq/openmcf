output "index_endpoint_id" {
  description = "Fully qualified index endpoint resource path (projects/{project}/locations/{location}/indexEndpoints/{id}) — the value a GcpVertexAiDeployedIndex passes as its index_endpoint reference"
  value       = google_vertex_ai_index_endpoint.this.id
}

output "index_endpoint_name" {
  description = "The GCP-assigned numeric index endpoint ID (the last path segment of index_endpoint_id)"
  value       = google_vertex_ai_index_endpoint.this.name
}

output "public_endpoint_domain_name" {
  description = "Domain name for querying deployed indexes over the public internet (populated only when public_endpoint_enabled is true)"
  value       = google_vertex_ai_index_endpoint.this.public_endpoint_domain_name
}

output "create_time" {
  description = "RFC3339 timestamp of when the index endpoint was created"
  value       = google_vertex_ai_index_endpoint.this.create_time
}

output "update_time" {
  description = "RFC3339 timestamp of when the index endpoint was last updated"
  value       = google_vertex_ai_index_endpoint.this.update_time
}
