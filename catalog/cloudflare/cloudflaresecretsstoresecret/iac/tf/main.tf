# One secret inside the account Secrets Store. The value is write-only at
# Cloudflare (never returned, never drift-detected) and is marked sensitive
# in state. account_id, store_id, and name are create-only; value, scopes,
# and comment update in place (a merge-patch). The spec's CEL wall already
# guarantees scopes arrive in Cloudflare's canonical alphabetical order --
# the API returns them sorted, and an unsorted config would drift forever
# against the provider's ordered list.
resource "cloudflare_secrets_store_secret" "main" {
  account_id = var.spec.account_id
  store_id   = var.spec.store_id
  name       = var.spec.name
  value      = var.spec.value
  scopes     = var.spec.scopes

  comment = local.comment
}
