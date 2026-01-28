# Cloudflare Load Balancer Monitor (Health Check)
# This is an account-level resource that probes origins
resource "cloudflare_load_balancer_monitor" "health_check" {
  type           = "https"
  method         = "GET"
  path           = local.health_probe_path
  expected_codes = "2xx"
  timeout        = 5
  interval       = 60  # 60 seconds minimum for Pro plan
  retries        = 2

  description = "Health check for ${local.resource_name}"
}

# Cloudflare Load Balancer Pool
# This is an account-level resource that groups origins
resource "cloudflare_load_balancer_pool" "main" {
  name    = "${local.resource_name}-pool"
  monitor = cloudflare_load_balancer_monitor.health_check.id
  enabled = true

  # Create an origin block for each origin in the spec
  dynamic "origins" {
    for_each = local.origins
    content {
      name    = origins.value.name
      address = origins.value.address
      enabled = true
      weight  = coalesce(origins.value.weight, 1)
    }
  }

  description = "Pool for ${local.resource_name}"
}

# Cloudflare Load Balancer
# This is a zone-level resource that ties everything together
resource "cloudflare_load_balancer" "main" {
  zone_id = local.zone_id
  name    = var.spec.hostname

  # Pool configuration
  default_pool_ids = [cloudflare_load_balancer_pool.main.id]
  fallback_pool_id = cloudflare_load_balancer_pool.main.id

  # Proxy configuration
  proxied = local.proxied

  # Traffic steering
  steering_policy = local.steering_policy

  # Session affinity
  session_affinity = local.session_affinity

  description = "Load balancer for ${var.spec.hostname}"

  # Wait for the pool to be healthy before creating the load balancer
  depends_on = [
    cloudflare_load_balancer_pool.main
  ]
}

