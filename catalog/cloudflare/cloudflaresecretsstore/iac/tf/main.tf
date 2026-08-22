# Cloudflare account-level Secrets Store: the vault that store secrets,
# Worker bindings, and AI Gateway authentication consume from. Both arguments
# are create-only at the API (the provider's Update is an empty stub and
# every field forces replacement) -- a name change replaces the store AND
# every secret inside it. Cloudflare also allows only one store per account:
# if one already exists (e.g. dashboard-created), this create fails and the
# existing store should be adopted by import instead.
resource "cloudflare_secrets_store" "main" {
  account_id = var.spec.account_id
  name       = var.spec.name
}
