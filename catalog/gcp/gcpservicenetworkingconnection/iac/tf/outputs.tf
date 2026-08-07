# Name of the VPC peering GCP created on the network (e.g.
# servicenetworking-googleapis-com) — what an operator looks for on the
# network's peerings list when auditing private services access.
output "peering" {
  description = "Name of the VPC peering created for this connection"
  value       = google_service_networking_connection.this.peering
}

# The peered network as the connection resolved it — confirms which network
# the producer is attached to without chasing the reference chain.
output "network" {
  description = "Self-link (or name) of the peered VPC network"
  value       = google_service_networking_connection.this.network
}
