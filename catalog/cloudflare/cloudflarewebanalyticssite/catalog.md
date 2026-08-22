# Cloudflare Web Analytics Site

A Web Analytics (RUM) site: privacy-first real-user monitoring for one website, with include/exclude measurement rules. Real create, update, delete.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Web Analytics site** -- one `cloudflare_web_analytics_site`
- **Measurement rules** -- one `cloudflare_web_analytics_rule` per declared `rules[]` row

## Prerequisites

- **A Cloudflare API token** with Account Settings → Write
- For `autoInstall`: a zone proxied through Cloudflare

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWebAnalyticsSite
metadata:
  name: www-example-rum
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  host: www.example.com
```

```shell
planton apply -f web-analytics-site.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `host` / `zoneTag` | string | The site to measure. | Exactly one must be set. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `autoInstall` | bool | Inject the beacon at the edge. | Proxied zones only (pair with `zoneTag`). |
| `enabled` | bool | Whether measurement is active. | Cloudflare's default is enabled. |
| `lite` | bool | Lightweight beacon variant. | Reduced metric set. |
| `rules` | list | Include/exclude rows. | `host`, `paths[]`, `inclusive`, `isPaused`; managed as one provider object per row. |

## Destroy Semantics

Real delete for the site and every rule. Measurement stops; the site's historical analytics are no longer reachable.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `siteTag` | string | The Cloudflare-assigned site tag |
| `siteToken` | string | The beacon's measurement token (secret-marked) |
| `snippet` | string | The ready-to-embed script tag (secret-marked) |
| `rulesetId` | string | The parent object the rules live under |

## Related Components

- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- the zone a zone-measured site points at
- [Cloudflare Logpush Job](/docs/catalog/cloudflare/cloudflarelogpushjob) -- edge-side request logs for the same traffic
- [Cloudflare Notification Policy](/docs/catalog/cloudflare/cloudflarenotificationpolicy) -- the web-analytics metrics alert family
