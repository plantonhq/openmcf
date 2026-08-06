# GCP Health Check

Creates a Compute Engine health check — the probe load balancers consult before sending traffic to a backend and managed instance groups consult before auto-healing an instance. One kind covers both scopes: leave `region` empty for a global check (global external load balancers) or set it for a regional one (regional backend services and MIG auto-healing).

## What Gets Created

When you deploy a GcpHealthCheck resource, Planton provisions exactly one of:

- **Global health check** (`google_compute_health_check`) — when `region` is empty
- **Regional health check** (`google_compute_region_health_check`) — when `region` is set

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **IAM permissions** — any role carrying `compute.healthChecks.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpHealthCheck
metadata:
  name: web-backend-probe
spec:
  projectId:
    value: my-gcp-project-123
  http:
    portSpecification: USE_SERVING_PORT
    requestPath: /healthz
```

```shell
planton apply -f health-check.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `http` / `https` / `http2` / `tcp` / `ssl` / `grpc` / `grpcTls` | object | — | Required, exactly one. The probe protocol block with its port/path/response knobs. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the check. Can reference a GcpProject. Immutable. |
| `healthCheckName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable. |
| `region` | `string` | `""` (global) | Set for a regional health check. Immutable. |
| `checkIntervalSec` / `timeoutSec` | `int` | `5` / `5` | Probe cadence; timeout must not exceed the interval. |
| `healthyThreshold` / `unhealthyThreshold` | `int` | `2` / `2` | Consecutive verdicts to flip health state. |
| `enableLogging` | `bool` | `false` | Log every health state transition. |
| `sourceRegions` | `list(string)` | `[]` | Global-only: probe from exactly 3 named regions (interval ≥ 30 s, HTTP/HTTPS/TCP only). |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value backend services reference in their health check list |
| `health_check_name` | Name of the health check in GCP |
| `type` | Probe protocol GCP computed (HTTP, HTTPS, HTTP2, TCP, SSL, GRPC, GRPC_TLS) |
| `region` | Region of a regional check; empty for global |

## Related Components

- [GcpCloudArmorPolicy](/docs/catalog/gcp/gcpcloudarmorpolicy) — protects the backends this check keeps in rotation
- [GcpFirewallRule](/docs/catalog/gcp/gcpfirewallrule) — admits Google's prober ranges to instance-group backends
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the project that owns the check
