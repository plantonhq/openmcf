# Cloudflare IP Access Rule

One allow, block, or challenge decision applied to traffic matching an IP, IP range, ASN, or country -- account-wide or on a single zone. Exactly one scope must be set. Only `mode` and `notes` update in place; changing what a rule matches does not stick at the API.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **IP Access rule** -- one `cloudflare_access_rule` on the account or the zone, with a single `{target, value}` configuration

## Prerequisites

- **A Cloudflare account** (account-wide rule) **or a Cloudflare zone** (zone-scoped rule)
- **A Cloudflare API token** with Account → Firewall Access Rules → Edit (account scope) or the zone equivalent
- **A selector you are willing to recreate** -- changing target or value later requires delete-and-recreate, not an in-place edit

## Quick Start

Block a documentation-range IPv4 address on every zone in the account:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareIpAccessRule
metadata:
  name: block-scanner
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  mode: block
  configuration:
    target: ip
    value: "192.0.2.10"
  notes: Block a known scanner
```

```shell
planton apply -f ip-access-rule.yaml
```

Do not set `zoneId` on the same manifest -- the spec requires exactly one scope. The provider would silently prefer the account if both arrived.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` XOR `zoneId` | string / StringValueOrRef | Exactly one. `zoneId` can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Exactly one required. `accountId` is a 32-character hex string when set. |
| `mode` | string | Action on matching traffic. | Required. One of `block`, `challenge`, `whitelist`, `js_challenge`, `managed_challenge`. |
| `configuration` | object | `{target, value}` -- what the rule matches. | Required. See targets below. Changing this on an existing rule does not stick. |

### Configuration targets

| Target | Value shape | Notes |
|--------|-------------|-------|
| `ip` | IPv4 address (`192.0.2.10`) | Use `ip6` for IPv6, `ip_range` for CIDR. |
| `ip6` | Fully-expanded IPv6 (`2001:0db8:0000:0000:0000:0000:0000:0001`) | Cloudflare rejects compressed `::` notation here. |
| `ip_range` | IPv4 `/16` or `/24`; IPv6 `/32`, `/48`, or `/64` | Those are the only prefix lengths Cloudflare accepts. |
| `asn` | `AS13335` | Autonomous system. |
| `country` | Two-character code (`US`; `T1` for Tor) | ISO 3166-1 alpha-2. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `notes` | string | unset | Why the rule exists. Updates in place with `mode`. |

## Examples

### Account-wide IP block

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareIpAccessRule
metadata:
  name: block-scanner
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  mode: block
  configuration:
    target: ip
    value: "192.0.2.10"
  notes: Block a known scanner
```

### Zone-scoped country challenge

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareIpAccessRule
metadata:
  name: challenge-us
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  mode: managed_challenge
  configuration:
    target: country
    value: US
  notes: Challenge visitors from the United States
```

## Destroy Semantics

Destroy is a real delete. The rule stops matching immediately. A configuration (target/value) change is not a replace -- it plans an in-place update that Cloudflare ignores, so the old match keeps serving. Recreate the rule to change what it matches.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `rule_id` | string | The created rule's ID |
| `zone_id` | string | The zone the rule was created on (empty for account-wide rules) |
| `account_id` | string | The account the rule was created on (empty for zone-scoped rules) |

## Related Components

- [Cloudflare Bot Management](/docs/catalog/cloudflare/cloudflarebotmanagement) -- zone-wide bot scoring; this kind is a static selector
- [Cloudflare Ruleset](/docs/catalog/cloudflare/cloudflareruleset) -- expression-based WAF and custom rules
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key for a zone-scoped rule
