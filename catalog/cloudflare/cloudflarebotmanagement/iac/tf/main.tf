# Cloudflare Bot Management: the zone's singleton bot-configuration surface.
# Create ADOPTS whatever configuration the zone already carries (the API's
# create is a PUT), unset fields are never sent (the zone keeps its current
# values), and destroy is a NO-OP at Cloudflare -- the state entry disappears
# but the live configuration stays exactly as last applied. To retire a
# setting, apply its off value BEFORE destroying.
#
# Plan gates are Cloudflare's, enforced at the API: fields the zone's plan does
# not include (SBFM on free plans, the Enterprise knobs) fail at apply, and on
# non-entitled zones the API omits those fields from responses entirely --
# the provider's own issue tracker records refresh drift when a manifest sets
# fields the zone is not entitled to. Manage only what the plan includes.
resource "cloudflare_bot_management" "main" {
  zone_id = var.spec.zone_id

  fight_mode                      = var.spec.fight_mode
  sbfm_definitely_automated       = var.spec.sbfm_definitely_automated
  sbfm_likely_automated           = var.spec.sbfm_likely_automated
  sbfm_verified_bots              = var.spec.sbfm_verified_bots
  sbfm_static_resource_protection = var.spec.sbfm_static_resource_protection
  optimize_wordpress              = var.spec.optimize_wordpress
  auto_update_model               = var.spec.auto_update_model
  suppress_session_score          = var.spec.suppress_session_score
  enable_js                       = var.spec.enable_js
  bm_cookie_enabled               = var.spec.bm_cookie_enabled
  ai_bots_protection              = var.spec.ai_bots_protection
  crawler_protection              = var.spec.crawler_protection
  content_bots_protection         = var.spec.content_bots_protection
  cf_robots_variant               = var.spec.cf_robots_variant
  is_robots_txt_managed           = var.spec.is_robots_txt_managed
}
