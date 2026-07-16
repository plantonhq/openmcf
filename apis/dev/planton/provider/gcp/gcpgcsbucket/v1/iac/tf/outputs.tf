# For GCS the resource ID equals the globally unique bucket name — the
# value every consumer (backend buckets, function sources, Dataproc
# staging, Pub/Sub sinks) references.
output "bucket_id" {
  description = "ID of the bucket (equals the bucket name)"
  value       = google_storage_bucket.this.id
}

output "bucket_name" {
  description = "Name of the bucket"
  value       = google_storage_bucket.this.name
}

output "url" {
  description = "Base URI of the bucket (gs://<name>)"
  value       = google_storage_bucket.this.url
}

output "self_link" {
  description = "API self link of the bucket"
  value       = google_storage_bucket.this.self_link
}

# GCS reports location upper-cased regardless of input case.
output "location" {
  description = "Location of the bucket as reported by GCS"
  value       = google_storage_bucket.this.location
}

output "project_number" {
  description = "Numeric project number of the owning project"
  value       = google_storage_bucket.this.project_number
}
