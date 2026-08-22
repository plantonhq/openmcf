# subscription.tf

# Zone plan subscription, provisioned only when the spec sets one. Creating a
# paid rate plan is a real billing action (the deploying token needs Billing
# Write scope). The spec's flat rate_plan/scope pair maps into the provider's
# nested rate_plan object; unset levers are sent as null so Cloudflare's
# defaults stay in effect.
resource "cloudflare_zone_subscription" "main" {
  count   = var.spec.subscription != null ? 1 : 0
  zone_id = cloudflare_zone.main.id

  frequency = var.spec.subscription.frequency != "" ? var.spec.subscription.frequency : null

  rate_plan = var.spec.subscription.rate_plan != "" ? {
    id    = var.spec.subscription.rate_plan
    scope = var.spec.subscription.scope != "" ? var.spec.subscription.scope : null
  } : null
}
