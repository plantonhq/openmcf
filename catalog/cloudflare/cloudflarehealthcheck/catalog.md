# Cloudflare Health Check

A standalone origin probe with healthy/unhealthy status. No load balancer required. Health checks are a paid zone feature (Pro+); the API enforces the plan gate at create. `type` is `HTTP`, `HTTPS`, or `TCP` -- `httpConfig` only for HTTP/HTTPS, `tcpConfig` only for TCP.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Health check** -- one `cloudflare_healthcheck` on the zone, with either `http_config` or `tcp_config` depending on `type`

## Prerequisites

- **A Cloudflare zone on a Pro plan or above** -- free zones are rejected at the API
- **A Cloudflare API token** with Zone → Health Checks → Edit. Prefer the API token over the global API key; some upstream provider tests blank `CLOUDFLARE_API_TOKEN` and that class of auth defect is real
- **An origin that can be reached from Cloudflare's probe regions** -- the check is created even if the origin never answers; status will sit unhealthy

## Quick Start

An HTTP probe against `example.com/health`:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareHealthcheck
metadata:
  name: origin-http
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: origin-http
  address: example.com
  type: HTTP
  httpConfig:
    path: /health
    expectedCodes:
      - "200"
```

```shell
planton apply -f healthcheck.yaml
```

Do not attach this to a load balancer -- that is `CloudflareLoadBalancerMonitor`. This kind watches an origin.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone the health check belongs to. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. |
| `name` | string | Short name shown in the dashboard and alerts. | Required, min length 1. |
| `address` | string | Origin hostname or IP. | Required. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `type` | string | `HTTP` | `HTTP`, `HTTPS`, or `TCP`. A deliberate tightening -- the provider accepts any string. |
| `checkRegions` | string[] | Cloudflare default region | One of `WNAM`, `ENAM`, `WEU`, `EEU`, `NSAM`, `SSAM`, `OC`, `ME`, `NAF`, `SAF`, `IN`, `SEAS`, `NEAS`, `ALL_REGIONS` (Enterprise). |
| `consecutiveFails` | int32 | unset (API: 1) | Failed probes before unhealthy. |
| `consecutiveSuccesses` | int32 | unset (API: 1) | Successful probes before healthy again. |
| `interval` | int32 | unset (API: 60) | Seconds between probes. |
| `retries` | int32 | unset (API: 2) | Immediate retries on timeout. |
| `timeout` | int32 | unset (API: 5) | Probe timeout in seconds. |
| `suspended` | bool | unset | Pause probing without deleting the check. |
| `httpConfig` | object | unset | HTTP/HTTPS only. `method` (GET/HEAD), `path`, `port`, `expectedCodes`, `expectedBody`, `followRedirects`, `allowInsecure`, `headers` (map of name → `{values}`). The provider argument is `header`. |
| `tcpConfig` | object | unset | TCP only. `method` (`connection_established`), `port`. |

`httpConfig` on a TCP check (or `tcpConfig` on HTTP/HTTPS) is rejected at validation.

## Examples

### HTTP origin probe

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareHealthcheck
metadata:
  name: origin-http
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: origin-http
  address: example.com
  type: HTTP
  httpConfig:
    path: /health
    expectedCodes:
      - "200"
```

### TCP port probe

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareHealthcheck
metadata:
  name: origin-tcp
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
  name: origin-tcp
  address: origin.example.com
  type: TCP
  tcpConfig:
    port: 5432
```

## Destroy Semantics

Destroy is a real delete. The probe stops and its history is removed. Set `suspended: true` to pause without losing the ID.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `healthcheck_id` | string | The created health check's ID |
| `zone_id` | string | The zone the health check belongs to |

## Related Components

- [Cloudflare Load Balancer Monitor](/docs/catalog/cloudflare/cloudflareloadbalancermonitor) -- the account-scoped monitor a pool consumes; do not mix the two
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key; the zone plan gates create
