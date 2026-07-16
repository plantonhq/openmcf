# GcpSubnetwork - Terraform Module

This Terraform module provisions a subnetwork in a custom-mode GCP VPC. It is the Terraform-side implementation of the Planton `GcpSubnetwork` resource kind and has feature parity with the Pulumi module.

## Overview

The module enables the Compute Engine API (never disabling it on destroy) and creates a `google_compute_subnetwork` with the full address plan: primary IPv4 range, secondary (alias) ranges, purpose/role (proxy-only, PSC), dual-stack IPv6, Private Google Access (v4 and v6), the secondary-range removal latch, and VPC Flow Logs. The subnetwork resource runs on the `google-beta` provider because `allow_subnet_cidr_routes_overlap` is preview-stage on the released 6.x line; everything else is GA and identical in beta.

`name`, `project`, `region`, `network`, and `description` are immutable (ForceNew). The primary range is asymmetric: expansion updates in place, shrinkage recreates.

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
cd apis/dev/planton/provider/gcp/gcpsubnetwork/v1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpSubnetwork spec | — |

The `spec` object includes: `vpc_self_link` (required), `subnetwork_name` (required), `region` (required), `ip_cidr_range` (required except IPv6-only), `project_id` (empty falls back to the provider default project), `purpose`/`role`, `secondary_ip_ranges`, `private_ip_google_access` + `private_ipv6_google_access`, `stack_type` + `ipv6_access_type` + `external_ipv6_prefix`, `allow_subnet_cidr_routes_overlap`, `send_secondary_ip_range_if_empty`, and `log_config` (flow logs).

## Outputs

| Name | Description |
|------|-------------|
| `subnetwork_self_link` | Self-link — the value subnet consumers reference |
| `subnetwork_name` | Name in GCP |
| `region` / `ip_cidr_range` | Placement and primary range |
| `secondary_ranges` | Names + CIDRs of secondary ranges |
| `gateway_address` | IPv4 address of the default gateway |
| `subnetwork_id` | Server-assigned numeric ID (as a string) |
| `internal_ipv6_prefix` / `external_ipv6_prefix` | Allocated IPv6 prefixes |
