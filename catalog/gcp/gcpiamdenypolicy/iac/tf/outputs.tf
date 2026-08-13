# The policy's identifier as {url-encoded-parent}/{policy_name} — the
# handle gcloud and the v2 policies API reference the policy by.
output "policy_name" {
  description = "The identifier of the deny policy"
  value       = google_iam_deny_policy.this.id
}

# The policy's current etag — changes on every update; useful for
# optimistic-concurrency tooling reading the policy out of band.
output "etag" {
  description = "The etag of the deny policy"
  value       = google_iam_deny_policy.this.etag
}
