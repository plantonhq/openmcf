# Exactly one of the two secret resources exists (see the count guards in
# main.tf), so each output selects whichever branch was created with
# one(concat(...)).

# The full resource name — the handle consumers use to reference the secret
# (e.g. a Cloud Run valueFromSecret mount).
output "secret_name" {
  description = "Full resource name of the secret"
  value       = one(concat(google_secret_manager_secret.this[*].name, google_secret_manager_regional_secret.this[*].name))
}

# The short secret ID (last segment of secret_name).
output "secret_id" {
  description = "The secret ID"
  value       = one(concat(google_secret_manager_secret.this[*].secret_id, google_secret_manager_regional_secret.this[*].secret_id))
}

# The version created from initial_version; empty when none was configured.
# Consumers pinning an exact version (instead of the "latest" alias)
# reference this value.
output "latest_version_name" {
  description = "Resource name of version 1 when initial_version was configured; empty otherwise"
  # Never coalesce() here — HCL's coalesce skips empty strings too and
  # errors when every argument is null/empty (the no-initial-version case).
  value = (
    var.spec.initial_version != null
    ? one(concat(google_secret_manager_secret_version.initial[*].name, google_secret_manager_regional_secret_version.initial[*].name))
    : ""
  )
}
