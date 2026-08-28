# Cloudflare account-level mTLS certificate upload: the trust material that
# Authenticated Origin Pulls rows, zone TLS CA associations, and Workers mTLS
# bindings reference. Every argument is create-only at the API -- any change
# replaces the upload and the certificate id changes, so consumers must
# re-point at the new id (rotate by replace, never in place).
resource "cloudflare_mtls_certificate" "main" {
  account_id = var.spec.account_id
  ca         = var.spec.ca
  # Canonicalized to end with exactly one trailing newline: Cloudflare stores
  # the PEM in that form and echoes it back on every read (measured
  # 2026-08-28), while the provider neither normalizes the echo nor updates
  # in place -- a non-canonical config would re-plan after import and its
  # apply can never converge ("Provider produced inconsistent result").
  certificates = "${trimspace(var.spec.certificates)}\n"

  name        = local.name
  private_key = local.private_key
}
