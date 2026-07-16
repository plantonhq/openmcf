output "self_link" {
  description = "Self-link URI of the network endpoint group — the value a backend service references in backends[].group."
  value       = google_compute_region_network_endpoint_group.this.self_link
}

output "network_endpoint_group_name" {
  description = "Name of the network endpoint group as it exists in GCP."
  value       = google_compute_region_network_endpoint_group.this.name
}

output "network_endpoint_type" {
  description = "The endpoint type of the NEG (SERVERLESS, PRIVATE_SERVICE_CONNECT, INTERNET_IP_PORT, INTERNET_FQDN_PORT, or GCE_VM_IP_PORTMAP)."
  value       = google_compute_region_network_endpoint_group.this.network_endpoint_type
}

output "region" {
  description = "Region the network endpoint group lives in."
  # Emit the spec region (plain name like us-central1), not the provider's
  # .region attribute which may be a region self-link on newer provider lines.
  value = var.spec.region
}
