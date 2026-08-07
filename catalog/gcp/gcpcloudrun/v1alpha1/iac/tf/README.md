# GcpCloudRun - Terraform Module

This Terraform module provisions a Cloud Run service (`google_cloud_run_v2_service`). It is the Terraform-side implementation of the Planton `GcpCloudRun` resource kind and has feature parity with the Pulumi module.

## Overview

The module enables the Cloud Run Admin API (`run.googleapis.com`) so a fresh project works first try, then creates the service with the full released-provider surface: multiple containers (the serving container plus sidecars with `depends_on` startup ordering), literal and Secret Manager environment variables, startup and liveness probes (HTTP/TCP/gRPC), all five volume types (Cloud SQL sockets, Secret Manager files, in-memory or disk scratch, GCS FUSE, NFS), per-revision and service-level scaling (including MANUAL mode), revision traffic splitting with tags, direct VPC egress or a Serverless VPC Access connector, GPU node selection, CMEK image encryption, Binary Authorization, session affinity, and custom audiences.

Public access is the additive-IAM path: when `allow_unauthenticated` is true a `google_cloud_run_v2_service_iam_member` grants `roles/run.invoker` to `allUsers`; `invoker_iam_disabled` is the org-policy alternative that switches the IAM check off instead (the spec rejects setting both). `deletion_protection` defaults to true — a destroy fails until the manifest opts out.

An empty `project_id` falls back to the provider's default project; empty optional strings become null so the provider omits them from the API payload. Enum values arrive from the spec as the API's own names and pass through untranslated. The module runs on the plain `google` provider — every modeled field is GA on the released 6.x line.

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
cd catalog/gcp/gcpcloudrun/v1alpha1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpCloudRun spec | — |

The `spec` object includes: `region` (required), `containers` (required; image, command/args, env with Secret Manager references, single serving port with `h2c` support, cpu/memory limits + `cpu_idle`/`startup_cpu_boost`, volume mounts, startup/liveness probes, `depends_on`), `volumes` (cloud_sql_instance / secret / empty_dir / gcs / nfs — exactly one source each), `service_account` (empty = Compute Engine default), `scaling` + `service_scaling`, `max_instance_request_concurrency`, `timeout_seconds`, `execution_environment`, `session_affinity`, `encryption_key`, `revision`, `vpc_access` (connector XOR network_interfaces + egress), `node_selector`/`gpu_zonal_redundancy_disabled`, `ingress`, `allow_unauthenticated` XOR `invoker_iam_disabled`, `custom_audiences`, `traffic`, `launch_stage`, `binary_authorization`, `deletion_protection` (default true), and `project_id` (empty falls back to the provider default project).

## Outputs

| Name | Description |
|------|-------------|
| `url` | Canonical serving URL of the service |
| `service_name` | Name of the service as created in GCP (the serverless-NEG handle) |
| `revision` | Latest ready revision name |
| `location` | Region the service is deployed in |
| `uid` | Server-assigned unique identifier |
| `urls` | Every URL serving this service |
