# Path-prefix route

One snippet rule that invokes `redirect_legacy` when the path starts with `/legacy`. This list is the zone's entire snippet routing table -- every apply replaces it, and destroy wipes every snippet rule in the zone.

## When to Use

- First (or only) snippet route on a zone
- Pair with the `CloudflareSnippet` **01-redirect** preset, which creates `redirect_legacy`
- A starting point to grow the table one rule at a time, in this same manifest

## Key Configuration Choices

- **One manifest, one zone** -- a second CloudflareSnippetRules against the same zone silently overwrites this table.
- **enabled is omitted and still true** -- this spec defaults it on. Cloudflare's own default is false; do not assume an omitted field means disabled.
- **snippet_name is a literal** -- `redirect_legacy`. Switch to `value_from` once the snippet is a Planton resource.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone whose snippet routing table is managed | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |
| `rules[0].snippet_name.value` | The snippet to invoke | A CloudflareSnippet's `status.outputs.snippet_name`, or a name that already exists in the zone |
| `rules[0].expression` | The match expression | Your path or header inventory |

## Related Presets

None on this kind -- pair with the `CloudflareSnippet` **01-redirect** preset for the code this rule invokes.
