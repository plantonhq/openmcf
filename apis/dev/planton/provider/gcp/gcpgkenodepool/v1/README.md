# GCP GKE Node Pool

Deploys a Google Kubernetes Engine node pool via `google_container_node_pool` (Terraform) or Pulumi `container.NodePool` — a group of identically configured Compute Engine VMs attached to a GKE cluster, with the full node configuration surface: machine shape, disks, identity, scheduling constraints, GPU accelerators, security posture, and node-level tuning.

## Overview

Node pools are the compute half of the GKE composition boundary: the `GcpGkeCluster` owns the control plane and cluster-wide configuration; every node pool is a first-class resource with its own lifecycle, referencing the cluster by its `name` and `location` outputs. Production clusters run several pools — general-purpose on-demand, scale-to-zero Spot for fault-tolerant batch, GPU pools for ML — each sized and tainted for its workloads.

The pool inherits its project and location from the parent cluster (both resolve by reference), and honors the ambient-project contract: an empty `project_id` falls back to the provider's default project.

## Purpose

- **Heterogeneous compute, honest boundaries**: match infrastructure to workloads — pools are created, resized, and destroyed without ever touching the cluster control plane.
- **Pre-deploy coherence**: the API's cross-field rules (per-zone XOR total autoscaling limits, surge XOR blue-green settings, spot XOR preemptible, fast-socket-needs-gVNIC, specific-reservation key pairing) are enforced by validation before any cloud call.
- **Cost levers as first-class fields**: Spot VMs, scale-to-zero autoscaling, flex-start, max run duration, and reservation consumption are all modeled — not hidden behind escape hatches.

## Key Features

- **Sizing**: fixed count or cluster-autoscaler management with per-zone (`min_nodes`/`max_nodes`) or total (`total_min_nodes`/`total_max_nodes`) bounds and the `BALANCED`/`ANY` location policy
- **Upgrades**: surge (max_surge/max_unavailable) or blue-green strategy with batch rollout pacing and soak windows
- **Machine shape**: machine type, boot disk size/type (incl. hyperdisk), image type, min CPU platform, local SSDs (SCSI count, NVMe ephemeral-storage backing with data cache, raw-block NVMe)
- **GPUs**: accelerator type/count, GKE-managed driver installation, MIG partitioning, time-sharing/MPS GPU sharing, NCCL Fast Socket + gVNIC for distributed training
- **Scheduling**: Kubernetes node labels, taints with all three effects, network tags, compact placement (+ TPU topology), queued provisioning (Dynamic Workload Scheduler), flex-start, max run duration
- **Identity & security**: node service account by reference, OAuth scopes, shielded VM options, confidential nodes (SEV/SEV-SNP/TDX), CMEK boot-disk encryption by `GcpKmsKey` reference, workload metadata mode
- **Networking**: dedicated per-pool pod ranges (create-new or use-existing), per-pool private nodes override, TIER_1 egress bandwidth, pod CIDR overprovision control
- **Node tuning**: kubelet (CPU manager, CFS quota, PID limits, log rotation, image GC, the insecure read-only port), Linux sysctls, cgroup mode, hugepages, logging variant, image streaming (GCFS)
- **Capacity**: Spot and legacy preemptible VMs, Compute Engine reservation affinity, secondary boot disks for image preloading

## Stack Outputs

| Output | Description |
|---|---|
| `node_pool_name` | Name of the pool as created in GKE |
| `instance_group_urls` | Managed instance group URLs backing the pool (one per zone for regional clusters) |
| `min_nodes` | Effective minimum size (autoscaling minimum, or the fixed count) |
| `max_nodes` | Effective maximum size (autoscaling maximum, or the fixed count) |
| `current_node_count` | Nodes per zone at the last deploy (a snapshot — the autoscaler moves it) |
| `node_pool_id` | Fully qualified resource ID: `projects/{project}/locations/{location}/clusters/{cluster}/nodePools/{name}` |
| `location` | The pool's location (the parent cluster's region or zone), exactly as provided in the spec |
| `version` | Kubernetes version running on the pool's nodes |

## Deliberately not modeled (recorded reasons)

| Excluded Feature | Why |
|---|---|
| `node_drain_config`, `ignore_node_count_changes`, `gpudirect_strategy`, `node_image_config`, `sandbox_config` (gVisor), `architecture_taint_behavior` | Exist only on the provider's unreleased main line, not on the released 6.x major the GCP modules pin. Revisit on the next provider major. |
| `boot_disk` block (provisioned IOPS/throughput), `containerd_config` (private registry CAs), `advanced_machine_features`, `host_maintenance_policy`, `enable_confidential_storage`, `local_ssd_encryption_mode`, `storage_pools`, `node_group` + `sole_tenant_config` (sole tenancy), `windows_node_config`, kubelet eviction tuning (`eviction_soft`/`eviction_minimum_reclaim`/grace periods), `allowed_unsafe_sysctls`, transparent-hugepage knobs, `memory_manager`/`topology_manager`, `crash_loop_back_off`, `single_process_oom_kill` | Deep niches (sole tenancy, Windows pools, hyperdisk performance tuning, NUMA pinning, kubelet eviction surgery) without real-world pull relative to their spec weight. Each returns on demand — the modeled kubelet/Linux surface covers the tuning production pools actually reach for. |
| `name_prefix` | Planton owns resource naming — the pool name is `node_pool_name` or `metadata.name`, never generated. |
| `resource_manager_tags` | Catalog-wide decision — tag bindings are an org-governance layer pending a first-class design. |
| `additional_node_network_configs` / `additional_pod_network_configs`, `network_config.subnetwork` | Multi-networking (extra node NICs / pod networks) is an advanced cluster feature modeled when the cluster's `enable_multi_networking` gains real pull; the pool's subnetwork override rides the same wave. |

## Related Components

- **GcpGkeCluster** — the control plane this pool attaches to (references its `name` and `location` outputs)
- **GcpServiceAccount** — the node identity (`service_account` reference)
- **GcpKmsKey** — CMEK key for boot-disk encryption
- **GcpGkeWorkloadIdentityBinding** — grants workloads IAM identities; pair with `GKE_METADATA` workload metadata mode
- **GcpDataprocCluster** — a Dataproc-on-GKE virtual cluster schedules Dataproc roles onto named node pools by `node_pool_id`

## Additional Resources

- [Node Pools Documentation](https://cloud.google.com/kubernetes-engine/docs/concepts/node-pools)
- [NodePools REST API](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters.nodePools)
