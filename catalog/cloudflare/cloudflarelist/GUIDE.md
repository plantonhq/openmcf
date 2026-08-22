# CloudflareList guide

Operational judgment for Cloudflare Lists. The README covers what each field is; this covers how the pieces interact.

## The list is the container; items are a different kind

Never put entries in this kind. Inline `items` on `cloudflare_list` and CloudflareListItem are competing writers — the provider will fight itself and you will lose entries. Create the list empty, then add CloudflareListItem resources (one per entry, independent lifecycle).

## Names are expression identifiers, not labels

The list name appears in rule expressions as `$name`. That is why hyphens are forbidden (`^[a-zA-Z][a-zA-Z0-9_]*$`) and why the name is immutable. Pick a short, stable, lowercase identifier before the first apply; renaming replaces the list and every item has to be rewritten.

## Kind is a create-time contract

`ip`, `redirect`, `hostname`, and `asn` accept different item shapes. Changing kind replaces the list. A redirect list is what Bulk Redirect rulesets consume via `from_list`; an ip list is what WAF/custom rules consume via `ip.src in $name`. Pick the kind from the consumer, not from the data.

## Empty lists are valid

A list with zero items is a real, referenceable object. That is the right shape for a Bulk Redirect container you will fill later, and for an ip list you want to exist before the first block rule is written.
