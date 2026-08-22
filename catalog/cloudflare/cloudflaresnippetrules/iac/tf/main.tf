# Cloudflare snippet rules: the zone's snippet routing table. A zone singleton
# -- Cloudflare's API replaces the WHOLE list on every apply, so this manifest
# is the entire table (a second manifest against the same zone would silently
# overwrite this one). Destroy deletes ALL snippet rules in the zone, including
# rules created outside this manifest; the snippets themselves survive.
#
# enabled defaults to TRUE here even though Cloudflare's own default is FALSE:
# a declared rule should run. The explicit coalesce (rather than passing null
# through) is what makes the spec's promised default real on this engine --
# a null would let the provider default the rule to disabled.
resource "cloudflare_snippet_rules" "main" {
  zone_id = var.spec.zone_id

  rules = [
    for rule in var.spec.rules : {
      expression   = rule.expression
      snippet_name = rule.snippet_name
      description  = rule.description != "" ? rule.description : null
      enabled      = rule.enabled != null ? rule.enabled : true
    }
  ]
}
