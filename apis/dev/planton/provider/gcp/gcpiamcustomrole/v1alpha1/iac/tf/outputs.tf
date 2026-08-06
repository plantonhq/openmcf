# The fully-qualified role name (projects/<project>/roles/<role_id>) — the
# grantable handle. Feed it directly into IAM grants (GcpProjectIamMember's
# role field) exactly as IAM policies expect it.
output "name" {
  description = "The fully-qualified role name (projects/<project>/roles/<role_id>)"
  value       = google_project_iam_custom_role.this.name
}

# The bare role ID within the project, echoed for tooling that addresses the
# role by short ID.
output "role_id" {
  description = "The role ID within the project"
  value       = google_project_iam_custom_role.this.role_id
}

# Whether the role is currently soft-deleted (GCP retains deleted roles for
# up to 14 days, during which grants are rejected).
output "deleted" {
  description = "Whether the role is currently soft-deleted"
  value       = google_project_iam_custom_role.this.deleted
}
