# Cloudflare Access service token: a machine credential (client-ID /
# client-secret pair) presented in the CF-Access-Client-ID / -Client-Secret
# headers. The secret is returned ONLY at creation and rotation -- capture the
# client_secret output at deploy time; it can never be read back later.
#
# Rotation: incrementing client_secret_version mints a new secret; the previous
# one keeps working until previous_client_secret_expires_at. The two travel as
# a pair (spec-enforced); both unset is the normal non-rotating state.
resource "cloudflare_zero_trust_access_service_token" "main" {
  account_id = local.account_id != "" ? local.account_id : null
  zone_id    = local.zone_id != "" ? local.zone_id : null

  name = var.spec.name

  # Empty means Cloudflare's default of one year (8760h).
  duration = var.spec.duration != "" ? var.spec.duration : null

  client_secret_version             = var.spec.client_secret_version
  previous_client_secret_expires_at = var.spec.previous_client_secret_expires_at != "" ? var.spec.previous_client_secret_expires_at : null
}
