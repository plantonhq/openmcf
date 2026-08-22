# A single DNS record inside an existing DigitalOcean DNS zone. The per-type
# fields (priority/weight/port/flags/tag) carry the spec's presence semantics:
# unset arrives as null and is simply not sent, matching the provider's own
# GetOk-guarded request building. Spec CEL rules already guarantee the fields
# each record type requires are present.
resource "digitalocean_record" "dns_record" {
  domain = local.domain
  type   = local.type
  name   = var.spec.name
  value  = local.value

  # When null, the ttl attribute is Computed: DigitalOcean applies its
  # default (1800 seconds) and the applied value reads back into state.
  ttl = var.spec.ttl_seconds

  priority = var.spec.priority
  weight   = var.spec.weight
  port     = var.spec.port
  flags    = var.spec.flags
  tag      = var.spec.tag != "" ? var.spec.tag : null
}
