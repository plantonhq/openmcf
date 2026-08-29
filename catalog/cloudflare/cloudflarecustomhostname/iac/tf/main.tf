# A Cloudflare for SaaS custom hostname. The ssl object (with defaults coalesced and
# unset optionals omitted) is assembled in locals so both engines send identical
# values and rely on the provider's defaults for everything left unset.
resource "cloudflare_custom_hostname" "main" {
  zone_id              = local.zone_id
  hostname             = var.spec.hostname
  custom_origin_server = local.custom_origin_server
  custom_origin_sni    = local.custom_origin_sni
  custom_metadata      = local.custom_metadata
  ssl                  = local.ssl

  # ssl.certificate_authority cannot converge when unset: Cloudflare assigns a
  # CA server-side AT RANDOM (measured live 2026-08-29: ssl_com then google on
  # consecutive creates), writing one is Enterprise-gated (400 code 1459
  # "Certificate Authority selection is only available on an Enterprise plan",
  # probe-measured), and the provider schema ships no state-preserving plan
  # modifier on the attribute -- so every refresh-inclusive plan on an
  # ssl-bearing config proposes an update forever. The provider's OWN
  # acceptance tests ignore_changes this exact attribute; we mirror that
  # recipe (scoped to the one attribute, same in the Pulumi module).
  # Trade-off, documented on the spec field: an in-place change of
  # certificate_authority after create is NOT applied by this module --
  # Enterprise users changing CA must recreate the hostname.
  lifecycle {
    ignore_changes = [ssl.certificate_authority]
  }
}
