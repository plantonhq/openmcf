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

  # auto_install is ALWAYS sent: Cloudflare echoes a server-side false for
  # it even when never sent (measured on host AND zone sites), so an
  # omitted send drifts forever (config null vs state false). Unset in the
  # spec means false -- exactly Cloudflare's default, so semantics are
  # unchanged. enabled/lite stay conditional: they are no_refresh (never
  # read back), so no echo can drift them.
  auto_install = try(var.spec.auto_install, null) != null ? var.spec.auto_install : false
  enabled      = try(var.spec.enabled, null)
  lite         = try(var.spec.lite, null)
}

# Read-after-create (measured 2026-08-27): the create response omits
# `snippet` (GET-only), and the resource's computed attributes sit null in
# state until the first refresh -- so the snippet/ruleset_id outputs and
# the folded rules ride this data source, not the resource. The `ruleset`
# object itself is IDENTITY-DEPENDENT, not read-dependent: zone-linked
# sites carry it in every response (create included); host-identified
# sites have NO ruleset, ever (which is why rules require zone_tag -- the
# spec walls it). NOTE: reading site_info requires the Account Settings
# READ permission -- Account Settings Write alone creates sites but cannot
# read them back (measured 403/10000).
data "cloudflare_web_analytics_site" "main" {
  account_id = var.spec.account_id
  site_id    = cloudflare_web_analytics_site.main.id
}

# Every rule field is ALWAYS sent (measured 2026-08-27): Cloudflare's rule
# form validates each field's presence and rejects omissions one by one
# (400 code 10001 "form.host.invalid" / "form.is_paused.invalid" -- the
# provider passes nulls straight through because upstream never exercises
# rules). An empty spec host means "every host", which the API spells "*".
resource "cloudflare_web_analytics_rule" "main" {
  count = length(try(var.spec.rules, []))

  account_id = var.spec.account_id
  ruleset_id = data.cloudflare_web_analytics_site.main.ruleset.id

  host      = try(var.spec.rules[count.index].host, "") != "" ? var.spec.rules[count.index].host : "*"
  paths     = length(try(var.spec.rules[count.index].paths, [])) > 0 ? var.spec.rules[count.index].paths : null
  inclusive = try(var.spec.rules[count.index].inclusive, null) != null ? var.spec.rules[count.index].inclusive : false
  is_paused = try(var.spec.rules[count.index].is_paused, null) != null ? var.spec.rules[count.index].is_paused : false
}
