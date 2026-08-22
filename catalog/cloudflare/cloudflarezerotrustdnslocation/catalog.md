# Cloudflare Zero Trust DNS Location

A Gateway DNS location: a named entry point whose DNS traffic Gateway filters. Cloudflare assigns the resolver endpoints (DoH subdomain, destination IPs); policies match on the location.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **DNS location** -- one `cloudflare_zero_trust_dns_location`

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step)
- **A Cloudflare API token** with Account → Zero Trust → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustDnsLocation
metadata:
  name: hq-office
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: hq-office
```

```shell
planton apply -f dns-location.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `name` | string | The location's name. | Required. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `endpoints` | object | The resolver endpoint tree. | All four types (doh/dot/ipv4/ipv6) declared when set. |
| `networks` | list | Source IPv4 CIDRs for the IPv4 endpoint. | Cloudflare caps IPv4 CIDRs at /24. |
| `maxTtl` | object | TTL capping. | mode inherit/override/disabled; `ttlSecs` (60-36000) exactly with override. |
| `dnsDestinationIpsId` | string | Destination-IP mapping UUID. | Leave unset for the shared-pool auto-assign. |
| `clientDefault` | bool | The account's default location. | A real account-level lever. |
| `ecsSupport` | bool | Honor EDNS Client Subnet. | |

## Destroy Semantics

Real delete: the location and its endpoints disappear immediately.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `location_id` | string | The Cloudflare-assigned UUID |
| `doh_subdomain` | string | The unique DoH subdomain (https://<sub>.cloudflare-gateway.com/dns-query) |
| `ip` | string | The IPv4 destination for the plain-DNS endpoint |

## Related Components

- [Cloudflare Zero Trust Gateway Policy](/docs/catalog/cloudflare/cloudflarezerotrustgatewaypolicy) -- the filtering rules
- [Cloudflare Zero Trust Gateway Settings](/docs/catalog/cloudflare/cloudflarezerotrustgatewaysettings) -- the account posture
