# GcpHealthCheck - Terraform Module

This Terraform module provisions a GCP Compute Engine health check. It is the Terraform-side implementation of the Planton `GcpHealthCheck` resource kind and has feature parity with the Pulumi module.

## Overview

The module creates exactly one of `google_compute_health_check` (global, when `spec.region` is empty) or `google_compute_region_health_check` (regional, when it is set) — GCP models the two scopes as separate API collections with an identical probe surface. Both resources run on the `google-beta` provider because the gRPC-with-TLS protocol block is preview-stage on the released 6.x line; everything else is GA and identical in beta.

`name` and `project` are immutable (ForceNew); all probe knobs (cadence, thresholds, protocol settings) update in place. Ports left unset fall through to the API's protocol defaults (http/tcp 80, https/http2/ssl 443).

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
cd apis/dev/planton/provider/gcp/gcphealthcheck/v1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpHealthCheck spec | — |

The `spec` object includes: exactly one protocol block (`http`/`https`/`http2`/`tcp`/`ssl`/`grpc`/`grpc_tls`), `project_id` (empty falls back to the provider default project), `health_check_name` (empty defaults to `metadata.name`), `region` (empty = global), `check_interval_sec`/`timeout_sec`/`healthy_threshold`/`unhealthy_threshold` (defaults 5/5/2/2), `enable_logging`, and global-only `source_regions` (exactly 3 regions).

## Outputs

| Name | Description |
|------|-------------|
| `self_link` | Self-link URI — the value backend services reference |
| `health_check_name` | Name of the health check in GCP |
| `type` | Probe protocol GCP computed (HTTP, TCP, GRPC, ...) |
| `region` | Region of a regional check; empty for global |
