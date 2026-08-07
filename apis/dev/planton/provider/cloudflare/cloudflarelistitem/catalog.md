# List Item on Cloudflare

Adds a single entry to a Cloudflare List. Items have independent lifecycles, so a list can grow or shrink one entry at a time without rewriting the whole set -- ideal when entries are owned by different teams or pipelines. The entry's shape (an IP/CIDR, an ASN, a hostname, or a redirect rule) must match the parent list's kind. List Items integrate with Planton's Provider Connections for Cloudflare credential management and reference their parent list through Planton's foreign-key system.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **List Item** -- one entry written into the referenced list, of the shape that matches the list's kind
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Account Filter Lists edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A parent List** -- a CloudflareList to write into. Reference it (recommended) so Planton orders creation and renders the dependency, or supply a literal list ID for a list managed elsewhere.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **Matching list kind** -- the entry value must match the parent list's kind. Cloudflare rejects a mismatched entry (e.g. an IP value in a hostname list) at apply time.

## Deploy

### Console

Open the deployment store, find **List Item on Cloudflare**, and click **Deploy**. The creation wizard captures the owning account and parent list (with a live connection diagram), then the entry value via a type-aware editor. Start from the **IP entry** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareListItem
metadata:
  name: office-ip-1
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  listId:
    valueFrom:
      kind: CloudflareList
      name: office-ips
      fieldPath: status.outputs.list_id
  ip: "203.0.113.0/24"
  comment: Approved office egress range
```

```shell
planton apply -f cloudflare-list-item.yaml
```

This adds the `203.0.113.0/24` CIDR to the `office-ips` list, referencing the list's output so the dependency is wired automatically. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a list item. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Parent list (`listId`)** -- The list this entry is written to. Reference a CloudflareList (resolves to its `status.outputs.list_id`) so Planton orders creation and the graph shows the edge, or paste a literal list ID.

**Entry value (the `item` oneof)** -- Exactly one of: `ip` (an IPv4/IPv6 address or CIDR), `asn` (an Autonomous System Number), `hostname` (a hostname, optionally wildcarded with `excludeExactHostname`), or `redirect` (a Bulk Redirect rule with source/target URLs, a status code, and behavior flags). The chosen shape must match the list's kind.

**Comment (`comment`)** -- An optional note about why the entry exists (up to 500 characters).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | Description |
|------------|-------|-------------|
| CloudflareList | `listId` | The parent list this entry is written to (literal ID or a reference to a list's `list_id` output) |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `item_id` | The Cloudflare-assigned identifier of the list item | Verification, dashboards |
| `list_id` | The list ID the entry was written to | Verification, traceability |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**IP entry** -- A single IP/CIDR added to an `ip` list, referenced from firewall rules. Use to extend an allowlist or blocklist one address at a time. Start from the **IP entry** preset.

**Redirect entry** -- A source-to-target redirect added to a `redirect` list backing a Bulk Redirect ruleset. Use to migrate individual legacy URLs. Start from the **Redirect entry** preset.

## Works With

- [**List on Cloudflare**](/cloud-catalog/cloudflare-list) -- the parent container this entry is written into (via `listId`)
