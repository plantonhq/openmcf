# CloudflareWebAnalyticsSite guide

The judgment this guide protects you from: this component's rules are write-only from the provider's point of view, and several site fields never come back from the API at all.

## Manage rules here or in the dashboard -- never both

The provider's read for `cloudflare_web_analytics_rule` is an empty function: it makes the API call for nothing and never compares what came back. So rules are effectively write-only. A rule someone adds or edits in the Cloudflare dashboard is invisible to IaC -- no drift shows in any plan -- and it stays invisible until the next apply re-asserts your declared rows and quietly overwrites the difference. Pick one owner for the rule set. If it is this manifest, then dashboard edits are temporary by definition.

## Four site fields also never come back

`enabled`, `host`, `lite`, and `zone_tag` are marked no-refresh at the provider: they are sent on write and never populated from a read. State carries whatever you configured, so drift on them is undetectable by design. On import the provider works around this for one field only -- it reaches into the computed ruleset to backfill `zone_tag` -- which is worth knowing when an imported site looks like it lost its other settings.

## The site token is public by nature, and still worth protecting

The beacon runs in visitors' browsers, so the token inevitably ships inside your pages -- it is not a secret in the way an API key is. These outputs still mark `site_token` and `snippet` sensitive, because that keeps them out of plan logs and CI output where they would sit next to real credentials and get treated as noise. Nothing breaks if the token is seen; the marking is hygiene, not a control.

## host and zone_tag are different products of the same feature

With `host` you measure any site and embed the snippet yourself. With `zone_tag` you measure a Cloudflare zone, and `auto_install` can put the beacon on every proxied page with no code change. `auto_install` needs the zone to actually be proxied (orange-clouded) -- on a DNS-only zone Cloudflare has no response to inject into, so the setting is accepted and does nothing.

## Deleting the site deletes the history

A real delete here is not just "stop measuring": the site's collected analytics stop being reachable. If the data matters, export what you need before retiring the site -- and note there is no provider-side deletion guard to catch a careless destroy.

## Rules are ordered, and order is positional

The module manages one provider object per declared row, keyed by position. Reordering rows therefore rewrites those objects rather than reordering them in place. That is harmless (the end state is what the list says) but it makes plans noisier than expected -- expect churn if you insert a row at the top of a long list.

## Pairs well with

- [CloudflareDnsZone](../cloudflarednszone/README.md) -- the zone a zone-measured site references.
- [CloudflareLogpushJob](../cloudflarelogpushjob/README.md) -- the server-side view of the same traffic, including non-browser requests.
- [CloudflareNotificationPolicy](../cloudflarenotificationpolicy/README.md) -- the `web_analytics_metrics_update` alert family.
