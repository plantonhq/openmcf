# A Web Analytics (RUM) site plus its folded measurement rules. The site is
# identified by host OR zone (spec validation enforces exactly one).
#
# Cloudflare stores include/exclude rules as separate objects under the
# site's ruleset; this module manages one rule object per declared row,
# keyed by position. The provider never reads rules back after writing them
# (its refresh is deliberately blind), so each apply re-asserts exactly the
# declared rows.
resource "cloudflare_web_analytics_site" "main" {
  account_id = var.spec.account_id

  host     = try(var.spec.host, "") != "" ? var.spec.host : null
  zone_tag = try(var.spec.zone_tag, "") != "" ? var.spec.zone_tag : null

  auto_install = try(var.spec.auto_install, null)
  enabled      = try(var.spec.enabled, null)
  lite         = try(var.spec.lite, null)
}

resource "cloudflare_web_analytics_rule" "main" {
  count = length(try(var.spec.rules, []))

  account_id = var.spec.account_id
  ruleset_id = cloudflare_web_analytics_site.main.ruleset.id

  host      = try(var.spec.rules[count.index].host, "") != "" ? var.spec.rules[count.index].host : null
  paths     = length(try(var.spec.rules[count.index].paths, [])) > 0 ? var.spec.rules[count.index].paths : null
  inclusive = try(var.spec.rules[count.index].inclusive, null)
  is_paused = try(var.spec.rules[count.index].is_paused, null)
}
