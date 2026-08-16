# Cloudflare Zero Trust List

A reusable named set of values -- domains, IPs, URLs, emails, serial numbers, and kin -- that Gateway policies and device-posture rules reference by ID. The list type is immutable: changing it replaces the list and breaks every policy that referenced the old ID.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Zero Trust list** -- one `cloudflare_zero_trust_list` on the account, with its items as a set (order is not preserved)

## Prerequisites

- **A Cloudflare account with Zero Trust enabled**
- **A Cloudflare API token** with Account → Zero Trust → Edit
- **A Gateway policy (or posture rule)** that will reference `list_id` -- the list does nothing on its own

## Quick Start

A domain list two Gateway policies can share:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustList
metadata:
  name: blocked-domains
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: blocked-domains
  type: DOMAIN
  items:
    - value: gambling.example.com
    - value: casino.example.net
```

```shell
planton apply -f zt-list.yaml
```

Reference `status.outputs.list_id` from a Gateway policy `traffic` expression. Do not use `CloudflareList` here -- that is the older Ruleset-family list.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | Owning account. | Required. 32-character hex. |
| `name` | string | Display name. | Required, min length 1. |
| `type` | string | What kind of values the items carry. Immutable at Cloudflare. | Required. One of `SERIAL`, `URL`, `DOMAIN`, `EMAIL`, `IP`, `CATEGORY`, `LOCATION`, `DEVICE`, `AAGUID`. Use the uppercase form. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `description` | string | unset | Purpose. |
| `items` | object[] | empty | Set of `{value, description}`. `value` is required. Order is not significant and is not preserved. URL-type values are normalized by the API (a known plan-drift source). |

## Examples

### Domain blocklist

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustList
metadata:
  name: blocked-domains
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: blocked-domains
  type: DOMAIN
  description: Domains Gateway policies block for all users
  items:
    - value: gambling.example.com
      description: policy-blocked domain
    - value: casino.example.net
```

### IP allowlist

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustList
metadata:
  name: corp-cidrs
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: corp-cidrs
  type: IP
  items:
    - value: 203.0.113.0/24
    - value: 198.51.100.10
```

### Email allowlist

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustList
metadata:
  name: contractor-emails
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: contractor-emails
  type: EMAIL
  items:
    - value: ada@contractor.example.com
    - value: linus@contractor.example.com
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `list_id` | string | UUID referenced by Gateway policies and posture rules |

## Related Components

- [Cloudflare Zero Trust Gateway Policy](/docs/catalog/cloudflare/cloudflarezerotrustgatewaypolicy) -- references this list from `traffic` / `identity`
- [Cloudflare List](/docs/catalog/cloudflare/cloudflarelist) -- the older Ruleset-family list; do not mix the two
