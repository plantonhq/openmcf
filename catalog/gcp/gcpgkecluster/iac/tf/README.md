# GcpGkeCluster - Terraform Module

This Terraform module provisions a Google Kubernetes Engine cluster (`google_container_cluster`). It is the Terraform-side implementation of the Planton `GcpGkeCluster` resource kind and has feature parity with the Pulumi module.

## Overview

The module enables the Kubernetes Engine API (`container.googleapis.com`) so a fresh project works first try, then creates the cluster with the full released-provider surface: VPC-native IP allocation (named ranges or GKE-managed), private topology (peering- or PSC-based control planes), master authorized networks and the DNS endpoint, release channels and maintenance scheduling, node auto-provisioning, Dataplane V2 features, CMEK etcd encryption, Binary Authorization, Security Posture, per-component observability with managed Prometheus, Pub/Sub notifications, BigQuery usage export, the addon set, fleet registration, and Autopilot mode.

On Standard clusters the API-mandated default node pool is removed at create time — compute comes from `GcpGkeNodePool` resources. On Autopilot clusters the node-management fields are omitted entirely (GKE owns nodes). An empty `project_id` falls back to the provider's default project; empty optional strings become null so the provider omits them from the API payload. The module runs on the plain `google` provider — every modeled field is GA on the released 6.x line.

`deletion_protection` defaults to true, matching GCP: a destroy plan fails until the spec explicitly sets it false.

## Usage with Planton CLI

```shell
planton tofu init --manifest ../../e2e/manifest.yaml
planton tofu plan --manifest ../../e2e/manifest.yaml
planton tofu apply --manifest ../../e2e/manifest.yaml --auto-approve
planton tofu destroy --manifest ../../e2e/manifest.yaml --auto-approve
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../../e2e/manifest.yaml`.

## Direct Terraform Usage

```bash
cd catalog/gcp/gcpgkecluster/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpGkeCluster spec | — |

The `spec` object includes: `location` + `network` + `subnetwork` (required), `cluster_name` (empty = metadata.name), `description`, `node_locations`, `resource_labels`, `deletion_protection`, `ip_allocation` (named ranges xor CIDR blocks, stack type, additional pod ranges, overprovision toggle), `datapath_provider`, `default_max_pods_per_node`, the Dataplane V2 toggles (FQDN/Cilium policy, multi-networking, intranode visibility, in-transit encryption), `private_cluster`, `master_authorized_networks`, `control_plane_endpoints`, `release_channel` + `min_master_version`, `maintenance_policy`, `cluster_autoscaling` (NAP), `enable_vertical_pod_autoscaling` + `hpa_profile`, `workload_identity_enabled`, `enable_shielded_nodes`, `database_encryption`, `binary_authorization_evaluation_mode`, `security_posture`, `confidential_nodes`, `anonymous_authentication_mode`, `logging` + `monitoring`, `notification_pubsub`, `enable_cost_management`, `resource_usage_export`, `addons`, `enable_autopilot` + `allow_net_admin`, `fleet_project`, and `project_id` (empty falls back to the provider default project).

## Outputs

| Name | Description |
|------|-------------|
| `endpoint` | Kubernetes API server IP (private endpoint on private-only control planes) |
| `cluster_ca_certificate` | Base64 CA certificate — public trust material clients install as a trust anchor |
| `workload_identity_pool` | `PROJECT_ID.svc.id.goog`; empty when Workload Identity is disabled on a Standard cluster |
| `cluster_id` | `projects/{project}/locations/{location}/clusters/{name}` |
| `name` | Cluster name in GCP — the handle node pools and gcloud use |
| `location` | Region or zone, exactly as provided in the spec |
| `self_link` | Server-defined URL of the cluster resource |
| `master_version` | Kubernetes version currently on the control plane |
