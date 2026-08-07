# Server-assigned schedule ID — the last path segment of the schedule's
# resource name, which is what Admin API calls address the schedule by.
output "schedule_id" {
  description = "Server-assigned backup schedule ID (last path segment)"
  value       = element(split("/", google_firestore_backup_schedule.this.name), 5)
}

# The database the schedule protects — confirms the parent without chasing
# the reference chain.
output "database" {
  description = "Firestore database name the schedule protects"
  value       = var.spec.database
}
