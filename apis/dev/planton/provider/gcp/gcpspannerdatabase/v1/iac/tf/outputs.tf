output "database_id" {
  description = "Fully qualified database ID (projects/{project}/instances/{instance}/databases/{name})"
  # Built from the created resource's resolved attributes so the output is
  # correct under the ambient-project fallback (spec project may be empty).
  value = "projects/${google_spanner_database.this.project}/instances/${google_spanner_database.this.instance}/databases/${google_spanner_database.this.name}"
}

output "database_name" {
  description = "Short database name referenced by backup schedules"
  value       = google_spanner_database.this.name
}

output "state" {
  description = "Database state (CREATING or READY)"
  value       = google_spanner_database.this.state
}
