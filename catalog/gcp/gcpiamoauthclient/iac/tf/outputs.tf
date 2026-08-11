# The system-generated OAuth client ID applications present in OAuth flows
# (distinct from the user-chosen resource ID).
output "client_id" {
  description = "The system-generated OAuth client ID"
  value       = google_iam_oauth_client.this.client_id
}

# The client's full resource name:
# projects/{project}/locations/{location}/oauthClients/{id}.
output "client_name" {
  description = "The full resource name of the OAuth client"
  value       = google_iam_oauth_client.this.name
}

# The client's lifecycle state.
output "state" {
  description = "The lifecycle state of the OAuth client"
  value       = google_iam_oauth_client.this.state
}

# The system-generated secret of the FIRST credential in spec.credentials
# (empty when no credentials are defined) — the single-credential case is
# the operating norm; rotation adds a second credential and swaps consumers
# over. A live, long-lived credential — marked sensitive so it never prints
# in plans.
output "client_secret" {
  description = "The client secret of the first credential (empty when no credentials are defined)"
  value = (
    length(var.spec.credentials) > 0
    ? google_iam_oauth_client_credential.this[var.spec.credentials[0].credential_id].client_secret
    : ""
  )
  sensitive = true
}
