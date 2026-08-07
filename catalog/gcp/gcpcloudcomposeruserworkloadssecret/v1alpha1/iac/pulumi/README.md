# GcpCloudComposerUserWorkloadsSecret - Pulumi Module

This Pulumi (Go) module provisions a Kubernetes Secret in a Cloud Composer environment's workloads namespace (`composer.UserWorkloadsSecret`). It is the Pulumi-side implementation of the Planton `GcpCloudComposerUserWorkloadsSecret` resource kind and has feature parity with the Terraform module.

## Overview

The Secret's data updates in place; name, environment, region, and project are immutable. Values are base64-encoded secret material (the Kubernetes Secret contract). The module wraps the data map with `ToSecret`, so it is held encrypted in Pulumi state and never surfaced in stack outputs. No API enablement here: the Composer API is enabled by the environment this Secret is delivered into. An empty `project_id` falls back to the provider's default project.

## Usage with Planton CLI

```shell
planton pulumi up --manifest ../hack/manifest.yaml --module-dir .
planton pulumi destroy --manifest ../hack/manifest.yaml --module-dir .
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Pulumi Usage

```bash
cd catalog/gcp/gcpcloudcomposeruserworkloadssecret/v1alpha1/iac/pulumi
make build
pulumi up --stack dev
```

## Module Layout

- `main.go` — entrypoint; loads the stack input and calls the module
- `module/main.go` — provider setup and resource orchestration
- `module/locals.go` — resolved stack input values
- `module/user_workloads_secret.go` — the Secret resource with the `ToSecret`-wrapped data map
- `module/outputs.go` — stack output keys (must match `stack_outputs.proto`)

## Outputs

| Name | Description |
|------|-------------|
| `name` | Fully qualified resource name (`projects/{p}/locations/{r}/environments/{e}/userWorkloadsSecrets/{n}`) |
| `secret_name` | The Kubernetes Secret name — what DAGs reference |

The Secret's data is deliberately never exported.
