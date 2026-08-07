# GcpGcsBucket - Pulumi Module

This Pulumi (Go) module provisions a Cloud Storage bucket (`storage.Bucket`) plus additive per-bucket IAM grants. It is the Pulumi-side implementation of the Planton `GcpGcsBucket` resource kind and has feature parity with the Terraform module.

## Overview

The module enables the Cloud Storage API (`disable_on_destroy=false`) so a fresh project works first try and teardown never disables the API project-wide. User labels are merged beneath the platform attribution labels (`planton-ai_*`), identically to the Terraform module.

**Safety is the sharp edge**: `force_destroy` honors the spec (default false — destroying a non-empty bucket fails instead of erasing data); a locked retention policy is irreversible and blocks deletion until objects pass retention; the soft-delete block is sent only when the spec sets it, so unset specs follow GCP's server-side 7-day default without a perpetual diff. Numeric lifecycle conditions ride on presence — a set `0` is sent via the provider's send-zero flags, identically to the Terraform module.

## Usage with Planton CLI

```shell
planton pulumi up --manifest ../hack/manifest.yaml --module-dir .
planton pulumi destroy --manifest ../hack/manifest.yaml --module-dir .
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Pulumi Usage

```bash
cd catalog/gcp/gcpgcsbucket/v1alpha1/iac/pulumi
make build
pulumi up --stack dev
```

## Module Layout

- `main.go` — Pulumi entrypoint (loads the stack input, calls the module)
- `module/main.go` — provider setup + orchestration
- `module/locals.go` — label merge
- `module/gcs_bucket.go` — API enablement + the bucket + additive IAM members + outputs
- `module/outputs.go` — output key constants

## Outputs

| Output | Description |
|--------|-------------|
| `bucket_id` | Bucket ID (equals the bucket name) |
| `bucket_name` | Bucket name |
| `url` | `gs://<name>` |
| `self_link` | API self link |
| `location` | Location as reported by GCS |
| `project_number` | Numeric owning project |
