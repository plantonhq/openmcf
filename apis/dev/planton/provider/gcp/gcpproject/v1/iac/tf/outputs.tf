output "project_id" {
  description = "The unique project ID (globally unique across all GCP) — the value every other kind's project_id reference resolves"
  value       = google_project.this.project_id
}

output "project_number" {
  description = "The numeric identifier of the project (assigned by Google)"
  value       = google_project.this.number
}

output "name" {
  description = "The display name of the project"
  value       = google_project.this.name
}
