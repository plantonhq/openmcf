output "self_link" {
  description = "Self-link URI of the backend service — the value URL maps reference as a default service or path-rule target."
  value       = google_compute_backend_service.this.self_link
}

output "backend_service_name" {
  description = "Name of the backend service as it exists in GCP."
  value       = google_compute_backend_service.this.name
}

output "generated_id" {
  description = "Server-assigned numeric ID of the backend service."
  value       = tostring(google_compute_backend_service.this.generated_id)
}

output "fingerprint" {
  description = "Server-computed fingerprint, used for optimistic concurrency control on out-of-band updates."
  value       = google_compute_backend_service.this.fingerprint
}
