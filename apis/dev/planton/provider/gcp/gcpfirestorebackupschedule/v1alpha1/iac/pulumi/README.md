# GcpFirestoreBackupSchedule - Pulumi Module

This Pulumi (Go) module provisions a Firestore backup schedule (`firestore.BackupSchedule`). It is the Pulumi-side implementation of the Planton `GcpFirestoreBackupSchedule` resource kind and has feature parity with the Terraform module.

## Overview

Recurrence (daily or weekly day) is immutable; retention updates in place. Backups already taken outlive the schedule — deleting this resource stops future backups but never deletes existing ones. The module enables the Firestore API so a fresh project works first try.

## Usage with Planton CLI

```shell
planton pulumi up --manifest ../hack/manifest.yaml --module-dir .
planton pulumi destroy --manifest ../hack/manifest.yaml --module-dir .
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Pulumi Usage

```bash
cd apis/dev/planton/provider/gcp/gcpfirestorebackupschedule/v1alpha1/iac/pulumi
make build
pulumi up --stack dev
```

## Module Layout

- `main.go` — entrypoint; loads the stack input and calls the module
- `module/main.go` — provider setup and resource orchestration
- `module/locals.go` — resolved stack input values
- `module/backup_schedule.go` — API enablement + the backup schedule
- `module/outputs.go` — stack output keys (must match `stack_outputs.proto`)

## Outputs

| Name | Description |
|------|-------------|
| `schedule_id` | Server-assigned schedule ID (last path segment) |
| `database` | Firestore database name the schedule protects |
