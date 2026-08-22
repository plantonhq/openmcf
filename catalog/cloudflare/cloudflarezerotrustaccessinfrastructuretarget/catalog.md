# Cloudflare Zero Trust Access Infrastructure Target

An infrastructure target: a server (hostname + private IP) that Access infrastructure applications broker SSH access to. Targets are the inventory; applications select them by hostname or IP.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Infrastructure target** -- one `cloudflare_zero_trust_access_infrastructure_target`

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step)
- **A Cloudflare API token** with Account → Zero Trust → Edit
- **A tunnel path to the server** (a `CloudflareZeroTrustTunnel` routing its network) for live SSH brokering

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustAccessInfrastructureTarget
metadata:
  name: prod-db-1
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  hostname: prod-db-1
  ip:
    ipv4:
      ipAddr: 10.0.10.5
```

```shell
planton apply -f target.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `hostname` | string | The target's hostname. | ≤255 chars; letters/digits/dashes/periods; alphanumeric ends. |
| `ip` | object | The target's addressing. | At least one family declared. |

### The ip block

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `ipv4.ipAddr` / `ipv6.ipAddr` | string | The address for that family. | Required within a declared family. |
| `ipv4.virtualNetworkId` / `ipv6.virtualNetworkId` | string/ref | The virtual network the address lives in. | Omit for the account's default virtual network. |

## Destroy Semantics

Real delete: the target leaves the inventory immediately.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `target_id` | string | The Cloudflare-assigned UUID of the target |

## Related Components

- [Cloudflare Zero Trust Tunnel Virtual Network](/docs/catalog/cloudflare/cloudflarezerotrusttunnelvirtualnetwork) -- network segments for overlapping CIDRs
- [Cloudflare Zero Trust Access Application](/docs/catalog/cloudflare/cloudflarezerotrustaccessapplication) -- selects targets by hostname/IP
- [Cloudflare Zero Trust Tunnel](/docs/catalog/cloudflare/cloudflarezerotrusttunnel) -- the data path
