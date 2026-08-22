# CloudflareSnippetRules guide

The judgment this guide protects you from: this resource is the zone's entire snippet routing table, and destroy empties that table -- including rules you did not create. One manifest, one zone. A second manifest is a silent overwrite.

## One manifest owns the WHOLE table

Cloudflare's API is a full-replacement PUT. Every apply sends exactly the `rules` list in the manifest and the live table becomes that list. Rules added in the dashboard, by another team, or by a second Planton manifest disappear on the next apply.

Keep every snippet rule for a zone in this one resource. Split ownership and you will fight each other, one apply at a time, with no merge.

Rules evaluate in list order against Cloudflare's Rules language -- the same wirefilter expressions rulesets use.

## Destroy wipes everything, including outsiders

Destroying this resource deletes ALL snippet rules in the zone. Dashboard rules, API rules, rules from a previous manifest you forgot about -- gone. The snippets themselves survive; only the routing table empties.

If you need to stop invoking snippets without deleting the table's ownership, set each rule's `enabled: false` and apply. That is reversible. Destroy is not selective.

## enabled defaults TRUE here -- Cloudflare's default is FALSE

This is the footgun. Cloudflare's provider (and the API's older generation) defaults `enabled` to false. A rule that omits the field is created DISABLED and matches nothing. Migrated configurations have been burned by the flip.

This spec defaults `enabled` to true so a manifest that declares a rule gets a rule that runs. The module coalesces a null to true rather than passing it through -- that is what makes the promise real. Set `enabled: false` explicitly to stage a rule without activating it.

## snippet_name is a foreign key

`snippet_name` references a `CloudflareSnippet` (defaults to `status.outputs.snippet_name`). A literal name works too, and is how you invoke a snippet that already exists in the zone. A missing name fails at apply, not at validation -- Cloudflare is the wall.

## Pairs well with

- [CloudflareSnippet](../cloudflaresnippet/README.md) -- the code this table invokes; create the snippet first.
- [CloudflareRuleset](../cloudflareruleset/README.md) -- a different expression table (WAF, cache, redirects).
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`.
