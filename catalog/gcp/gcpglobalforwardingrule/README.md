# GCP Global Forwarding Rule

Deploys a global Compute Engine forwarding rule (`google_compute_global_forwarding_rule`) — the VIP node of a global load balancer. The forwarding rule is where traffic enters: it binds an IP address and port to a target proxy. It is also the entry point for Private Service Connect, forwarding a VPC's traffic privately to Google APIs or a producer's service attachment.

## What Gets Created

A single global forwarding rule in the chosen project — the resource DNS records point at (via its `ip_address` output).

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **A target proxy** — a `GcpTargetHttpsProxy` or `GcpTargetHttpProxy` (or another global target's URI)
- **Recommended: a reserved static IP** — a `GcpGlobalAddress`, so the VIP survives frontend rebuilds

## Quick Start

Create a file `forwarding-rule.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
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

Deploy:

```shell
planton apply -f forwarding-rule.yaml
```

This creates a port-443 frontend with a Google-assigned ephemeral IP; add `ipAddress` referencing a `GcpGlobalAddress` for a stable production VIP.

## Configuration Reference

### The frontend

| Field | Description |
|-------|-------------|
| `target` | Required. The proxy receiving traffic (`GcpTargetHttpsProxy` ref by default; HTTP proxy self-link; PSC bundles `all-apis`/`vpc-sc`; service attachment URIs). **Mutable in place** — the zero-downtime frontend swap |
| `ipAddress` | The VIP: a `GcpGlobalAddress` ref (its reserved IP), a literal IP, or empty for an ephemeral IP. Immutable |
| `portRange` | Port or contiguous range (`"443"`, `"8080-8090"`); non-overlapping ranges are how two rules share one IP. Immutable |
| `ipProtocol` | Default `TCP` (what every proxy-based LB and PSC uses). Immutable |
| `ipVersion` | `IPV4` (default) or `IPV6` for auto-assigned ephemeral IPs. Immutable |
| `loadBalancingScheme` | `EXTERNAL` (default), `EXTERNAL_MANAGED`, `INTERNAL_MANAGED`, `INTERNAL_SELF_MANAGED`, or `NONE` (Private Service Connect). Immutable |

### Network wiring (internal schemes + PSC)

| Field | Description |
|-------|-------------|
| `network` | VPC ref — required for PSC, used by internal schemes; rejected on external frontends (CEL) |
| `subnetwork` | Subnetwork ref for internal load balancing |
| `networkTier` | Global rules are `PREMIUM`-only (CEL-enforced) |

### Traffic Director / PSC extras

| Field | Description |
|-------|-------------|
| `metadataFilters` | xDS client scoping — `INTERNAL_SELF_MANAGED` only (CEL) |
| `serviceDirectoryRegistration` | Register a PSC Google-APIs endpoint in Service Directory — scheme `NONE` only (CEL) |
| `noAutomateDnsZone` | Skip the auto-created PSC DNS zone — scheme `NONE` only (CEL) |

### Lifecycle

| Field | Description |
|-------|-------------|
| `labels` | Organize/bill the rule. Mutable |
| `externalManagedBackendBucketMigrationState` / `...TestingPercentage` | The EXTERNAL → EXTERNAL_MANAGED backend-bucket canary migration, without recreating the VIP |
| `deletionPolicy` | What destroy does: `DELETE` (default), `PREVENT` (refuse), or `ABANDON` (keep serving, drop from management) |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `ip_address` | `string` | The VIP — the value DNS records point at (always the literal IP) |
| `self_link` | `string` | Self-link URI of the rule |
| `forwarding_rule_name` | `string` | Name of the rule in GCP |
| `forwarding_rule_id` | `string` | Server-assigned numeric ID |
| `psc_connection_id` | `string` | PSC connection id (PSC frontends only) |
| `psc_connection_status` | `string` | PSC connection status — `ACCEPTED` means the producer admitted the connection |

## Deployment Methods

### Pulumi (Go)

See [`iac/pulumi/README.md`](iac/pulumi/README.md).

### Terraform

See [`iac/tf/README.md`](iac/tf/README.md).

## Important Notes

- **Only `target` and `labels` mutate in place.** Everything else recreates the rule — and an ephemeral VIP changes on recreate, which is why production frontends reference a reserved `GcpGlobalAddress`.
- **The pair pattern**: a port-80 rule (→ `GcpTargetHttpProxy` serving a redirect map) and a port-443 rule (→ `GcpTargetHttpsProxy`) share one static IP.
- **PSC naming**: Private Service Connect rules for Google APIs are limited to 20-character letter/digit names.

## Related Components

- [GcpTargetHttpsProxy](/docs/catalog/gcp/gcptargethttpsproxy) — the default target
- [GcpTargetHttpProxy](/docs/catalog/gcp/gcptargethttpproxy) — the port-80 redirect target
- [GcpGlobalAddress](/docs/catalog/gcp/gcpglobaladdress) — the reserved static VIP
- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork) — the network for internal/PSC frontends
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the GCP project that owns the rule

## Additional Resources

- [Forwarding rule concepts](https://cloud.google.com/load-balancing/docs/forwarding-rule-concepts)
- [Private Service Connect for Google APIs](https://cloud.google.com/vpc/docs/configure-private-service-connect-apis)

## Support

For issues, questions, or contributions, please refer to the Planton documentation or open an issue in the repository.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
