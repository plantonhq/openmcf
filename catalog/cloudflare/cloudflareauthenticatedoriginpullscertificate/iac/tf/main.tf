# Cloudflare Authenticated Origin Pulls client certificate upload. The scope
# selects the API surface: zone replaces the zone-wide client certificate;
# hostname uploads a certificate for per-hostname associations to reference.
# Rotation is replacement on both surfaces. Never rotate only the private key
# against the same certificate: the zone-scoped API silently ignores a
# key-only change (a measured provider defect at v5.23.0 -- its Update is
# empty and the key does not force replacement), so key and certificate must
# always change together.

resource "cloudflare_authenticated_origin_pulls_certificate" "zone" {
  count = local.is_zone_scope ? 1 : 0

  zone_id     = var.spec.zone_id
  certificate = var.spec.certificate
  private_key = var.spec.private_key
}

resource "cloudflare_authenticated_origin_pulls_hostname_certificate" "hostname" {
  count = local.is_zone_scope ? 0 : 1

  zone_id     = var.spec.zone_id
  certificate = var.spec.certificate
  private_key = var.spec.private_key
}
