output "dataset_id" {
  description = "The short dataset ID referenced by SQL queries and downstream tables"
  value       = google_bigquery_dataset.this.dataset_id
}

output "self_link" {
  description = "The fully qualified URI of the dataset"
  value       = google_bigquery_dataset.this.self_link
}

output "project" {
  description = "The GCP project that contains this dataset"
  # Read from the created resource so the output is correct under the
  # ambient-project fallback (the spec project may be empty).
  value = google_bigquery_dataset.this.project
}

output "creation_time" {
  description = "The creation time of the dataset in milliseconds since epoch"
  value       = google_bigquery_dataset.this.creation_time
}

output "location" {
  description = "Geographic location of the dataset (tables joined in one query must share it)"
  value       = google_bigquery_dataset.this.location
}

output "etag" {
  description = "Entity tag, changing on every metadata modification"
  value       = google_bigquery_dataset.this.etag
}
