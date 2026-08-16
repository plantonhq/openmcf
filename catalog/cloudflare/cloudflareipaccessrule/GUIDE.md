# CloudflareIpAccessRule guide

The judgment this guide protects you from: the scope is a one-of, and changing what a rule matches looks like an in-place update that silently does nothing. Pick the scope deliberately, and recreate the rule when the selector changes.

## Exactly one of account_id or zone_id

Account rules apply to every zone. Zone rules apply to one zone and can override an account rule for that zone. The provider accepts both at once and silently prefers the account -- this spec requires exactly one so the manifest states its intent. Setting both is rejected at validation; setting neither is too.

`zone_id` can be a literal or a `value_from` reference to a `CloudflareDnsZone` (defaults to `status.outputs.zone_id`).

## Configuration changes do not stick

Only `mode` and `notes` update in place. Cloudflare's API accepts a request that edits `configuration.target` or `configuration.value` and then ignores it -- the provider's own tests document this. A plan that shows an in-place configuration update will apply "successfully" and leave the old match serving.

To change what a rule matches: create a new rule with the new selector, then destroy the old one. Do not edit `target` or `value` on a live rule and expect the change to take.

## IPv6 is long form; CIDR prefixes are a short list

`ip6` values must be fully expanded -- eight colon-separated groups, e.g. `2001:0db8:0000:0000:0000:0000:0000:0001`. Compressed `::` notation is rejected here even though it is valid IPv6 everywhere else.

`ip_range` accepts only IPv4 `/16` or `/24`, and IPv6 `/32`, `/48`, or `/64`. A `/32` IPv4 or a `/20` will fail validation before it reaches the API. Use `ip` for a single IPv4 address and `ip6` for a single IPv6 address -- do not put a host in `ip_range`.

Country values are two characters (`US`). `T1` matches Tor exit nodes. ASN values look like `AS13335`.

## Destroy is a real delete

Destroy removes the rule. Matching traffic stops being allowed, blocked, or challenged by this object. There is no abandon-in-place behavior -- unlike Bot Management, the last-applied values do not linger.

## Pairs well with

- [CloudflareBotManagement](../cloudflarebotmanagement/README.md) -- zone-wide bot scoring; reach for this kind when the decision is a static IP/ASN/country selector.
- [CloudflareRuleset](../cloudflareruleset/README.md) -- expression-based WAF when a single selector is not enough.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from` for a zone-scoped rule.
