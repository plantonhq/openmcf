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
  # The provider stores the plain config name at apply but the API's
  # read-back (refresh) returns the fully qualified
  # projects/{p}/instanceConfigs/{name} path, so exporting the attribute
  # verbatim drifts short -> full on the first refresh. The output
  # contract is the plain config name (what spec authors and API callers
  # use), so normalize to the last path segment — stable across both
  # forms and identical to the Pulumi module's export.
  value = element(reverse(split("/", google_spanner_instance.this.config)), 0)
}
