# Cloudflare AI Gateway: the caching, rate-limiting, guardrails,
# spend-limits, and routing control plane in front of AI model traffic. The
# gateway's id is the user-chosen URL slug (create-only: renaming replaces
# the gateway and changes the endpoint URL every client calls). The five
# scalar arguments below the ids are Required by Cloudflare with no
# defaults -- the spec forces an explicit choice for each.
resource "cloudflare_ai_gateway" "main" {
  account_id = var.spec.account_id
  id         = var.spec.gateway_id

  cache_invalidate_on_update = var.spec.cache_invalidate_on_update
  cache_ttl                  = var.spec.cache_ttl
  collect_logs               = var.spec.collect_logs
  rate_limiting_interval     = var.spec.rate_limiting_interval
  rate_limiting_limit        = var.spec.rate_limiting_limit

  rate_limiting_technique = local.rate_limiting_technique

  # The spec's retry{} / log_management{} groups fan out to the provider's
  # flat arguments (see locals.tf).
  retry_backoff           = local.retry_backoff
  retry_delay             = local.retry_delay
  retry_max_attempts      = local.retry_max_attempts
  log_management          = local.log_management
  log_management_strategy = local.log_management_strategy

  authentication          = var.spec.authentication
  logpush                 = var.spec.logpush
  logpush_public_key      = local.logpush_public_key
  zdr                     = var.spec.zdr
  workers_ai_billing_mode = local.workers_ai_billing_mode
  store_id                = local.store_id

  dlp          = local.dlp
  guardrails   = local.guardrails
  otel         = local.otel
  stripe       = local.stripe
  spend_limits = local.spend_limits
}

# One dynamic-routing object per dynamic_routes row, attached to the gateway
# above. The provider forces replacement on any change to a route's elements
# list -- editing a graph recreates that route object (requests re-resolve on
# the next call), while renaming a route is its only in-place update. Keyed
# by route name so a graph edit replaces only its own route.
resource "cloudflare_ai_gateway_dynamic_routing" "routes" {
  for_each = local.dynamic_routes

  account_id = var.spec.account_id
  gateway_id = cloudflare_ai_gateway.main.id
  name       = each.key
  elements   = each.value
}
