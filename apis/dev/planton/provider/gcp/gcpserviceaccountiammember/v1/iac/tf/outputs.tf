# The grant tuple, fully resolved (after reference resolution) — echoed so
# downstream tooling and audits see exactly what was applied without
# re-resolving references.
output "service_account_id" {
  description = "The fully-qualified resource name of the service account whose IAM policy received the grant"
  value       = google_service_account_iam_member.this.service_account_id
}

output "role" {
  description = "The role that was granted"
  value       = google_service_account_iam_member.this.role
}

output "member" {
  description = "The member the role was granted to, in IAM member format"
  value       = google_service_account_iam_member.this.member
}

# The etag of the service account IAM policy after this grant was applied — a
# fingerprint of the policy version, useful for audit correlation.
output "etag" {
  description = "The etag of the service account IAM policy after the grant"
  value       = google_service_account_iam_member.this.etag
}
