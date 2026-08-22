# Redirect snippet

A single-file snippet that returns a 302. The snippet does nothing until a `CloudflareSnippetRules` rule invokes `redirect_legacy` by name. Create is an upsert -- deploying this name on a zone that already has it silently overwrites the existing snippet.

## When to Use

- A legacy-path redirect too small for a Worker
- First snippet on a zone
- A shared snippet two or more rules can invoke

## Key Configuration Choices

- **snippet_name: redirect_legacy** -- letters, numbers, and underscores only. Hyphens are rejected. This name is the identity.
- **main_module: main.js** -- must name one of `files`. The provider argument is `metadata`.
- **Byte-stable content** -- keep the one-liner as-is. Formatting churn reads back as drift because the provider rebuilds stored source from a multipart read.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone the snippet is deployed to | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |
| `files[0].content` | The redirect target URL inside `Response.redirect` | Your destination hostname |

## Related Presets

None on this kind -- pair with the `CloudflareSnippetRules` **01-path-prefix** preset to invoke this snippet on a path.
