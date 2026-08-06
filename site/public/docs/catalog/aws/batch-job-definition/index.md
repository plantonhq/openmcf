---
title: "Batch Job Definition"
description: "Batch Job Definition deployment documentation"
icon: "package"
order: 100
componentName: "awsbatchjobdefinition"
---

# AWS Batch Job Definition

Deploys an AWS Batch job definition — the versioned container blueprint jobs are submitted from: image, command, sizing, IAM identities, retries, and timeout. Every spec change registers a new immutable revision, and the revision-carrying ARN output rolls referencing consumers automatically.

## What Gets Created

When you deploy an AwsBatchJobDefinition resource, Planton provisions:

- **Job Definition revision** — an `aws_batch_job_definition` of type `container` with the structured container properties (sizing via resource requirements, identities, environment, secrets, logging, volumes), retry strategy, and timeout

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **For Fargate jobs**: an execution role carrying `AmazonECSTaskExecutionRolePolicy` — use `AwsIamRole`
- **For private ECR images on EC2**: the compute environment's instance role needs pull permissions

## Quick Start

Create a file `job-definition.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBatchJobDefinition
metadata:
  name: nightly-etl
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsBatchJobDefinition.nightly-etl
spec:
  region: us-west-2
  container:
    image: public.ecr.aws/amazonlinux/amazonlinux:2023
    command: ["echo", "hello"]
    vcpus: 1
    memoryMib: 2048
```

Deploy:

```shell
planton apply -f job-definition.yaml
```

Submit jobs against the definition (with a queue) via SubmitJob or an EventBridge Batch target.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Region the definition is registered in. | Required; non-empty |
| `container.image` | string | Full image reference. | Required; ≤255 chars |
| `container.vcpus` | double | vCPUs (fractional Fargate sizes: 0.25-16). | > 0 |
| `container.memoryMib` | int32 | Memory hard limit in MiB. | >= 4 |

### Optional Fields — Container

| Field | Type | Description |
|-------|------|-------------|
| `command` | list(string) | CMD override; supports `Ref::<key>` placeholders resolved from `parameters`. |
| `jobRole` | StringValueOrRef → AwsIamRole | The code's runtime AWS identity. |
| `executionRole` | StringValueOrRef → AwsIamRole | The agent's setup identity. **Required for Fargate.** |
| `environment` | map(string) | Plain env vars (names must not start with `AWS_BATCH`). |
| `secrets` | map(string) | Env vars from Secrets Manager / SSM ARNs, resolved at job start. |
| `logConfiguration` | object | Driver override; without it logs land in `/aws/batch/job` automatically. |
| `volumes` + `mountPoints` | list | EFS volumes (file system + access point refs) and EC2 host paths. |
| `gpus`, `ulimits`, `linuxParameters`, `privileged`, `user`, `readonlyRootFilesystem` | — | EC2-only container controls (GPU reservation, limits, devices/tmpfs/swap, hardening). |
| `runtimePlatform`, `fargatePlatformVersion`, `assignPublicIp`, `ephemeralStorageGib` | — | Fargate-only (ARM64/Windows, platform version, public IP, 21-200 GiB scratch). |
| `repositoryCredentialsSecretArn` | string | Private non-ECR registry credentials (a Secrets Manager ARN). |

### Optional Fields — Top Level

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `platformCapabilities` | list(string) | `["EC2"]` (AWS) | `EC2` and/or `FARGATE`; decides which container knobs are legal. |
| `parameters` | map(string) | — | Default `Ref::` placeholder values; overridable per job. |
| `retryStrategy.attempts` | int32 | 1 | 1-10 attempts. |
| `retryStrategy.evaluateOnExit` | list, max 5 | — | Ordered RETRY/EXIT conditions on exit code / reason / status reason (trailing `*` globs). |
| `timeout.attemptDurationSeconds` | int32 | — | Hard wall-clock limit per attempt (≥60). |
| `schedulingPriority` | int32 | — | 0-9999; consulted only on fair-share queues. |
| `propagateTags` | bool | false | Propagate definition tags to the ECS task. |
| `deregisterOnNewRevision` | bool | true | Keep exactly one ACTIVE revision. |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `job_definition_arn` | string | Revision-carrying ARN — changes every revision, rolling referencing consumers. |
| `arn_without_revision` | string | Revisionless ARN for latest-ACTIVE consumers. |
| `job_definition_name` | string | The name revisions register under. |
| `revision` | int64 | The registered revision number. |

## Presets

| Name | Description |
|------|-------------|
| [01-fargate-container-job](presets/01-fargate-container-job.yaml) | Fargate job with execution role and Spot-safe retries |
| [02-ec2-gpu-job](presets/02-ec2-gpu-job.yaml) | EC2 GPU job with ulimits and shared memory |

## Related Components

- [AwsBatchJobQueue](/docs/catalog/aws/batch-job-queue) — where jobs from this definition are submitted
- [AwsBatchComputeEnvironment](/docs/catalog/aws/batch-compute-environment) — the capacity the jobs run on
- [AwsIamRole](/docs/catalog/aws/iam-role) — the job and execution identities
- [AwsElasticFileSystem](/docs/catalog/aws/elastic-file-system) / [AwsEfsAccessPoint](/docs/catalog/aws/efs-access-point) — durable shared volumes
- [AwsEventBridgeRule](/docs/catalog/aws/event-bridge-rule) — schedules submissions referencing this definition's revision-carrying ARN
