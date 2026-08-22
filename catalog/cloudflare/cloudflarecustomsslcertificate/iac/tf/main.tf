# Cloudflare custom SSL certificate: a bring-your-own certificate uploaded to
# the zone. Custom certificates are a Business/Enterprise zone feature --
# Cloudflare enforces the plan gate at create. Rotation is replacement:
# changing the certificate or private key destroys and re-creates the upload
# (Cloudflare serves the previous certificate until the replacement deploys).
#
# priority is deliberately absent: at provider v5.23.0 it is read-only (the
# v4 reprioritization surface was dropped), so certificate priority cannot be
# managed here.
resource "cloudflare_custom_ssl" "main" {
  zone_id     = var.spec.zone_id
  certificate = var.spec.certificate
  private_key = var.spec.private_key

  type          = var.spec.type
  bundle_method = var.spec.bundle_method
  deploy        = var.spec.deploy

  policy           = local.policy
  custom_csr_id    = local.custom_csr_id
  geo_restrictions = local.geo_restrictions
}
