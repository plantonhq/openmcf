# Cloudflare standalone health check: scheduled origin probes with
# healthy/unhealthy status, no load balancer required. Health checks are a paid
# zone feature -- Cloudflare enforces the plan gate at create.
#
# The unused config block is never sent (the spec's CEL guarantees http_config
# only on HTTP/HTTPS and tcp_config only on TCP), because both blocks are
# Computed upstream and sending the wrong one reads back as drift.
resource "cloudflare_healthcheck" "main" {
  zone_id = var.spec.zone_id
  name    = var.spec.name
  address = var.spec.address

  type = var.spec.type

  check_regions = length(var.spec.check_regions) > 0 ? var.spec.check_regions : null

  consecutive_fails     = var.spec.consecutive_fails
  consecutive_successes = var.spec.consecutive_successes
  interval              = var.spec.interval
  retries               = var.spec.retries
  timeout               = var.spec.timeout
  suspended             = var.spec.suspended

  http_config = local.http_config
  tcp_config  = local.tcp_config
}
