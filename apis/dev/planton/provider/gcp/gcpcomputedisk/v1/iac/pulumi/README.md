# GcpComputeDisk - Pulumi Module

This Pulumi (Go) module provisions a zonal Compute Engine persistent disk (`compute.Disk`). It is the Pulumi-side implementation of the Planton `GcpComputeDisk` resource kind and has feature parity with the Terraform module.

## Overview

The module enables the Compute Engine API (`disable_on_destroy=false`) so a fresh project works first try and teardown never disables the API project-wide. `disk_name` falls back to `metadata.name`; user labels are merged beneath the platform attribution labels (`planton-ai_*`), identically to the Terraform module. Optional fields are set only when the spec carries them, so the provider applies its own defaults (default disk type, Google-managed encryption, `READ_WRITE_SINGLE` access).

**Immutability is the sharp edge**: name, zone, type, sources, encryption, and architecture are immutable — changing them replaces the disk and its data; size grows in place but never shrinks. At most one source (image / snapshot / source_disk) is enforced pre-deploy by the spec's CEL; none creates an empty disk. Deleting a disk still attached to a running instance fails; `create_snapshot_before_destroy` takes a final snapshot during destroy (CMEK disks reuse their key for it). For CMEK, the Compute Engine service agent must hold `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key before create.

## Usage with Planton CLI

```shell
planton pulumi up --manifest ../hack/manifest.yaml --module-dir .
planton pulumi destroy --manifest ../hack/manifest.yaml --module-dir .
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Pulumi Usage

```bash
cd apis/dev/planton/provider/gcp/gcpcomputedisk/v1/iac/pulumi
make build
pulumi up --stack dev
```

## Module Layout

- `main.go` — Pulumi entrypoint (loads the stack input, calls the module)
- `module/main.go` — provider setup + orchestration
- `module/locals.go` — disk-name fallback + label merge
- `module/disk.go` — API enablement + the disk + outputs
- `module/outputs.go` — output key constants

## Outputs

| Output | Description |
|--------|-------------|
| `name` | Name of the disk in GCP |
| `disk_id` | Server-assigned unique numeric identifier |
| `self_link` | Self-link URL — the attachment composition key |
| `zone` | Zone the disk lives in (plain zone name) |
| `size_gb` | Provisioned size in GB |
| `type` | Disk type, normalized to the plain type name (identical on both engines) |
