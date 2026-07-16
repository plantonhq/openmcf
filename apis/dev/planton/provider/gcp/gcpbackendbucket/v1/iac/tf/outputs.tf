# The self-link — the value URL maps reference as a default service or
# path-rule target; the composition handle for routing static content.
output "self_link" {
  description = "Self-link URI of the backend bucket"
  value       = google_compute_backend_bucket.this.self_link
}

# The name as it exists in GCP.
output "backend_bucket_name" {
  description = "Name of the backend bucket"
  value       = google_compute_backend_bucket.this.name
}

# The origin bucket currently being served, echoed for tooling that
# resolves the serving chain.
output "bucket_name" {
  description = "Name of the Cloud Storage bucket serving as the origin"
  value       = google_compute_backend_bucket.this.bucket_name
}
