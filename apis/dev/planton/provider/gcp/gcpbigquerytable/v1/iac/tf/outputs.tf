output "table_id" {
  description = "The short table ID referenced by SQL queries and foreign keys"
  value       = google_bigquery_table.this.table_id
}

output "self_link" {
  description = "The fully qualified URI of the table"
  value       = google_bigquery_table.this.self_link
}

output "project" {
  description = "The GCP project that contains the table"
  # Read from the created resource so the output is correct under the
  # ambient-project fallback (the spec project may be empty).
  value = google_bigquery_table.this.project
}

output "dataset_id" {
  description = "The dataset that contains the table"
  value       = google_bigquery_table.this.dataset_id
}

output "type" {
  description = "The kind of table GCP materialized: TABLE, VIEW, MATERIALIZED_VIEW, or EXTERNAL"
  value       = google_bigquery_table.this.type
}

output "location" {
  description = "Geographic location of the table (inherited from the dataset)"
  value       = google_bigquery_table.this.location
}

output "creation_time" {
  description = "The creation time of the table in milliseconds since epoch"
  value       = google_bigquery_table.this.creation_time
}

# The dotted {project}.{dataset}.{table} handle, pre-assembled from the
# created resource's attributes (correct under the ambient-project fallback)
# so consumers that address tables in SQL-style dotted form (Pub/Sub
# BigQuery delivery, query tooling) never do string assembly.
output "qualified_name" {
  description = "The dotted fully qualified table name ({project}.{dataset}.{table})"
  value       = "${google_bigquery_table.this.project}.${google_bigquery_table.this.dataset_id}.${google_bigquery_table.this.table_id}"
}
