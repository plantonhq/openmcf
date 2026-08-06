output "instance_id" {
  description = "Fully qualified instance ID (projects/{project}/instances/{name})"
  # Built from the created resource's resolved project so the output is
  # correct under the ambient-project fallback (spec project may be empty).
  value = "projects/${google_spanner_instance.this.project}/instances/${google_spanner_instance.this.name}"
}

output "instance_name" {
  description = "Short instance name referenced by databases and backup schedules"
  value       = google_spanner_instance.this.name
}

output "state" {
  description = "Instance state (CREATING or READY)"
  value       = google_spanner_instance.this.state
}

output "config" {
  description = "Instance configuration (geographic topology, e.g. regional-us-central1)"
  value       = google_spanner_instance.this.config
}
