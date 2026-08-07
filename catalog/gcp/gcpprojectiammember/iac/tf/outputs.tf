# The grant tuple, fully resolved (after reference resolution and any
# provider-default project fallback) — echoed so downstream tooling and audits
# see exactly what was applied without re-resolving references.
output "project_id" {
  description = "The project ID whose IAM policy received the grant"
  value       = google_project_iam_member.this.project
}

output "role" {
  description = "The role that was granted"
  value       = google_project_iam_member.this.role
}

output "member" {
  description = "The member the role was granted to, in IAM member format"
  value       = google_project_iam_member.this.member
}

# The etag of the project IAM policy after this grant was applied — a
# fingerprint of the policy version, useful for audit correlation.
output "etag" {
  description = "The etag of the project IAM policy after the grant"
  value       = google_project_iam_member.this.etag
}
