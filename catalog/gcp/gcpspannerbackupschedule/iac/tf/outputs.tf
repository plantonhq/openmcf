output "schedule_id" {
  description = "Fully qualified backup schedule ID (projects/{project}/instances/{instance}/databases/{database}/backupSchedules/{name})"
  # Built from the created resource's resolved attributes so the output is
  # correct under the ambient-project fallback (spec project may be empty).
  value = "projects/${google_spanner_backup_schedule.this.project}/instances/${google_spanner_backup_schedule.this.instance}/databases/${google_spanner_backup_schedule.this.database}/backupSchedules/${google_spanner_backup_schedule.this.name}"
}

output "schedule_name" {
  description = "Short backup schedule name within the database"
  value       = google_spanner_backup_schedule.this.name
}
