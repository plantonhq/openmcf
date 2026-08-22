# Cloudflare Snippet

A small JavaScript module at the zone's edge, invoked by snippet rules. The snippet name is the identity: create is an upsert, so a name collision silently overwrites the existing snippet. Code travels inline; `mainModule` must name one of the files.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Snippet** -- one `cloudflare_snippet` on the zone, with inline files and a `metadata.main_module` entry point

## Prerequisites

- **A Cloudflare zone** -- typically a CloudflareDnsZone whose `zone_id` output this resource references
- **A Cloudflare API token** with Zone → Snippets → Edit
- **A unique `snippetName`** -- deploying a name that already exists in the zone silently overwrites it
- **Headroom under the plan's snippet-count limit** -- free plans allow a small number; the API enforces the cap at create

## Quick Start

A 302 redirect snippet two rules can share:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareSnippet
metadata:
  name: redirect-legacy
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  snippetName: redirect_legacy
  files:
    - name: main.js
      content: "export default { async fetch(request) { return Response.redirect(\"https://www.example.com\", 302); } };"
  mainModule: main.js
```

```shell
planton apply -f snippet.yaml
```

Reference `status.outputs.snippet_name` from a `CloudflareSnippetRules` rule. The snippet does nothing until a rule invokes it.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone the snippet is deployed to. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. |
| `snippetName` | string | Identity within the zone. Changing it replaces the snippet. | Required. Letters, numbers, and underscores only. |
| `files` | object[] | Source files. Most snippets are a single file. | At least one. Each entry needs `name` and `content`. `mainModule` must name one of them. |
| `mainModule` | string | The entry file's name. The provider argument is `metadata`. | Required. Must equal one of `files[].name`. |

### Optional Fields

None. Content must be byte-stable -- Cloudflare rebuilds stored source from a multipart read, and gratuitous formatting differences read back as drift.

## Examples

### Legacy-path redirect

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareSnippet
metadata:
  name: redirect-legacy
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  snippetName: redirect_legacy
  files:
    - name: main.js
      content: "export default { async fetch(request) { return Response.redirect(\"https://www.example.com\", 302); } };"
  mainModule: main.js
```

## Destroy Semantics

Destroy is a real delete. The snippet is removed from the zone. Rules that referenced `snippetName` keep pointing at that name and start invoking nothing. Create is an upsert -- a second manifest with the same `snippetName` on the same zone silently overwrites the first.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `snippet_name` | string | The snippet's name -- its identity within the zone, and what snippet rules reference |
| `zone_id` | string | The zone the snippet is deployed to |

## Related Components

- [Cloudflare Snippet Rules](/docs/catalog/cloudflare/cloudflaresnippetrules) -- the zone's routing table that invokes this snippet by name
- [Cloudflare Worker](/docs/catalog/cloudflare/cloudflareworker) -- the full Worker when you need bindings, cron, or a custom domain
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key
