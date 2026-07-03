# The full provider resource name (projects/<number>/locations/global/
# workloadIdentityPools/<pool_id>/providers/<provider_id>) — the AUDIENCE for
# the keyless-auth handshake: OIDC tokens minted for this provider set `aud`
# to this value, and web-identity provider configurations consume exactly this
# string as their audience.
output "name" {
  description = "The full provider resource name — the audience string for token exchange"
  value       = google_iam_workload_identity_pool_provider.this.name
}

# The bare provider ID, echoed for tooling that addresses the provider by
# short ID.
output "workload_identity_pool_provider_id" {
  description = "The provider ID (final component of the resource name)"
  value       = google_iam_workload_identity_pool_provider.this.workload_identity_pool_provider_id
}

# The provider lifecycle state: ACTIVE, or DELETED while soft-deleted
# (~30 days, during which the ID cannot be reused).
output "state" {
  description = "The provider lifecycle state (ACTIVE or DELETED)"
  value       = google_iam_workload_identity_pool_provider.this.state
}
