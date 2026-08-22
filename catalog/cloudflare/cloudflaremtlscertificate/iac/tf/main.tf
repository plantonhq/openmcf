# Cloudflare account-level mTLS certificate upload: the trust material that
# Authenticated Origin Pulls rows, zone TLS CA associations, and Workers mTLS
# bindings reference. Every argument is create-only at the API -- any change
# replaces the upload and the certificate id changes, so consumers must
# re-point at the new id (rotate by replace, never in place).
resource "cloudflare_mtls_certificate" "main" {
  account_id   = var.spec.account_id
  ca           = var.spec.ca
  certificates = var.spec.certificates

  name        = local.name
  private_key = local.private_key
}
