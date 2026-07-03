# GCP Regional Addresses — Deep Dive

## Regional vs global scope

Google Cloud exposes two address reservation APIs with similar fields but
different scope:

- **`google_compute_global_address`** (`GcpGlobalAddress`) — global scope.
  Used for HTTP(S) load balancer frontends, global VPC peering ranges, and
  Private Service Connect endpoints.
- **`google_compute_address`** (`GcpAddress`) — **regional** scope. Used for
  Cloud NAT, regional load balancers, VM static IPs, internal LB VIPs,
  regional VPC peering ranges, and IPsec interconnect VLAN attachments.

This component models the regional resource. If you need a global static IP
or a global VPC peering range for managed services, use `GcpGlobalAddress`
instead.

## The two faces of regional addresses

### Face 1: External static IPs

Reserve a public IPv4 or IPv6 address in a specific region. Attach it to Cloud
NAT, a regional external load balancer, or a VM instance. The address persists
when the consuming resource is recreated — but only within that region.

**Configuration:** `address_type: EXTERNAL`, required `region`, optionally
`network_tier: PREMIUM|STANDARD` and `ipv6_endpoint_type: VM|NETLB` for IPv6.

### Face 2: Internal regional IPs

Reserve a private IP or CIDR range within a VPC network or subnetwork.
Purposes distinguish the use case:

| Purpose | Requires | Use case |
|---|---|---|
| `GCE_ENDPOINT` | `subnetwork` | VM instances, alias IP ranges |
| `DNS_RESOLVER` | `subnetwork` | DNS resolver endpoints |
| `SHARED_LOADBALANCER_VIP` | — (INTERNAL only) | Internal LB shared VIP |
| `VPC_PEERING` | `network` | Regional VPC peering ranges |
| `IPSEC_INTERCONNECT` | `network` | HA VPN over Cloud Interconnect |

**Note:** `PRIVATE_SERVICE_CONNECT` is global-only — use `GcpGlobalAddress`.

## Schema-level validation

Planton's protobuf schema encodes GCP's cross-field constraints via CEL:

- `network_tier` cannot be set on INTERNAL addresses (internal traffic is
  always Premium tier).
- `purpose` requires INTERNAL `address_type`.
- `VPC_PEERING` / `IPSEC_INTERCONNECT` require `network`.
- `GCE_ENDPOINT` / `DNS_RESOLVER` require `subnetwork`.
- `SHARED_LOADBALANCER_VIP` requires INTERNAL.

These catch the most common misconfiguration errors at authoring time — before
any API call — rather than at apply time in Terraform or Pulumi.

## ForceNew awareness

Every field except labels is ForceNew. Changing `address_name`, `region`,
`address_type`, `network`, `subnetwork`, `purpose`, or `prefix_length`
destroys and recreates the reservation. For EXTERNAL addresses, recreation
means a **new IP** — update DNS, allowlists, and all references before
applying.

## Composition via StringValueOrRef

`project_id`, `network`, and `subnetwork` use Planton's foreign-key pattern:

```yaml
spec:
  network:
    valueFrom:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  subnetwork:
    valueFrom:
      kind: GcpSubnetwork
      name: app-subnet
      fieldPath: status.outputs.subnetwork_self_link
```

Downstream resources reference the address via `status.outputs.self_link`.

## Region output contract

The `region` stack output exports the **plain spec region name** (e.g.
`us-central1`), not a provider self-link. Downstream composition and E2E
verification use this to confirm scope compatibility without parsing URLs.

## Deliberately not modeled (recorded reasons)

- **`ip_collection`** — BYOIP via Public Delegated Prefixes is an enterprise
  workflow with its own provisioning chain; defer until a concrete consumer
  need appears in the catalog.

## Related components

- [GcpGlobalAddress](../gcpglobaladdress/v1/README.md) — global-scope addresses
- [GcpVpc](../gcpvpc/v1/README.md) — VPC network for INTERNAL addresses
- [GcpSubnetwork](../gcpsubnetwork/v1/README.md) — subnetwork for GCE endpoints
- [GcpServiceNetworkingConnection](../gcpservicenetworkingconnection/v1/README.md) — uses global VPC peering ranges, not regional addresses
