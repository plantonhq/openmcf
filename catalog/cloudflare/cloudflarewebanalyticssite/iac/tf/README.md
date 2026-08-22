# CloudflareWebAnalyticsSite Terraform Module

Terraform IaC module for Web Analytics (RUM) sites.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareWebAnalyticsSiteSpec (generated)
locals.tf     — Naming/labels
main.tf       — cloudflare_web_analytics_site + one cloudflare_web_analytics_rule per rules[] row
outputs.tf    — site_tag, site_token, snippet, ruleset_id
```

## Behavior

A plain CRUD resource: real create/update/delete for the site and every rule. The site is identified by `host` OR `zone_tag` (spec validation enforces exactly one).

Rules fan out one provider object per declared row, keyed by position, with `ruleset_id` wired from the created site. The provider never reads rules back (its Read is an empty stub), so each apply re-asserts exactly the declared rows -- dashboard edits to the rule set are invisible and will be overwritten. Four site attributes (`enabled`, `host`, `lite`, `zone_tag`) are also never refreshed by the provider.

`site_token` and `snippet` are marked sensitive in the outputs: they carry the beacon credential, so keeping them out of plan logs is hygiene.

Import the site as `{account_id}/{site_tag}` (the site tag, not the id copy); rules ship no importer upstream.

## Outputs

| Name | Description |
|------|-------------|
| `site_tag` | The Cloudflare-assigned site tag |
| `site_token` | The beacon's measurement token (sensitive) |
| `snippet` | The ready-to-embed script tag (sensitive) |
| `ruleset_id` | The parent object the measurement rules live under |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
