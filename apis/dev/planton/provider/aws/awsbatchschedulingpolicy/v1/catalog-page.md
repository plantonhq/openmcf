# AWS Batch Scheduling Policy

Deploys an AWS Batch fair-share scheduling policy — reusable fairness rules that divide a job queue's compute capacity across share identifiers (teams, workload classes) instead of first-in-first-out. One policy can govern many queues.

## What Gets Created

When you deploy an AwsBatchSchedulingPolicy resource, Planton provisions:

- **Scheduling Policy** — an `aws_batch_scheduling_policy` with the configured fair-share dials (compute reservation, share decay, weighted share distributions)

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **Appropriate IAM permissions** for `batch:*SchedulingPolicy*` operations

## Quick Start

Create a file `fair-share.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsBatchSchedulingPolicy
metadata:
  name: team-fair-share
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsBatchSchedulingPolicy.team-fair-share
spec:
  region: us-west-2
  shareDecaySeconds: 3600
  shareDistributions:
    - shareIdentifier: teamData
      weightFactor: 0.5
    - shareIdentifier: teamMl
      weightFactor: 1.0
```

Deploy:

```shell
planton apply -f fair-share.yaml
```

Attach the policy to an `AwsBatchJobQueue` via its `schedulingPolicy` field; jobs then submit with a `shareIdentifier` to participate.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region; attachable only to same-region queues. | Required; non-empty |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `computeReservation` | int32 | — | 0-99: % of capacity held back for shares not currently running (`(r/100)^N` scaling). |
| `shareDecaySeconds` | int32 | — | 0-604800: sliding window of past usage counted against each share. |
| `shareDistributions` | list(object), max 500 | — | Weight per share identifier. |
| `shareDistributions[].shareIdentifier` | string | — | Identifier jobs carry at submission; end with `*` for a prefix match. |
| `shareDistributions[].weightFactor` | double | 1.0 (AWS) | 0.0001-999.9999. LOWER weight = MORE capacity (0.5 gets twice a 1.0 share). |

Every dial updates in place on a live policy.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `scheduling_policy_arn` | string | What job queues reference through `schedulingPolicy`. |
| `scheduling_policy_name` | string | The policy's name (from `metadata.name`). |

## Presets

| Name | Description |
|------|-------------|
| [01-team-fair-share](presets/01-team-fair-share.yaml) | Weighted team shares with decay and new-submitter headroom |

## Related Components

- [AwsBatchJobQueue](/docs/catalog/aws/batch-job-queue) — attaches this policy for fair-share ordering (replaceable, never removable)
- [AwsBatchComputeEnvironment](/docs/catalog/aws/batch-compute-environment) — the capacity the shares divide
