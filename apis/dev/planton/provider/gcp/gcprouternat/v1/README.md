# GCP Router NAT

Deploys a GCP Cloud Router with a Cloud NAT gateway (`google_compute_router` + `google_compute_router_nat`) — the managed egress path that lets instances without external IPs reach the internet (public NAT) or other private networks (private NAT). Covers IP allocation (auto or referenced static reservations, with connection draining), per-subnetwork scoping, port tuning, connection timeouts, NAT rules, and translation logging.

## What Gets Created

When you deploy a GcpRouterNat resource, Planton provisions:

- **Compute Engine API enablement** — a `google_project_service` resource that activates `compute.googleapis.com` on the target project
- **Cloud Router** — a regional `google_compute_router` attached to the specified VPC network, optionally with a BGP ASN and keepalive interval
- **Cloud NAT Gateway** — a `google_compute_router_nat` on the router, configured with the chosen allocation strategy, subnetwork coverage, port/timeout tuning, rules, and log settings

Static NAT IPs are **referenced, never created**: each `natIps` entry points at a [GcpAddress](/docs/catalog/gcp/gcpaddress) reservation, which is its own composable node with its own lifecycle. The literal IP of each entry is that address node's `address` output.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** — referenced via `projectId` (or the provider's default project)
- **An existing VPC network** — referenced via `vpcSelfLink` (typically a GcpVpcNetwork resource)
- **GcpAddress reservations** (EXTERNAL, same region) if you need stable, allowlistable egress IPs — omit `natIps` entirely for auto-allocation
- **GcpSubnetwork references** if you want to restrict NAT to specific subnetworks — omit for region-wide coverage

## Quick Start

Create a file `router-nat.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRouterNat
metadata:
  name: my-nat
spec:
  projectId:
    value: my-gcp-project-123
  routerName: egress-router
  natName: egress-nat
  region: us-central1
  vpcSelfLink:
    valueFrom:
      kind: GcpVpcNetwork
      name: my-vpc
      fieldPath: status.outputs.network_self_link
```

Deploy:

```shell
planton apply -f router-nat.yaml
```

This creates a Cloud Router and NAT gateway covering all subnetworks in `us-central1` with auto-allocated IPs and `ERRORS_ONLY` logging.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `routerName` | `string` | Name of the Cloud Router to create. | RFC 1035, 1-63 chars. Immutable. |
| `natName` | `string` | Name of the NAT configuration on the router. | RFC 1035, 1-63 chars. Immutable. |
| `region` | `string` | GCP region for the router and NAT. | Required. Immutable. |
| `vpcSelfLink` | `StringValueOrRef` | The VPC network the router attaches to. | Required. Can reference a GcpVpcNetwork. Immutable. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default project | GCP project. Can reference a GcpProject. |
| `type` | `string` | `PUBLIC` | `PUBLIC` (internet egress) or `PRIVATE` (NAT between NCC spokes, no external IPs). Immutable. |
| `natIps` | `StringValueOrRef[]` | `[]` (auto-allocate) | GcpAddress reservations (by self link) the NAT translates through. Non-empty selects manual allocation — the shape for stable, allowlistable egress IPs. |
| `drainNatIps` | `StringValueOrRef[]` | `[]` | IPs being drained: existing connections continue, new connections stop. Each entry must already be in `natIps`. |
| `autoNetworkTier` | `string` | `PREMIUM` | Network tier for auto-allocated IPs: `PREMIUM` or `STANDARD`. |
| `sourceSubnetworkIpRangesToNat` | `string` | all subnetworks, all ranges | `ALL_SUBNETWORKS_ALL_IP_RANGES`, `ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES`, or `LIST_OF_SUBNETWORKS` (implied by listing `subnetworks`). |
| `subnetworks` | `object[]` | `[]` | Per-subnetwork scoping: subnetwork ref + which ranges to NAT (`ALL_IP_RANGES`, `PRIMARY_IP_RANGE`, `LIST_OF_SECONDARY_IP_RANGES` + names). |
| `minPortsPerVm` | `int32` | 64 static / 32 dynamic | Minimum ports reserved per VM. Power of two when dynamic allocation is on. |
| `maxPortsPerVm` | `int32` | — | Port ceiling per VM; only with dynamic allocation. Power of two. |
| `enableDynamicPortAllocation` | `bool` | `false` | Grow a busy VM's ports from min toward max on demand. Mutually exclusive with endpoint-independent mapping. |
| `enableEndpointIndependentMapping` | `bool` | `false` | Same internal ip:port maps to the same NAT ip:port regardless of destination (RFC 5128). |
| `endpointTypes` | `string[]` | `ENDPOINT_TYPE_VM` | Which resource class this NAT serves: `ENDPOINT_TYPE_VM`, `ENDPOINT_TYPE_SWG`, or `ENDPOINT_TYPE_MANAGED_PROXY_LB` (exactly one). Immutable. |
| `udpIdleTimeoutSec` | `int32` | 30 | UDP idle timeout. |
| `icmpIdleTimeoutSec` | `int32` | 30 | ICMP idle timeout. |
| `tcpEstablishedIdleTimeoutSec` | `int32` | 1200 | Established-TCP idle timeout. |
| `tcpTransitoryIdleTimeoutSec` | `int32` | 30 | Half-open TCP idle timeout. |
| `tcpTimeWaitTimeoutSec` | `int32` | 120 | TIME_WAIT linger before a NAT port is reusable. |
| `rules` | `object[]` | `[]` | NAT rules: route matching egress (CEL `match` on destination) through dedicated IPs (public) or subnetwork ranges (private). Lower `ruleNumber` wins. |
| `logFilter` | `enum` | `ERRORS_ONLY` | `DISABLED`, `ERRORS_ONLY`, `TRANSLATIONS_ONLY`, or `ALL`. |
| `routerAsn` | `uint32` | GCP-assigned | Private BGP ASN (64512-65534 or 4200000000+); only needed when the router also serves BGP sessions. |
| `routerKeepaliveInterval` | `int32` | 20 | BGP keepalive in seconds (20-60). |

### Validation Rules

- **Dynamic port allocation and endpoint-independent mapping are mutually exclusive** — the API rejects the combination; the spec rejects it pre-deploy.
- **`maxPortsPerVm` requires dynamic allocation**, must be ≥ `minPortsPerVm`, and both must be powers of two when dynamic allocation is on.
- **Private NAT carries no external IPs** — `natIps`, `drainNatIps`, and `autoNetworkTier` must be empty when `type` is `PRIVATE`; private rules use subnetwork ranges, public rules use IPs.
- **`drainNatIps` requires `natIps`** — an IP must be in the manual set before it can be drained (and the API rejects drain entries on a brand-new NAT).
- **`LIST_OF_SUBNETWORKS` and the `subnetworks` list imply each other.**

## Zero-Downtime Operations

Everything except names, region, network, `endpointTypes`, and `type` updates in place:

- **NAT IP rotation**: add the new GcpAddress ref to `natIps`, move the old one to `drainNatIps`, wait for connections to bleed off, then remove it.
- **Fleet-wide egress tuning**: port floors/ceilings, timeouts, rules, and logging all apply to a live gateway without disturbing traffic.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `name` | `string` | Name of the Cloud NAT gateway as created in GCP |
| `router_self_link` | `string` | Self-link URL of the Cloud Router carrying this NAT |
| `nat_ips` | `string[]` | Self links of the static IPs the NAT translates through (manual allocation only; empty for auto-allocation, where GCP manages an unlisted pool) |

## Deployment Methods

### Pulumi (Go)

See [`iac/pulumi/README.md`](iac/pulumi/README.md).

### Terraform

See [`iac/tf/README.md`](iac/tf/README.md).

## Important Notes

- **One node, two resources**: the router and the NAT are provisioned together — the 90% case is exactly one NAT per router. A router without NAT (dedicated BGP for Interconnect/VPN, with interfaces and peers) is a different concern and is not modeled here.
- **Auto-allocated IPs can change** as GCP scales the pool. When third parties allowlist your egress IPs, use `natIps` with GcpAddress reservations — the literal IPs live on those nodes' outputs and survive NAT reconfiguration.
- **Pair with Private Google Access** on your subnetworks so traffic to Google APIs stays internal and never consumes NAT ports.

### Deliberately not modeled (recorded reasons)

| Excluded Feature | Why |
|---|---|
| Router interfaces / BGP peers / MD5 auth keys | Interconnect/VPN BGP surface belongs to a dedicated router kind if demand appears; NAT-only routers never need it. |
| `initial_nat_ips` + `google_compute_router_nat_address` | An alternative address-attachment workflow for incremental IP management; the `natIps`/`drainNatIps` in-place update path covers rotation without a second resource pattern. |
| NAT64 (`source_subnetwork_ip_ranges_to_nat64`, `nat64_subnetwork`) | IPv6-to-IPv4 translation for IPv6-only subnets is a niche topology; defer until a concrete consumer need appears. |
| Router `params.resource_manager_tags` | Absent from the released provider line the modules pin. |
| `deletion_policy` | Client-side lever that conflicts with Planton-managed destroy. |

## Related Components

- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork) — provides the VPC network the router attaches to
- [GcpAddress](/docs/catalog/gcp/gcpaddress) — EXTERNAL reservations referenced as stable NAT IPs
- [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork) — subnetworks scoped for NAT coverage
- [GcpProject](/docs/catalog/gcp/gcpproject) — the GCP project where the router and NAT are created
- [GcpGkeCluster](/docs/catalog/gcp/gcpgkecluster) — private GKE clusters that depend on Cloud NAT for outbound internet access

## Additional Resources

- [Cloud NAT overview](https://cloud.google.com/nat/docs/overview)
- [Routers API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/routers)
