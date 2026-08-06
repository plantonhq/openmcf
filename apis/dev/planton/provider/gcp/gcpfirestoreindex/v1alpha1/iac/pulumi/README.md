# GcpFirestoreIndex - Pulumi Module

This Pulumi (Go) module provisions a Firestore composite index (`firestore.Index`). It is the Pulumi-side implementation of the Planton `GcpFirestoreIndex` resource kind and has feature parity with the Terraform module.

## Overview

Every index property is immutable — changing anything replaces the index (Firestore rebuilds in the background). Fields map to order, array_config, or vector_config (with flat index). The module enables the Firestore API so a fresh project works first try.

## Usage with Planton CLI

```shell
planton pulumi up --manifest ../hack/manifest.yaml --module-dir .
planton pulumi destroy --manifest ../hack/manifest.yaml --module-dir .
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Pulumi Usage

```bash
cd apis/dev/planton/provider/gcp/gcpfirestoreindex/v1alpha1/iac/pulumi
make build
pulumi up --stack dev
```

## Module Layout

- `main.go` — entrypoint; loads the stack input and calls the module
- `module/main.go` — provider setup and resource orchestration
- `module/locals.go` — resolved stack input values
- `module/index.go` — API enablement + the composite index
- `module/outputs.go` — stack output keys (must match `stack_outputs.proto`)

## Outputs

| Name | Description |
|------|-------------|
| `index_id` | Fully qualified index resource name |
| `collection` | Collection (group) the index applies to |
