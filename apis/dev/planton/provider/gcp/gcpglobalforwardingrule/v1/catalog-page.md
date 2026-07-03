# GCP Global Forwarding Rule

Creates a global Compute Engine forwarding rule — the VIP node where traffic enters a global load balancer. It binds an IP address and port to a target proxy, and doubles as the Private Service Connect entry point (private paths to Google APIs or producer service attachments).

## What Gets Created

A single `google_compute_global_forwarding_rule` in the chosen GCP project.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **A target proxy** — a `GcpTargetHttpsProxy` or `GcpTargetHttpProxy`
- **Recommended: a reserved static IP** — a `GcpGlobalAddress`
- **IAM permissions** — any role carrying `compute.globalForwardingRules.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpGlobalForwardingRule
metadata:
  name: web-frontend-443
spec:
  projectId:
    value: my-gcp-project-123
  target:
    value: https://www.googleapis.com/compute/v1/projects/my-gcp-project-123/global/targetHttpsProxies/web-https-frontend
  portRange: "443"
```

```shell
planton apply -f forwarding-rule.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `target` | `StringValueOrRef` | — | Required. Proxy self-link, PSC bundle (`all-apis`/`vpc-sc`), or service attachment URI. Mutable in place |
| `ipAddress` | `StringValueOrRef` | ephemeral | `GcpGlobalAddress` ref, literal IP, or empty. Immutable |
| `portRange` | `string` | — | Port or contiguous range. Immutable |
| `ipProtocol` | `string` | `TCP` | TCP/UDP/ESP/AH/SCTP/ICMP. Immutable |
| `ipVersion` | `string` | `IPV4` | For auto-assigned IPs only. Immutable |
| `loadBalancingScheme` | `string` | `EXTERNAL` | Plus `EXTERNAL_MANAGED`, `INTERNAL_MANAGED`, `INTERNAL_SELF_MANAGED`, `NONE` (PSC). Immutable |
| `network` / `subnetwork` | refs | — | Internal schemes + PSC only (CEL-enforced) |
| `networkTier` | `string` | `PREMIUM` | Global rules are PREMIUM-only (CEL-enforced) |
| `metadataFilters` | list | `[]` | Traffic Director xDS scoping (CEL-enforced) |
| `serviceDirectoryRegistration` | object | — | PSC-for-Google-APIs discovery (CEL-enforced) |
| `noAutomateDnsZone` | `bool` | `false` | Skip the auto-created PSC DNS zone (CEL-enforced) |
| `labels` | map | `{}` | Mutable |
| `externalManagedBackendBucketMigrationState` / `...TestingPercentage` | string / double | — | EXTERNAL → EXTERNAL_MANAGED canary migration |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `ip_address` | The VIP — the value DNS records point at |
| `self_link` | Self-link URI of the rule |
| `forwarding_rule_name` | Name of the rule in GCP |
| `forwarding_rule_id` | Server-assigned numeric ID |
| `psc_connection_id` | PSC connection id (PSC frontends only) |
| `psc_connection_status` | PSC connection status (PSC frontends only) |

## Related Components

- [GcpTargetHttpsProxy](/docs/catalog/gcp/gcptargethttpsproxy) — the default target
- [GcpTargetHttpProxy](/docs/catalog/gcp/gcptargethttpproxy) — the port-80 redirect target
- [GcpGlobalAddress](/docs/catalog/gcp/gcpglobaladdress) — the reserved static VIP
- [GcpVpc](/docs/catalog/gcp/gcpvpc) — the network for internal/PSC frontends
