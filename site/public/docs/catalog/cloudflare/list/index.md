---
title: "List"
description: "List deployment documentation"
icon: "package"
order: 100
componentName: "cloudflarelist"
---

# List on Cloudflare

Deploys a reusable, named Cloudflare List -- a collection you reference from rule expressions such as a WAF rule's `ip.src in $office_ips` or a Bulk Redirect ruleset's `from_list`. A list fixes a single entry shape at creation (IP/CIDR, redirect, hostname, or ASN); its entries are managed independently as List Item Cloud Resources, so one list can hold a handful of curated values or a large, separately-owned set. Lists are account-scoped and integrate with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **List** -- an account-scoped, named collection of the chosen kind (ip, redirect, hostname, or asn), empty until entries are added
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Account Filter Lists edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **Account-level access** -- lists are created at the account level (not per-zone), so the API token must be scoped to the account.

## Deploy

### Console

Open the deployment store, find **List on Cloudflare**, and click **Deploy**. The creation wizard captures the owning account, the list type, an expression-friendly name, and an optional description. Start from the **IP allowlist** preset in the [Presets](#presets) tab for a typical firewall list.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareList
metadata:
  name: office-ips
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  kind: ip
  name: office_ips
  description: Approved office egress IPs
```

```shell
planton apply -f cloudflare-list.yaml
```

This creates an empty IP list named `office_ips`, ready for entries and referenceable from rules as `$office_ips`. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a list. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Kind (`kind`)** -- Fixes which entry shape the list accepts: `ip` (addresses/CIDRs), `redirect` (Bulk Redirect rules), `hostname` (optionally wildcarded), or `asn` (Autonomous System Numbers). Immutable -- changing it replaces the list.

**Name (`name`)** -- The identifier used in filter and rule expressions, written as `$name`. Prefer short, lowercase, expression-friendly names. Immutable -- because rules reference it literally, renaming would break them, so changing it replaces the list.

**Description (`description`)** -- An optional human-readable summary shown in the Cloudflare dashboard.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign-key dependencies -- a list is defined entirely by its own account, kind, and name.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `list_id` | The Cloudflare-assigned identifier of the list | Referenced by a CloudflareListItem's `listId` field |
| `name` | The list name (echoed; the identifier used in rule expressions) | Referenced by a CloudflareRuleset Bulk Redirect rule's `from_list` |
| `kind` | The list kind (ip, redirect, hostname, or asn) | Verification, dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**IP allowlist** -- An `ip` list referenced from a WAF or firewall rule to allow (or block) a managed set of addresses. Use for office egress IPs, partner ranges, or known-bad CIDRs. Start from the **IP allowlist** preset.

**Bulk redirect list** -- A `redirect` list a Bulk Redirect ruleset points at via `from_list`, to migrate large numbers of legacy URLs in one place. Start from the **Bulk redirect** preset.

## Works With

- [**List Item on Cloudflare**](/cloud-catalog/cloudflare-list-item) -- adds entries to this list (via `listId`), one independently-managed entry at a time
- [**Ruleset on Cloudflare**](/cloud-catalog/cloudflare-ruleset) -- a Bulk Redirect rule references a redirect list by name (via `from_list`)
