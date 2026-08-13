# The service account email — the most common reference handle; workload
# configs (GKE, Cloud Run, Cloud Functions, Compute) attach identity by email.
output "email" {
  description = "The email address of the created service account"
  value       = google_service_account.main.email
}

# The IAM member string ("serviceAccount:<email>") in the exact format GCP IAM
# policies expect — feed it directly into IAM grants without string assembly.
output "member" {
  description = "The IAM member string for this service account (serviceAccount:<email>)"
  value       = google_service_account.main.member
}

# The stable numeric ID GCP assigns; never reused across delete/recreate,
# unlike the email — use where a tamper-proof identity reference matters.
output "unique_id" {
  description = "The stable, unique numeric ID of the service account"
  value       = google_service_account.main.unique_id
}

# The fully-qualified resource name (projects/<project>/serviceAccounts/<email>)
# used by APIs that address the account as a resource (keys, IAM-on-SA, WI bindings).
output "name" {
  description = "The fully-qualified resource name of the service account"
  value       = google_service_account.main.name
}

# The base64-encoded private key, populated only when user_managed_key
# requested the GENERATE flow (empty in the upload flow — GCP never returns
# a private key for an uploaded public key). A live, long-lived credential —
# marked sensitive so it never prints in plans.
output "key_base64" {
  description = "The base64-encoded private key for the service account (generate flow only)"
  value       = local.create_key ? google_service_account_key.main[0].private_key : null
  sensitive   = true
}
