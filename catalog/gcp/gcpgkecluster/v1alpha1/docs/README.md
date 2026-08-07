# GcpGkeCluster — Research and Design Documentation

## 1. Introduction

### What Is a GKE Cluster?

Google Kubernetes Engine is GCP's managed Kubernetes: Google runs the control plane (API server, etcd, scheduler, controllers), handles its availability and upgrades, and integrates the cluster with the surrounding cloud — VPC networking, IAM, Cloud Logging/Monitoring, KMS, Pub/Sub, BigQuery. A "cluster" in the GKE API is far more than a control plane switch: `google_container_cluster` carries the entire cluster-wide configuration surface — networking mode and IP allocation, private topology, control-plane endpoints and access control, release channels and maintenance scheduling, node auto-provisioning, encryption and supply-chain security, per-component observability, the addon set, fleet membership, and the Autopilot operating mode.

### Two Operating Modes

- **Standard**: you own node pools — machine types, GPUs, spot mixes, per-pool autoscaling. Node pools are separate API resources with independent lifecycles.
- **Autopilot**: GKE provisions and manages nodes invisibly, bills per pod resource request, and enforces a hardened posture (shielded nodes, Workload Identity, Dataplane V2 — always on). Autopilot is the same API resource with `enable_autopilot: true` and a set of node-management fields the API rejects.

## 2. Deployment Methods Landscape

### Level 0: Cloud Console

Point-and-click cluster creation. Useful for exploring what the API offers; unusable for fleets — no reproducibility, no review, and the console quietly makes choices (default node pool, default network) that fight later automation.

### Level 1: gcloud CLI

`gcloud container clusters create` with ~200 flags. Scriptable but imperative: no drift detection, no plan/preview, and cluster mutation is a different command surface than creation.

### Level 2: Terraform / OpenTofu

`google_container_cluster` is one of the largest resources in the provider (dozens of top-level attributes and nested blocks). Declarative with plan/apply and drift detection. The sharp edges the raw resource leaves exposed:

- The API mandates a default node pool at create time; the community idiom (`remove_default_node_pool = true` + `initial_node_count = 1`) must be known, not discovered.
- `deletion_protection` defaults to true and fails destroys — surprising in CI until understood, dangerous when reflexively hardcoded to false.
- Cross-field constraints (Autopilot conflicts, maintenance-window exclusivity, secondary-range-name vs CIDR-block exclusivity) surface only at apply time — after a 20-minute cluster create has already started.
- Inline `node_pool` blocks are ForceNew on the cluster: adding a pool inline recreates the entire cluster.

### Level 3: Pulumi

`container.Cluster` bridges the same provider schema into real languages. Same surface, same sharp edges, plus language-level composition.

### Level 4: Planton

A validated protobuf spec compiled to BOTH engines with identical behavior. The spec encodes the API's cross-field rules as pre-deploy validations, removes the default-node-pool trap entirely (always removed; node pools are `GcpGkeNodePool` resources), keeps GCP's safety defaults honest (`deletion_protection` true by default), and wires every cross-resource input as a typed reference (`GcpVpcNetwork`, `GcpSubnetwork`, `GcpKmsKey`, `GcpServiceAccount`, `GcpPubSubTopic`, `GcpBigQueryDataset`).

## 3. The Planton Approach

### Composition Boundaries

- **Cluster ⟷ node pools**: node pools are first-class `GcpGkeNodePool` resources referencing the cluster's `name` output. The cluster never carries node `node_config`; the default pool is always removed at create. This mirrors the API's own lifecycle boundary (inline pools are ForceNew; the separate node-pool resource is not).
- **Cluster ⟷ network**: `network` and `subnetwork` are required references. The provider would default to the legacy auto-created "default" network; the spec deliberately does not — default-network clusters are a documented GCP anti-pattern and do not compose into reviewable infrastructure. This is an above-floor requiredness decision, not a coverage gap.
- **Cluster ⟷ NAT**: GCP's cluster API has no NAT input. Private-node egress composes as a `GcpRouterNat` on the network; charts order the cluster on the NAT with an explicit dependency, not a phantom spec field.
- **Autopilot as a mode, not a kind**: Autopilot is the same provider resource with one immutable flag and a conflict set. Splitting it into its own kind would duplicate the entire cluster surface to express one boolean; the conflict set is enforced by validation rules instead.

### Pre-Deploy Coherence Rules

The spec turns the API's apply-time failures into manifest-time errors:

| Rule | What it prevents |
|---|---|
| Autopilot ⟂ cluster_autoscaling / default_max_pods_per_node / intranode visibility / Calico / shielded-nodes flag / dns-cache & stateful-ha addons | The API's Autopilot conflict set, rejected before a 20-minute create begins |
| `allow_net_admin` requires Autopilot | The field is Autopilot-only |
| ip_allocation: range NAME xor CIDR block (per range) | The provider's ConflictsWith pairs |
| private endpoint requires private nodes | API rejection |
| `master_ipv4_cidr_block` xor `private_endpoint_subnetwork` | Peering-based vs PSC-based control planes are exclusive |
| maintenance: exactly one of daily/recurring window | The provider's ExactlyOneOf |
| NAP enabled requires resource limits (and limits require enabled) | An unbounded NAP is an unbounded bill; orphaned limits are dead config |
| `ENCRYPTED` state requires a KMS key | API rejection |
| notifications enabled require a topic; usage export requires a dataset | API rejections |
| Enum-valued strings validated against the released provider's accepted values | Typo-driven apply failures |

### Modeled Surface (the 90/10 floor)

Verified against the RELEASED provider line (google 6.50.0 schema dump), not the provider's main branch:

| Family | Modeled |
|---|---|
| Core | name (defaults to metadata.name), location, description, node_locations, resource labels, deletion_protection (default true), fleet registration |
| Networking | VPC-native ip_allocation (named ranges xor GKE-managed CIDRs, dual-stack, additional pod ranges, overprovision toggle), datapath provider (Dataplane V2), intranode visibility, L4 ILB subsetting, FQDN + Cilium cluster-wide policy, multi-networking, private IPv6 Google access, in-transit pod encryption, default SNAT disable, Calico network policy, Cloud DNS config, Gateway API channel, service external IPs, egress bandwidth tier, L4 LB firewall reconciliation toggle |
| Private topology | private nodes, private-only endpoint, peering-based (/28) and PSC-based control planes, master global access, master authorized networks (+ GCP-public-CIDR and private-endpoint enforcement), control-plane DNS/IP endpoints |
| Upgrades | release channels (incl. EXTENDED), min master version, daily/recurring maintenance windows, exclusions with scopes |
| Autoscaling | node auto-provisioning (limits, profile, locations, full auto-provisioning defaults), vertical pod autoscaling, HPA profile |
| Security | Workload Identity (default on), shielded nodes, CMEK etcd encryption, Binary Authorization, Security Posture, authenticator groups, legacy ABAC, mesh certificates, Secret Manager CSI, confidential nodes, anonymous-auth mode, identity service |
| Observability | logging/monitoring component lists, managed Prometheus + auto-monitoring, dataplane observability, Pub/Sub notifications, cost allocation, BigQuery usage export |
| Addons | HTTP LB, HPA, PD/Filestore/GCS-Fuse CSI, Backup for GKE, NodeLocal DNSCache, Config Connector, Stateful HA, Ray operator |
| Autopilot | enable_autopilot, allow_net_admin, full conflict-set validation |

### Deliberately Not Modeled (recorded reasons)

| Excluded | Reason |
|---|---|
| Default-pool `node_config`, inline `node_pool`, `node_pool_defaults`, `node_pool_auto_config` | Node-level configuration belongs to `GcpGkeNodePool`; the default pool is always removed; inline pools are ForceNew (hostile to composition). NAP node defaults ARE modeled (cluster-level concern). |
| `deletion_policy`, `ignore_node_count_changes`, `dataplane_optimization_mode`, `autopilot_privileged_admission`, `secret_sync_config`, `autopilot_cluster_policy_config`, `node_creation_config`, `maintenance_policy.disruption_budget`, `ip_allocation_policy.auto_ipam_config`/`network_tier_config`, pod-snapshot / slice-controller / slurm-operator / agent-sandbox addons | Absent from the released 6.x provider line (main-branch-only). Revisit on the next provider major. |
| `tpu_config`, `pod_security_policy_config`, `cluster_telemetry`, `protect_config`, `managed_opentelemetry_config`, `workload_alts_config` | Beta-provider-only on the released line. |
| `master_auth.client_certificate_config` | Legacy client-cert auth; deprecated practice superseded by IAM. |
| `user_managed_keys_config` | Bring-your-own cluster CA / etcd keys — deep-compliance niche with heavy key-management burden; CMEK etcd encryption covers the real requirement. |
| `enterprise_config` | Deprecated on the provider. |
| `rbac_binding_config`, `enable_k8s_beta_apis`, `gke_auto_upgrade_config`, `enable_kubernetes_alpha`, `enable_tpu`, routes-based `cluster_ipv4_cidr` | Niche/legacy without real-world pull; each returns on demand. |
| `logging_service` / `monitoring_service` strings | Legacy; superseded by the component configs. |

## 4. Implementation Notes

### Both Engines, One Contract

- The Terraform module runs on plain `google ~> 6.0` — every modeled field is GA on the released line (verified by schema dump), so no beta provider dependency exists to drift.
- Both modules enable `container.googleapis.com` before creating the cluster (`disable_on_destroy = false`), so a fresh project works on the first deploy without disabling the API for neighbors on destroy.
- Both modules translate the spec identically: empty optional strings are omitted (never sent as ""), presence-carrying messages gate their blocks, the spec's `NONE` release channel maps to the provider's `UNSPECIFIED`, and the Workload Identity pool name is computed from the effective project (ambient when `project_id` is empty).
- `remove_default_node_pool=true` + `initial_node_count=1` on Standard clusters; both fields are omitted on Autopilot (the API rejects them).
- User `resource_labels` merge beneath the platform attribution labels, so a user label can never clobber attribution.

### Immutability

The API replaces the cluster (and everything on it) when any of these change: name, location, description, network, subnetwork, the whole ip_allocation block, datapath provider, default max pods per node, confidential nodes, Autopilot mode, and the private control-plane placement fields. `enable_l4_ilb_subsetting` is one-way. Field comments carry these so a wizard or an agent can warn before a destructive edit.

### Outputs

`endpoint` + `cluster_ca_certificate` are what kubeconfig-builders need; `workload_identity_pool` is what `GcpGkeWorkloadIdentityBinding` composes against; `name`/`location` are the handles `GcpGkeNodePool` resolves by default; `cluster_id` is the fully-qualified ID downstream services (e.g. Dataproc on GKE) reference. The CA certificate is public trust material, not a secret — clients install it as a trust anchor.

## 5. Production Best Practices

1. **Regional for production** — zonal control planes take a brief API outage during upgrades.
2. **Private nodes + Cloud NAT** — no public node IPs; compose a `GcpRouterNat` for egress.
3. **Dataplane V2** (`ADVANCED_DATAPATH`) — native NetworkPolicy without Calico, observability, FQDN policy.
4. **Named secondary ranges** — plan pod/service space on the subnetwork; GKE-managed ranges are for sandboxes.
5. **Master authorized networks or DNS-endpoint-only** — never leave a public endpoint open to the internet, even credentialed.
6. **Keep `deletion_protection` on** for anything real; presets turn it off only for sandboxes.
7. **A maintenance window plus exclusions** — make upgrade timing a decision, not a surprise; use exclusions for change freezes.
8. **CMEK etcd encryption + Binary Authorization** where compliance demands it; both are modeled and composable with `GcpKmsKey`.
9. **Autopilot for teams that should not own nodes** — the hardened default; Standard when machine control (GPUs, spot, custom images) matters.
