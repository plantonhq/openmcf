# GCP GKE Cluster

Deploys a Google Kubernetes Engine cluster via `google_container_cluster` (Terraform) or Pulumi `container.Cluster` — the managed Kubernetes control plane plus every cluster-wide configuration surface: networking mode, private topology, upgrades and maintenance, autoscaling, security, observability, addons, and Autopilot.

## Overview

One GcpGkeCluster resource is one cluster. Node pools are separate composable resources — `GcpGkeNodePool` — referencing this cluster by its `name` output; the cluster's default node pool is always removed at create time so every pool is an explicitly managed, first-class node. In Autopilot mode GKE manages nodes entirely (billed per pod) and no node pool resources attach.

Clusters are always VPC-native (alias IP): pods and services draw from secondary ranges on the referenced subnetwork — either ranges you name in `ip_allocation`, or ranges GKE creates and manages when none are named. Legacy routes-based networking is deliberately not modeled.

## Purpose

- **Full control-plane surface, honest boundaries**: the cluster owns cluster-wide concerns; compute lives in `GcpGkeNodePool`; NAT, addresses, and subnets are their own composable nodes.
- **Pre-deploy coherence**: the GKE API's cross-field rules (Autopilot conflicts, window exclusivity, NAP limits, CMEK key pairing, private-endpoint requirements) are enforced by validation before any cloud call — a 20-minute cluster create never fails on a rule the manifest already broke.
- **Safety defaults that match GCP**: `deletion_protection` defaults to true (a destroy fails until it is explicitly turned off), shielded nodes follow GCP's default-on posture, and Workload Identity is on by default.

## Key Features

- **Two modes, one kind**: Standard (you own node pools) and Autopilot (`enable_autopilot`; GKE owns nodes) with the API's conflict set enforced pre-deploy
- **Private clusters**: private nodes, peering-based (`master_ipv4_cidr_block`) or PSC-based (`private_endpoint_subnetwork`) control planes, master global access, private-only endpoints
- **Control-plane access**: master authorized networks (CIDR allowlist, GCP-public-CIDR and private-endpoint enforcement toggles) and the modern DNS endpoint (`control_plane_endpoints`)
- **Upgrades on rails**: release channels incl. EXTENDED, `min_master_version`, daily/recurring maintenance windows, up to 20 maintenance exclusions with scopes
- **Node auto-provisioning (NAP)**: resource limits (the cost brake, required when enabled), autoscaling profiles, full auto-provisioning defaults (SA ref, disks, image, shielding, auto-repair/upgrade)
- **Dataplane V2**: `datapath_provider`, FQDN / Cilium cluster-wide network policy, dataplane observability metrics and relay, in-transit pod traffic encryption
- **Security**: CMEK etcd encryption (KMS key ref), Binary Authorization, Security Posture dashboard, authenticator groups, confidential nodes (SEV/SEV-SNP/TDX), anonymous-auth hardening, mesh certificates, Secret Manager CSI
- **Observability**: per-component logging/monitoring, managed Prometheus (+ auto-monitoring scope), Pub/Sub lifecycle notifications, cost allocation, BigQuery usage export
- **Addons**: HTTP LB, HPA, PD/Filestore/GCS-Fuse CSI drivers, Backup for GKE, NodeLocal DNSCache, Config Connector, Stateful HA, Ray operator
- **Fleet registration**: `fleet_project` for multi-cluster features

## Stack Outputs

| Output | Description |
|---|---|
| `endpoint` | Kubernetes API server IP (the private endpoint on private-only control planes) |
| `cluster_ca_certificate` | Base64 CA certificate clients use to validate the API server's TLS certificate (public trust material) |
| `workload_identity_pool` | `PROJECT_ID.svc.id.goog` — empty when Workload Identity is disabled on a Standard cluster |
| `cluster_id` | Fully qualified resource ID: `projects/{project}/locations/{location}/clusters/{name}` |
| `name` | The cluster name as created in GCP — the handle node pools and gcloud use |
| `location` | Region (regional) or zone (zonal), exactly as provided in the spec |
| `self_link` | Server-defined URL of the cluster resource |
| `master_version` | Kubernetes version currently running on the control plane |

## Deliberately not modeled (recorded reasons)

| Excluded Feature | Why |
|---|---|
| Default node pool `node_config` / inline `node_pool` blocks | Node pools are first-class `GcpGkeNodePool` resources; the default pool is always removed. Inline pools are ForceNew on the cluster — hostile to composition. |
| `deletion_policy`, `ignore_node_count_changes`, `dataplane_optimization_mode`, `autopilot_privileged_admission`, `secret_sync_config`, `autopilot_cluster_policy_config`, `node_creation_config`, `maintenance_policy.disruption_budget`, `ip_allocation_policy.auto_ipam_config` / `network_tier_config`, pod-snapshot / slice-controller / slurm-operator / agent-sandbox addons | Exist only on the provider's unreleased main line, not on the released 6.x major the GCP modules pin. Revisit on the next provider major. |
| `tpu_config`, `pod_security_policy_config`, `cluster_telemetry`, `protect_config`, `managed_opentelemetry_config`, `workload_alts_config` | Beta-provider-only surfaces on the released line. |
| `master_auth.client_certificate_config` | Legacy client-certificate authentication — superseded by IAM and deprecated practice; the CA certificate output covers trust needs. |
| `user_managed_keys_config` | Self-managed cluster CA / etcd / service-account signing keys — a deep-compliance niche with heavy operational burden; Google-managed keys plus CMEK etcd encryption cover the real cases. |
| `enterprise_config` | Deprecated on the provider. |
| `rbac_binding_config`, `enable_k8s_beta_apis`, `gke_auto_upgrade_config`, `enable_kubernetes_alpha`, `enable_tpu`, `cluster_ipv4_cidr` (routes-based), `default_snat_status` beyond the disable toggle, `network_performance_config` beyond the tier | Niche or legacy surfaces without real-world pull; each returns on demand. |
| `logging_service` / `monitoring_service` legacy strings | Superseded by the `logging` / `monitoring` component configs. |

## Related Components

- **GcpVpcNetwork** — the network the cluster lives in
- **GcpSubnetwork** — carries the primary node range and pod/service secondary ranges
- **GcpGkeNodePool** — compute for Standard clusters (references this cluster's `name` output)
- **GcpRouterNat** — outbound internet for private nodes
- **GcpKmsKey** — CMEK key for etcd secrets encryption and NAP boot disks
- **GcpServiceAccount** — node identity for NAP-created pools
- **GcpPubSubTopic** — receives cluster lifecycle notifications
- **GcpBigQueryDataset** — receives resource-usage export records
- **GcpGkeWorkloadIdentityBinding** — binds Kubernetes service accounts to IAM service accounts on the cluster's workload pool

## Additional Resources

- [GKE Documentation](https://cloud.google.com/kubernetes-engine/docs)
- [Clusters REST API](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters)
