# Cloudflare Snippet Rules

The zone's snippet routing table: an ordered list of expressions deciding which requests invoke which snippet. One manifest owns the WHOLE table -- every apply is a full-replacement PUT. Destroy wipes every snippet rule in the zone, including ones created outside this manifest.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Snippet rules** -- one `cloudflare_snippet_rules` on the zone, whose `rules` list is the entire routing table

## Prerequisites

- **A Cloudflare zone** -- typically a CloudflareDnsZone whose `zone_id` output this resource references
- **A Cloudflare API token** with Zone → Snippets → Edit
- **At least one CloudflareSnippet** -- rules invoke snippets by name; a missing name fails at apply
- **Sole ownership of the zone's snippet-rule table** -- a second manifest against the same zone silently overwrites this one

## Quick Start

One rule that invokes a snippet on a path prefix:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareSnippetRules
metadata:
  name: legacy-routes
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  rules:
    - expression: 'http.request.uri.path starts_with "/legacy"'
      snippetName:
        value: redirect_legacy
      description: Redirect legacy paths
```

```shell
planton apply -f snippet-rules.yaml
```

`enabled` is omitted and still true -- this spec defaults it on so a declared rule runs. Cloudflare's own default is false.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone whose snippet routing table is managed. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. |
| `rules` | object[] | The WHOLE table, evaluated in order. | At least one. Each rule needs `expression` and `snippetName`. |

### Rule Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `expression` | string | -- | Cloudflare Rules language (wirefilter). Validated by Cloudflare at apply. |
| `snippetName` | StringValueOrRef | -- | The snippet to invoke. Can reference a CloudflareSnippet via `valueFrom` (defaults to `status.outputs.snippet_name`). |
| `description` | string | unset | Shown in the dashboard. |
| `enabled` | bool | **true** | FOOTGUN: Cloudflare defaults this to false. This spec defaults it to true so a declared rule runs. Set `false` explicitly to stage a rule. |

## Examples

### Path-prefix route

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareSnippetRules
metadata:
  name: legacy-routes
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
  rules:
    - expression: 'http.request.uri.path starts_with "/legacy"'
      snippetName:
        valueFrom:
          kind: CloudflareSnippet
          name: redirect-legacy
          fieldPath: status.outputs.snippet_name
      description: Redirect legacy paths
```

## Destroy Semantics

Destroy wipes ALL snippet rules in the zone, including ones created outside this manifest (dashboard, API). The snippets themselves survive. Every apply also replaces the whole table -- rules you did not declare disappear.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `zone_id` | string | The zone whose snippet routing table is managed. The table is a zone singleton -- the zone ID is its identity. |

## Related Components

- [Cloudflare Snippet](/docs/catalog/cloudflare/cloudflaresnippet) -- the code this table invokes by `snippetName`
- [Cloudflare Ruleset](/docs/catalog/cloudflare/cloudflareruleset) -- expression-based WAF and cache rules, a different table
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key
