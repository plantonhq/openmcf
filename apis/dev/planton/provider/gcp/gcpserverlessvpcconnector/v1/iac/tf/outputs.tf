output "name" {
  description = "Short name of the connector as created in GCP"
  value       = google_vpc_access_connector.main.name
}

output "self_link" {
  description = "Fully qualified resource name (projects/*/locations/*/connectors/*) serverless workloads attach to"
  value       = google_vpc_access_connector.main.self_link
}

output "state" {
  description = "State of the connector (READY, CREATING, DELETING, ERROR, UPDATING)"
  value       = google_vpc_access_connector.main.state
}

output "region" {
  # Emit the plain spec region, not the resource attribute: on released 6.x
  # provider lines regional attributes can surface as region self-links
  # rather than plain names, and the output proto documents a plain name.
  description = "Region the connector lives in (plain region name)"
  value       = var.spec.region
}
