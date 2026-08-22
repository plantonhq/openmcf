# DigitalOcean CDN
#
# Provisions a CDN endpoint fronting a Spaces bucket -- the complete
# digitalocean_cdn resource surface at schema v1. The certificate is wired
# through the provider's certificate_name argument ONLY: the numeric
# certificate_id argument is deprecated upstream (Let's Encrypt renewals
# rotate the UUID) and its update path can silently detach the certificate
# when the custom domain changes, so it is deliberately never rendered.
#
# The origin is create-only (changing it replaces the endpoint); ttl,
# certificate, and custom domain update in place. The provider retries
# reads on 404 for up to 30 seconds -- CDN creation is eventually
# consistent at the edge.

resource "digitalocean_cdn" "cdn" {
  # The origin reference resolves to the Space's fully-qualified domain
  # name before the module runs.
  origin = var.spec.origin

  # Unset defers to DigitalOcean's default of 3600 seconds, which the
  # provider reads back (Optional+Computed -- no perpetual diff). An
  # explicit zero cannot exist: spec validation floors ttl at 1 because
  # the provider drops zeros on the way out.
  ttl = var.spec.ttl

  # The certificate reference resolves to the certificate's stable NAME
  # (which is also its resource id at the current provider pin). The
  # "needs-cloudflare-cert" sentinel passes through verbatim.
  certificate_name = try(var.spec.certificate, "") != "" ? var.spec.certificate : null

  custom_domain = try(var.spec.custom_domain, "") != "" ? var.spec.custom_domain : null
}
