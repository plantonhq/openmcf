# GCP Subnetworks: The Address Plan Everything Else Inherits

## Why the Subnet Is the Real Network Decision

A custom-mode VPC in GCP is nearly free of decisions — it is a named container with a routing mode. The subnet is where network architecture actually happens: which region, how much address space, which secondary ranges for containers, whether IPv6, whether the address space is workload space at all or reserved for load-balancer proxies or Private Service Connect. Every consumer above it — GKE clusters, compute instances, Cloud Run's Direct VPC egress, internal load balancers — inherits the subnet's plan and its constraints.

The constraints are unusually permanent. Name, project, region, network, and even description are immutable: changing any of them destroys and recreates the subnet, which is an outage for everything addressed in it. The one asymmetric knob is the primary range: GCP allows **expanding** `ip_cidr_range` in place (a /24 can become a /20 with zero downtime) but never shrinking. The operational lesson is baked into the sizing guidance: starting small "to be safe" is exactly backwards — undersizing is the mistake you cannot cheaply undo.

## Secondary Ranges: How Containers Get Real IPs

GCP's alias-IP mechanism gives containers first-class VPC addresses instead of NAT'd overlay space, and secondary ranges are its supply. A VPC-native GKE cluster draws pod IPs from one named secondary range and service IPs from another, selecting them **by name** — which is why `range_name` is a contract, not a label, and why this component's outputs export the ranges for FK reference (a GKE cluster's `cluster_secondary_range_name` resolves against them).

Pod ranges are the classic under-sizing trap: GKE reserves a /24 per node by default, so a /20 pod range caps a cluster at 16 nodes regardless of how much node address space remains. The GKE-ready preset ships /14 pods + /20 services for this reason.

Updates to secondary ranges carry a sharp edge the spec defuses: the API treats "no secondary ranges in the request" as ambiguous between "leave them alone" and "remove them all". The `send_secondary_ip_range_if_empty` latch (default false) pins the safe meaning — a partial manifest cannot silently wipe a cluster's pod ranges; removal requires stating the intent.

## Purpose: When a Subnet Is Not Workload Space

Most subnets are `PRIVATE` — regular workload space. The `purpose` field creates the special-role subnets other networking features hard-require:

- **`REGIONAL_MANAGED_PROXY`** — the proxy-only subnet. GCP's regional Application Load Balancers (internal and external) run on a managed Envoy fleet that allocates its addresses from exactly this subnet; a region cannot host a regional ALB until one exists in the VPC there. The `role` field (`ACTIVE`/`BACKUP`) exists for drain-and-swap migrations of that proxy address space. `GLOBAL_MANAGED_PROXY` is the cross-region equivalent.
- **`PRIVATE_SERVICE_CONNECT`** — address space backing published PSC services.
- **`PRIVATE_NAT`** — source ranges for Private NAT gateways; **`PEER_MIGRATION`** — staging space for subnet moves between peered VPCs.

These are one field on this kind rather than separate kinds because GCP models them as one resource: the purpose changes what the address space is *for*, not what the API object *is*.

## IPv6: Dual-Stack as a Stack Type

IPv6 arrives through two fields, and the split matters. `stack_type` (`IPV4_ONLY` → `IPV4_IPV6` → `IPV6_ONLY`) decides what address families the subnet carries — and moving between IPv4-only and dual-stack is an in-place update, so retrofitting IPv6 onto a live subnet is safe. `ipv6_access_type` decides what the IPv6 addresses *are*: `EXTERNAL` assigns internet-routable GUAs from Google's space (optionally pinned with `external_ipv6_prefix`); `INTERNAL` assigns ULAs routable only inside the VPC and requires the VPC to have its ULA range enabled first. The access type is immutable — it is the one IPv6 decision to get right up front. The spec enforces the pairing rules (IPv6 stack types require an access type; `IPV4_ONLY` forbids one; IPv6-only subnets carry no IPv4 range).

## Flow Logs: Observability with a Bill

VPC Flow Logs sample the subnet's TCP/UDP flows into Cloud Logging — the raw material for network forensics, egress-cost analysis, and compliance (PCI/SOC 2 commonly require them). The spec models the full provider surface: aggregation interval (fewer, larger entries at longer intervals), sampling fraction (1.0 for forensics, lower for trends), metadata scope (all / none / a custom field list), and a CEL `filter_expr` for logging only the flows that matter. The cost lever is real — full sampling on a busy subnet is a significant Logging line item, which is why logging is off unless the block is present and why the defaults mirror GCP's own (50% sampling, 5-second aggregation).

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `name` / `network` / `region` / `project` | ✅ | `vpc_self_link` → GcpVpc ref; project falls back to the provider default |
| `ip_cidr_range` | ✅ | Optional only for IPV6_ONLY (spec-enforced); expand-don't-shrink documented |
| `description` | ✅ | Immutable on this resource (provider ForceNew) |
| `purpose` + `role` | ✅ | All six purposes; role gated to REGIONAL_MANAGED_PROXY pre-deploy |
| `secondary_ip_range` | ✅ | name + CIDR; the `send_secondary_ip_range_if_empty` latch modeled |
| `private_ip_google_access` / `private_ipv6_google_access` | ✅ | |
| `stack_type` / `ipv6_access_type` / `external_ipv6_prefix` | ✅ | Pairing rules enforced pre-deploy |
| `allow_subnet_cidr_routes_overlap` | ✅ | Preview-stage on the released provider line — the TF module selects google-beta |
| `log_config` (full block) | ✅ | All five fields incl. CUSTOM_METADATA gating |
| `reserved_internal_range` (top-level + per-secondary) | ❌ | Internal-range reservations (Network Connectivity Center) are a niche address-management layer; model when that family arrives |
| `ip_collection` (BYOIP) | ❌ | Bring-your-own-IP public delegated prefixes — enterprise-BYOIP niche |
| `resolve_subnet_mask` | ❌ | ARP behavior for PEER_MIGRATION staging only |
| `params.resource_manager_tags` | ❌ | Write-only tag bindings, unmodeled across the catalog; adopting them is a catalog-wide decision |
| `enable_flow_logs` | ❌ | Deprecated in the provider; `log_config` presence is the honest on/off |
| computed: `gateway_address`, `subnetwork_id`, IPv6 prefixes | ✅ outputs | `state`/`fingerprint` omitted (drain-state and optimistic-locking internals nothing composes on) |

## Composition

The subnet sits one layer above the VPC and one below everything compute:

1. **GcpVpc** — the parent network (`vpc_self_link` references its `network_self_link` output).
2. **GcpSubnetwork** (this component) — the address plan.
3. **GcpGkeCluster** — references the subnet's self-link and selects pod/service secondary ranges by name; **GcpComputeInstance**, **GcpCloudRun** (Direct VPC egress by subnet name), **Dataproc**, **Composer**, and **Vertex AI** consumers reference it the same way.

Pair with **GcpRouterNat** for outbound internet from private-only subnets, and create a `REGIONAL_MANAGED_PROXY` sibling before any regional Application Load Balancer lands in the region.
