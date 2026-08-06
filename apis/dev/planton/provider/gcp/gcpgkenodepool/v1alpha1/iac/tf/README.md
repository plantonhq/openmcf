# GcpGkeNodePool - Terraform Module

This Terraform module provisions a Google Kubernetes Engine node pool (`google_container_node_pool`). It is the Terraform-side implementation of the Planton `GcpGkeNodePool` resource kind and has feature parity with the Pulumi module.

## Overview

The module enables the Kubernetes Engine API (`container.googleapis.com`) so a fresh project works first try, then creates the node pool with the full released-provider surface: fixed or autoscaled sizing (per-zone or total bounds, scale-to-zero), surge and blue-green upgrade strategies, compact placement and queued provisioning, per-pool pod ranges and private-nodes override, and the complete node configuration — machine shape, boot disk, local SSDs, GPUs with GKE-managed drivers, Spot/preemptible capacity, shielded and confidential VMs, CMEK boot-disk encryption, reservation affinity, secondary boot disks, kubelet tuning, and Linux OS tuning.

The pool addresses its parent by the cluster's `name` and `location` outputs (resolved to plain strings before the module runs) — no data-source lookup. An empty `project_id` falls back to the provider's default project; empty optional strings become null so the provider omits them from the API payload. The module runs on the plain `google` provider — every modeled field is GA on the released 6.x line.

For autoscaled pools, `lifecycle.ignore_changes` on `node_count` keeps plans from fighting the cluster autoscaler. GKE's `disable-legacy-endpoints=true` node metadata is enforced beneath any user metadata; platform attribution labels are merged over user resource labels.

## Usage with Planton CLI

```shell
planton tofu init --manifest ../hack/manifest.yaml
planton tofu plan --manifest ../hack/manifest.yaml
planton tofu apply --manifest ../hack/manifest.yaml --auto-approve
planton tofu destroy --manifest ../hack/manifest.yaml --auto-approve
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Terraform Usage

```bash
cd apis/dev/planton/provider/gcp/gcpgkenodepool/v1alpha1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpGkeNodePool spec | — |

The `spec` object includes: `cluster_name` + `location` (required; resolved from the cluster's outputs), `node_pool_name` (empty = metadata.name), `node_locations`, `version`, `max_pods_per_node`, `initial_node_count`, `node_count` XOR `autoscaling` (per-zone or total bounds + location policy), `management` (auto-repair/auto-upgrade, both default true), `upgrade_settings` (surge dials or blue-green rollout), `placement_policy`, `queued_provisioning_enabled`, `network_config` (pod range, private nodes, egress tier, overprovision), `node_config` (machine/disks/image/identity/scopes/labels/taints/spot/GPUs/shielded/confidential/CMEK/local SSDs/GCFS/gVNIC/fast-socket/workload-metadata/reservations/secondary disks/kubelet/Linux/logging-variant/flex-start/max-run-duration), and `project_id` (empty falls back to the provider default project).

## Outputs

| Name | Description |
|------|-------------|
| `node_pool_name` | Pool name as created in GKE |
| `instance_group_urls` | Managed instance group URLs backing the pool (one per zone) |
| `min_nodes` | Effective minimum size (autoscaling minimum, or the fixed count) |
| `max_nodes` | Effective maximum size (autoscaling maximum, or the fixed count) |
| `current_node_count` | Nodes per zone at the last deploy (autoscaler-owned at runtime) |
| `node_pool_id` | `projects/{project}/locations/{location}/clusters/{cluster}/nodePools/{name}` |
| `location` | The pool's region or zone, exactly as provided in the spec |
| `version` | Kubernetes version running on the pool's nodes |
