# The full resource name
# (projects/<number>/locations/global/workloadIdentityPools/<pool_id>) — the
# handle IAM principals are built from and the parent providers attach to.
output "name" {
  description = "The full pool resource name (projects/<number>/locations/global/workloadIdentityPools/<pool_id>)"
  value       = google_iam_workload_identity_pool.this.name
}

# The bare pool ID, echoed for tooling that addresses the pool by short ID —
# providers reference the pool through this output.
output "workload_identity_pool_id" {
  description = "The pool ID (final component of the resource name)"
  value       = google_iam_workload_identity_pool.this.workload_identity_pool_id
}

# The pool lifecycle state: ACTIVE, or DELETED while soft-deleted (~30 days,
# during which the ID cannot be reused).
output "state" {
  description = "The pool lifecycle state (ACTIVE or DELETED)"
  value       = google_iam_workload_identity_pool.this.state
}
