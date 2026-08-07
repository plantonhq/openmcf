# AWS Batch Job Queue

Deploys an AWS Batch job queue — the submission target that routes jobs onto up to three compute environments in preference order. The ordered mapping is what expresses Spot-first-with-On-Demand-overflow and zero-downtime compute replacement.

## What Gets Created

When you deploy an AwsBatchJobQueue resource, Planton provisions:

- **Job Queue** — an `aws_batch_job_queue` with the configured priority, state, compute-environment mapping, optional fair-share scheduling policy attachment, and optional stuck-job time-limit actions

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **At least one VALID compute environment** — use `AwsBatchComputeEnvironment` to provision
- **Optionally a scheduling policy** for fair-share ordering — use `AwsBatchSchedulingPolicy`

## Quick Start

Create a file `job-queue.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBatchJobQueue
metadata:
  name: etl-queue
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsBatchJobQueue.etl-queue
spec:
  region: us-west-2
  priority: 10
  computeEnvironmentOrder:
    - order: 1
      computeEnvironment:
        valueFrom:
          kind: AwsBatchComputeEnvironment
          name: etl-fargate
          fieldPath: status.outputs.compute_environment_arn
```

Deploy:

```shell
planton apply -f job-queue.yaml
```

Jobs submitted to `etl-queue` (SubmitJob or an EventBridge Batch target) now run on the referenced compute environment.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region; must match the compute environments'. | Required; non-empty |
| `computeEnvironmentOrder` | list(object) | Preference-ordered compute environment mapping. | 1-3 entries |
| `computeEnvironmentOrder[].order` | int32 | Preference position; lowest tried first. | >= 1 |
| `computeEnvironmentOrder[].computeEnvironment` | StringValueOrRef → AwsBatchComputeEnvironment | The environment's ARN. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `priority` | int32 | 0 | Scheduling preference vs other queues sharing an environment; higher wins. |
| `state` | string | `ENABLED` | `DISABLED` rejects new submissions while running jobs finish (the drain switch). |
| `schedulingPolicy` | StringValueOrRef → AwsBatchSchedulingPolicy | — | Fair-share ordering within the queue. Once set it can be replaced but never removed. |
| `jobStateTimeLimitActions` | list(object) | — | Auto-`CANCEL` jobs stuck in `RUNNABLE` past a threshold (600-86400s). `reason` selects WHICH stuck-job cause the action covers: `CAPACITY:INSUFFICIENT_INSTANCE_CAPACITY`, `MISCONFIGURATION:COMPUTE_ENVIRONMENT_MAX_RESOURCE`, or `MISCONFIGURATION:JOB_RESOURCE_REQUIREMENT`. |

All associated environments must be one family: EC2-based (`EC2`/`SPOT`) or Fargate-based (`FARGATE`/`FARGATE_SPOT`), never mixed.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `job_queue_arn` | string | The submission handle and EventBridge Batch target ARN. |
| `job_queue_name` | string | The queue's name (from `metadata.name`). |

## Presets

| Name | Description |
|------|-------------|
| [01-single-environment](../presets/01-single-environment.yaml) | One queue on one environment with stuck-job protection |
| [02-spot-overflow](../presets/02-spot-overflow.yaml) | Spot-first with On-Demand overflow — the canonical cost pattern |

## Related Components

- [AwsBatchComputeEnvironment](/docs/catalog/aws/batch-compute-environment) — the capacity this queue routes onto
- [AwsBatchJobDefinition](/docs/catalog/aws/batch-job-definition) — the container blueprint submitted to this queue
- [AwsBatchSchedulingPolicy](/docs/catalog/aws/batch-scheduling-policy) — fair-share ordering within the queue
- [AwsEventBridgeRule](/docs/catalog/aws/event-bridge-rule) — schedules or event-triggers job submissions to this queue
