# GCP Router NAT

Deploys a GCP Cloud Router with a Cloud NAT gateway to provide managed egress for instances without external IPs — internet egress (public NAT) or NAT between private networks (private NAT). Supports auto-allocated or referenced static NAT IPs with connection draining, per-subnetwork scoping, port and timeout tuning, NAT rules, and translation logging.

## What Gets Created

When you deploy a GcpRouterNat resource, Planton provisions:

- **Compute Engine API enablement** — activates `compute.googleapis.com` on the target project so a fresh project works first try
- **Cloud Router** — a regional `google_compute_router` attached to the specified VPC network
- **Cloud NAT Gateway** — a `google_compute_router_nat` on the router with the chosen allocation strategy, subnetwork coverage, tuning, rules, and log settings

Static NAT IPs are referenced [GcpAddress](/docs/catalog/gcp/gcpaddress) reservations, never created by this component — each reservation is its own composable node whose `address` output carries the literal IP.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** — referenced via `projectId` (or the provider's default project)
- **An existing VPC network** — referenced via `vpcSelfLink` (typically a GcpVpc resource)
- **GcpAddress reservations** (EXTERNAL, same region) for stable allowlistable egress IPs (optional — omit `natIps` for auto-allocation)
- **GcpSubnetwork references** to restrict NAT coverage (optional — omit for all subnetworks in the region)

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
      kind: GcpVpc
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
| `vpcSelfLink` | `StringValueOrRef` | The VPC network the router attaches to. | Required. Can reference a GcpVpc. Immutable. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default project | GCP project. Can reference a GcpProject. |
| `type` | `string` | `PUBLIC` | `PUBLIC` (internet egress) or `PRIVATE` (NAT between Network Connectivity Center spokes). Immutable. |
| `natIps` | `StringValueOrRef[]` | `[]` (auto-allocate) | GcpAddress reservations (by self link) for manual allocation — stable, allowlistable egress IPs. |
| `drainNatIps` | `StringValueOrRef[]` | `[]` | IPs being drained out of service: existing connections continue, new ones stop. Must already be in `natIps`. |
| `autoNetworkTier` | `string` | `PREMIUM` | Network tier for auto-allocated IPs. |
| `sourceSubnetworkIpRangesToNat` | `string` | all subnetworks, all ranges | `ALL_SUBNETWORKS_ALL_IP_RANGES`, `ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES`, or `LIST_OF_SUBNETWORKS`. |
| `subnetworks` | `object[]` | `[]` | Per-subnetwork scoping: subnetwork ref + `sourceIpRangesToNat` + `secondaryIpRangeNames`. |
| `minPortsPerVm` / `maxPortsPerVm` | `int32` | 64 (static) / — | Port floor and (with dynamic allocation) ceiling per VM; powers of two under dynamic allocation. |
| `enableDynamicPortAllocation` | `bool` | `false` | Ports grow from min toward max on demand. Mutually exclusive with endpoint-independent mapping. |
| `enableEndpointIndependentMapping` | `bool` | `false` | RFC 5128 endpoint-independent mapping for P2P/SIP workloads. |
| `endpointTypes` | `string[]` | `ENDPOINT_TYPE_VM` | `ENDPOINT_TYPE_VM`, `ENDPOINT_TYPE_SWG`, or `ENDPOINT_TYPE_MANAGED_PROXY_LB` (exactly one). Immutable. |
| `udpIdleTimeoutSec` / `icmpIdleTimeoutSec` / `tcpEstablishedIdleTimeoutSec` / `tcpTransitoryIdleTimeoutSec` / `tcpTimeWaitTimeoutSec` | `int32` | 30 / 30 / 1200 / 30 / 120 | Connection timeout tuning. |
| `rules` | `object[]` | `[]` | NAT rules (`ruleNumber`, CEL `match`, `action` with dedicated IPs or ranges). Lower number wins. |
| `logFilter` | `enum` | `ERRORS_ONLY` | `DISABLED`, `ERRORS_ONLY`, `TRANSLATIONS_ONLY`, or `ALL`. |
| `routerAsn` | `uint32` | GCP-assigned | Private BGP ASN; only for routers that also serve BGP sessions. |
| `routerKeepaliveInterval` | `int32` | 20 | BGP keepalive in seconds (20-60). |

## Examples

### All-Subnets NAT with Auto-Allocated IPs

The recommended default — outbound internet access for every subnetwork in the region with GCP managing the IP pool:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRouterNat
metadata:
  name: uscentral1-nat
spec:
  projectId:
    value: my-gcp-project-123
  routerName: dev-uscentral1-router
  natName: dev-uscentral1-nat
  region: us-central1
  vpcSelfLink:
    valueFrom:
      kind: GcpVpc
      name: my-vpc
      fieldPath: status.outputs.network_self_link
  logFilter: ERRORS_ONLY
```

### Static Egress IPs for Partner Allowlisting

Manual allocation referencing GcpAddress reservations — the egress IPs are stable and third parties can allowlist them:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRouterNat
metadata:
  name: prod-nat
spec:
  projectId:
    value: my-gcp-project-123
  routerName: prod-uscentral1-router
  natName: prod-uscentral1-nat
  region: us-central1
  vpcSelfLink:
    valueFrom:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  natIps:
    - valueFrom:
        kind: GcpAddress
        name: prod-egress-ip-a
        fieldPath: status.outputs.self_link
    - valueFrom:
        kind: GcpAddress
        name: prod-egress-ip-b
        fieldPath: status.outputs.self_link
  minPortsPerVm: 128
  logFilter: ERRORS_ONLY
```

### Subnetwork-Scoped NAT with Secondary-Range Selection

NAT only the GKE pod range of one subnetwork — the node range stays off the internet:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRouterNat
metadata:
  name: gke-egress-nat
spec:
  projectId:
    value: my-gcp-project-123
  routerName: gke-router
  natName: gke-nat
  region: europe-west1
  vpcSelfLink:
    valueFrom:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  sourceSubnetworkIpRangesToNat: LIST_OF_SUBNETWORKS
  subnetworks:
    - subnetwork:
        valueFrom:
          kind: GcpSubnetwork
          name: gke-nodes
          fieldPath: status.outputs.subnetwork_self_link
      sourceIpRangesToNat:
        - LIST_OF_SECONDARY_IP_RANGES
      secondaryIpRangeNames:
        - pods
  logFilter: ALL
```

### Dedicated NAT IP for One Partner (NAT Rules)

Give traffic to a partner's endpoints a dedicated, separately allowlistable source IP while everything else uses the shared pool:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRouterNat
metadata:
  name: partner-nat
spec:
  projectId:
    value: my-gcp-project-123
  routerName: prod-router
  natName: partner-nat
  region: us-central1
  vpcSelfLink:
    valueFrom:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  natIps:
    - valueFrom:
        kind: GcpAddress
        name: shared-egress-ip
        fieldPath: status.outputs.self_link
    - valueFrom:
        kind: GcpAddress
        name: partner-egress-ip
        fieldPath: status.outputs.self_link
  rules:
    - ruleNumber: 100
      match: destination.ip == '203.0.113.10' || destination.ip == '203.0.113.11'
      description: partner API egress pinned to a dedicated IP
      action:
        sourceNatActiveIps:
          - valueFrom:
              kind: GcpAddress
              name: partner-egress-ip
              fieldPath: status.outputs.self_link
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `name` | `string` | Name of the Cloud NAT gateway as created in GCP |
| `router_self_link` | `string` | Self-link URL of the Cloud Router carrying this NAT |
| `nat_ips` | `string[]` | Self links of the static IPs the NAT translates through (manual allocation only; empty for auto-allocation) |

## Related Components

- [GcpVpc](/docs/catalog/gcp/gcpvpc) — provides the VPC network that the Cloud Router attaches to
- [GcpAddress](/docs/catalog/gcp/gcpaddress) — EXTERNAL reservations referenced as stable NAT IPs
- [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork) — subnets that can be scoped for NAT coverage
- [GcpProject](/docs/catalog/gcp/gcpproject) — the GCP project where the router and NAT are created
- [GcpGkeCluster](/docs/catalog/gcp/gcpgkecluster) — private GKE clusters that depend on Cloud NAT for outbound internet access
- [GcpComputeInstance](/docs/catalog/gcp/gcpcomputeinstance) — private VMs that use Cloud NAT for egress
