# GcpGkeNodePool — Research and Design Documentation

## 1. Introduction

### What Is a GKE Node Pool?

A node pool is a group of Compute Engine VMs that share one configuration — machine type, disks, identity, scheduling constraints — attached to a GKE cluster. It is the unit of compute heterogeneity in Kubernetes on GCP: one cluster runs several pools, each shaped for its workloads. Frontend services on general-purpose on-demand nodes; fault-tolerant batch on scale-to-zero Spot; ML inference on GPU nodes with GKE-managed drivers; latency-sensitive engines on compact-placed, CPU-pinned pools.

The API makes node pools first-class resources with independent lifecycles: `projects.locations.clusters.nodePools` is its own REST surface, and pools are created, resized, upgraded, and destroyed without ever touching the cluster control plane — the pattern production stability depends on.

### The Composition Boundary

- **`GcpGkeCluster`** owns the control plane and cluster-wide configuration. Its default node pool is always removed at create time.
- **`GcpGkeNodePool`** owns compute. Every pool references the cluster by its `name` and `location` outputs — the names GKE actually assigned — so the reference survives any divergence between the Planton object name and the cloud name.
- Autopilot clusters take no node pools at all: GKE owns nodes there, and the cluster kind enforces that conflict set pre-deploy.

## 2. Deployment Methods Landscape

### Level 0: Cloud Console

The "Add Node Pool" form. Useful for discovering what the API offers; every "who changed that machine type" incident starts here.

### Level 1: gcloud CLI

`gcloud container node-pools create` with dozens of flags. Scriptable but imperative: no drift detection, no plan, and update is a different command surface than create.

### Level 2: Terraform / OpenTofu

`google_container_node_pool` is declarative with plan/apply. Its sharp edges:

- **The autoscaler fight**: for autoscaled pools the live `node_count` drifts by design; without an explicit `ignore_changes`, every plan tries to reset it.
- **The version fight**: an explicit `version` with `auto_upgrade` on makes the API and the plan permanently disagree.
- **Wide immutability**: nearly all of `node_config` is ForceNew — a machine-type edit silently plans a full pool replacement (drain and recreate).
- **Cross-field constraints** (per-zone XOR total autoscaling limits, blue-green settings vs strategy, spot vs preemptible) surface only at apply time.

### Level 3: Pulumi

`container.NodePool` bridges the same schema into real languages — same surface, same sharp edges, plus language-level composition.

### Level 4: Planton

A validated protobuf spec compiled to BOTH engines with identical behavior. The spec turns the API's apply-time failures into manifest-time errors, resolves the parent cluster by reference, wires identity and encryption to first-class resources (`GcpServiceAccount`, `GcpKmsKey`), and documents every immutable field so a wizard or an agent can warn before a destructive edit.

## 3. The Planton Approach

### Parenting by Reference

`cluster_name` resolves from the cluster's `name` output and `location` from its `location` output. A manifest never repeats what the cluster already knows, and the pool addresses its parent exactly the way GKE named it. `project_id` follows the ambient-project contract (empty = provider default) and must be the cluster's project — pools cannot cross projects.

### Pre-Deploy Coherence Rules

| Rule | What it prevents |
|---|---|
| Per-zone (min/max) XOR total (total_min/total_max) autoscaling limits | The provider's mutually exclusive addressing modes, rejected before a deploy |
| Autoscaling requires a maximum; min ≤ max in both modes | Unbounded pools and inverted bounds |
| `initial_node_count` only with autoscaling | A fixed-size pool's size IS node_count — the field would be dead config |
| Blue-green settings require the BLUE_GREEN strategy; surge dials conflict with it | The API's strategy/settings pairing |
| Rollout batch percentage XOR node count | The provider's exclusive batch sizing |
| `spot` XOR `preemptible` | Two capacity models that cannot combine |
| `fast_socket_enabled` requires `gvnic_enabled` | NCCL Fast Socket rides on gVNIC |
| SPECIFIC_RESERVATION requires key + values | The reservation-name pairing the API demands |
| `create_pod_range` requires a `pod_range` name | The new range must be named |
| Enum-valued strings validated against the released provider's accepted values | Typo-driven apply failures |

### Modeled Surface (the 90/10 floor)

Verified against the RELEASED provider line (google 6.50.0 schema dump), not the provider's main branch:

| Family | Modeled |
|---|---|
| Identity & parenting | project (ambient default), cluster name + location by reference, pool name (defaults to metadata.name), node locations, version, max pods per node |
| Sizing | fixed count XOR autoscaling (per-zone or total bounds, scale-to-zero, BALANCED/ANY location policy), initial node count |
| Lifecycle | auto-repair + auto-upgrade (default true), surge upgrades (max_surge/max_unavailable), blue-green upgrades (batch pacing + soak), compact placement (+ TPU topology), queued provisioning |
| Networking | dedicated per-pool pod ranges (create-new or use-existing), per-pool private nodes, TIER_1 egress bandwidth, pod CIDR overprovision control |
| Machine shape | machine type, boot disk size/type (incl. hyperdisk), image type, min CPU platform, local SSDs (SCSI, NVMe ephemeral-storage + data cache, raw-block NVMe), secondary boot disks |
| GPUs | accelerator type/count, GKE-managed driver install, MIG partitioning, time-sharing/MPS sharing, fast socket + gVNIC |
| Identity & security | service account by reference, OAuth scopes, shielded VM options, confidential nodes (SEV/SEV-SNP/TDX), CMEK boot-disk key by reference, workload metadata mode |
| Scheduling | Kubernetes labels, taints (all three effects), GCE network tags, Spot/preemptible, reservation affinity, flex-start, max run duration |
| Node tuning | kubelet (CPU manager, CFS quota, PID limit, read-only port, parallel pulls, log rotation, image GC), Linux sysctls, cgroup mode, hugepages, image streaming (GCFS), logging variant |

### Deliberately Not Modeled (recorded reasons)

| Excluded | Reason |
|---|---|
| `node_drain_config`, `ignore_node_count_changes`, `gpudirect_strategy`, `node_image_config`, `sandbox_config`, `architecture_taint_behavior` | Absent from the released 6.x provider line (main-branch-only). Revisit on the next provider major. |
| `boot_disk` performance block, `containerd_config`, `advanced_machine_features`, `host_maintenance_policy`, `enable_confidential_storage`, `local_ssd_encryption_mode`, `storage_pools`, sole tenancy (`node_group`/`sole_tenant_config`), `windows_node_config`, kubelet eviction surgery, `allowed_unsafe_sysctls`, transparent-hugepage/NUMA-manager knobs | Deep niches without real-world pull relative to their spec weight; each returns on demand. |
| `name_prefix` | Planton owns resource naming. |
| `resource_manager_tags` | Catalog-wide decision — tag bindings await a first-class design. |
| Multi-networking (`additional_node_network_configs`/`additional_pod_network_configs`) | Rides the cluster's `enable_multi_networking` wave when it gains pull. |

## 4. Implementation Notes

### Both Engines, One Contract

- Both modules enable `container.googleapis.com` before creating the pool (`disable_on_destroy = false`) — destroying a pool never disables the API for its own cluster.
- Both translate the spec identically: empty optional strings are omitted (never sent as ""), presence-carrying messages gate their blocks, `management` defaults to auto-repair + auto-upgrade on, and the `disable-legacy-endpoints=true` node metadata is enforced beneath user entries.
- User `resource_labels` merge beneath the platform attribution labels; Kubernetes node `labels` pass through untouched (two different label systems, kept distinct).
- Both engines ignore live `node_count` drift so the cluster autoscaler owns the size at runtime.
- The Terraform module runs on plain `google ~> 6.0` — every modeled field is GA on the released line (verified by schema dump).

### Immutability

The API replaces the pool (draining and recreating every node) when any of these change: name, location, initial_node_count, max_pods_per_node, placement policy, queued provisioning, pod range configuration, and nearly all of node_config — machine type, disks, image, service account, scopes, metadata, spot/preemptible, accelerators, shielded/confidential settings, local SSDs, CMEK key, reservation affinity. Mutable in place: node_count, autoscaling bounds, management, upgrade settings, node_locations, labels, taints, tags, and resource labels. Field comments carry these so a wizard or an agent can warn before a destructive edit.

### Outputs

`node_pool_id` is the fully qualified path downstream services address the pool by (Dataproc on GKE resolves it via FK); `instance_group_urls` are what instance-group LB backends and automation compose against; `min_nodes`/`max_nodes` are the effective bounds regardless of sizing mode; `current_node_count` is a deploy-time snapshot the autoscaler moves afterwards.

## 5. Production Best Practices

1. **Several pools, each with a job** — a general on-demand pool, a tainted Spot pool for batch, GPU pools that scale to zero. One giant pool is the anti-pattern node pools exist to fix.
2. **Taint special-purpose pools** — Spot, GPU, and compliance pools should require explicit tolerations so general workloads never land there by accident.
3. **Scale-to-zero what is expensive** — `minNodes: 0` with `ANY` location policy makes GPU and Spot pools free while idle.
4. **A dedicated node service account** — the Compute default SA is over-privileged; pair a minimal SA with Workload Identity (`GKE_METADATA`) for workload permissions.
5. **Secure boot on** unless a workload loads unsigned kernel modules; integrity monitoring stays on by default.
6. **`maxSurge: 1, maxUnavailable: 0`** for serving pools — upgrades never dip capacity; blue-green for pools where rollback matters more than cost.
7. **Never pin `version` with auto-upgrade on** — the two fight; pin only with auto-upgrade off and own upgrades deliberately.
8. **Plan machine-type changes as replacements** — nearly all of node_config is immutable; roll a new pool and drain the old one for zero-disruption migrations.
