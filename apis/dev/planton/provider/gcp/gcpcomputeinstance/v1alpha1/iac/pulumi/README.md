# GcpComputeInstance - Pulumi Module

This Pulumi (Go) module provisions a Compute Engine virtual machine (`compute.Instance`). It is the Pulumi-side implementation of the Planton `GcpComputeInstance` resource kind and has feature parity with the Terraform module.

## Overview

The module enables the Compute Engine API (`disable_on_destroy=false`) so a fresh project works first try and teardown never disables the API project-wide. The instance name falls back to `metadata.name`, user labels are merged beneath the platform attribution labels (`planton-ai_*`), and SSH keys fold into the metadata `ssh-keys` key (newline-joined) — all identically to the Terraform module.

**Spot is the sharp edge**: `provisioning_model = "SPOT"` requires the API's legacy preemptible flag and forbids automatic restart — both derived in the module so the spec's single switch stays honest. The boot disk honors the exactly-one-source contract (an existing disk attaches via `Source`; image/snapshot create a fresh disk through `InitializeParams`), attached disks are pre-existing `GcpComputeDisk` resources attached by reference, and `desired_status` starts/suspends/stops the VM in place. Updates that need a stop/restart (machine type, service account, ...) are performed only when `allow_stopping_for_update` is true.

## Usage with Planton CLI

```shell
planton pulumi up --manifest ../hack/manifest.yaml --module-dir .
planton pulumi destroy --manifest ../hack/manifest.yaml --module-dir .
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Pulumi Usage

```bash
cd apis/dev/planton/provider/gcp/gcpcomputeinstance/v1alpha1/iac/pulumi
make build
pulumi up --stack dev
```

## Module Layout

- `main.go` — Pulumi entrypoint (loads the stack input, calls the module)
- `module/main.go` — provider setup + orchestration + outputs export
- `module/locals.go` — instance-name fallback, label merge
- `module/instance.go` — API enablement + the instance (boot disk, attached/scratch disks, NICs, scheduling, shielded/confidential, GPUs, reservation affinity)
- `module/outputs.go` — output key constants

## Outputs

| Output | Description |
|--------|-------------|
| `instance_name` | Name of the instance |
| `instance_id` | Server-assigned unique numeric identifier |
| `self_link` | API self link |
| `internal_ip` | Primary internal IP (first interface) |
| `external_ip` | External IP of the first interface; `""` when the VM is private |
| `status` | Current status (RUNNING, TERMINATED, ...) |
| `zone` | Zone of the instance |
| `machine_type` | Machine type |
| `cpu_platform` | CPU platform the instance landed on |
