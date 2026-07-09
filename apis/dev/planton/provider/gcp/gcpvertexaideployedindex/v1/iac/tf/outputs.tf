output "name" {
  description = "Name of the DeployedIndex resource as the provider reports it"
  value       = google_vertex_ai_index_endpoint_deployed_index.this.name
}

output "deployed_index_id" {
  description = "The user-chosen deployment ID"
  value       = google_vertex_ai_index_endpoint_deployed_index.this.deployed_index_id
}

output "create_time" {
  description = "RFC3339 timestamp of when the deployment was created"
  value       = google_vertex_ai_index_endpoint_deployed_index.this.create_time
}

output "index_sync_time" {
  description = "RFC3339 timestamp up to which this deployment reflects the source index's updates"
  value       = google_vertex_ai_index_endpoint_deployed_index.this.index_sync_time
}

output "match_grpc_address" {
  description = "Private gRPC address for match queries inside the peered VPC (peered endpoints only)"
  value       = try(google_vertex_ai_index_endpoint_deployed_index.this.private_endpoints[0].match_grpc_address, "")
}

output "service_attachment" {
  description = "PSC service attachment consumers target with forwarding rules (PSC endpoints only)"
  value       = try(google_vertex_ai_index_endpoint_deployed_index.this.private_endpoints[0].service_attachment, "")
}

output "index_endpoint" {
  description = "Fully qualified resource path of the index endpoint this deployment lives on — query clients need it together with deployed_index_id"
  value       = google_vertex_ai_index_endpoint_deployed_index.this.index_endpoint
}
