# AwsBatchJobQueue

An AWS Batch job queue — the submission target that routes jobs onto one or more compute environments in preference order.

## What It Is

A job queue is the routing layer of AWS Batch. Jobs are submitted to a queue (never directly to compute), and the queue's `compute_environment_order` decides which [AwsBatchComputeEnvironment](../awsbatchcomputeenvironment/v1/README.md) runs them: the scheduler tries the lowest order first and falls forward when capacity runs out.

That ordered mapping is what makes the canonical Batch patterns expressible:

- **Spot-first, On-Demand overflow** — the Spot environment at order 1, On-Demand at order 2. Jobs ride cheap capacity and only spill to On-Demand under pressure.
- **Zero-downtime environment replacement** — associate the new environment, drain the old one, remove it. The queue (and its submitters) never notice.
- **Priority tiers on shared capacity** — two queues on one environment; the higher-`priority` queue's jobs are scheduled first.

## When to Use It

Every Batch deployment needs at least one queue — a compute environment alone cannot accept jobs. Create multiple queues when workloads need different priorities, different compute (GPU vs general), or different fairness rules.

## Key Facts

- **At most three compute environments per queue**, all of one family — EC2-based (`EC2`/`SPOT`) or Fargate-based (`FARGATE`/`FARGATE_SPOT`), never mixed.
- **Environments must be VALID** before they can be associated.
- **`priority` orders queues, not jobs.** When queues share an environment, higher-priority queues win. Within one queue, ordering is FIFO — or fair-share when a [scheduling policy](../awsbatchschedulingpolicy/v1/README.md) is attached.
- **A scheduling policy can be replaced but never removed** once set (AWS quirk); removing fair-share scheduling requires recreating the queue.
- **`job_state_time_limit_actions`** auto-cancel jobs stuck at the head of the queue in RUNNABLE — the guard against one mis-sized job blocking everything behind it. The `reason` field is a MATCHER against AWS's own stuck-job causes (insufficient capacity, over-max resource, unsatisfiable requirements), not a free-text message.
- **Deletion is disable-then-delete** — the queue drains before it is removed.

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | **Yes** | AWS region (must match the compute environments'). |
| `priority` | int | No (default 0) | Scheduling preference vs other queues sharing an environment; higher wins. |
| `state` | string | No (default `ENABLED`) | `DISABLED` rejects new submissions while running jobs finish. |
| `compute_environment_order` | list (1-3) | **Yes** | `{order, computeEnvironment ref}` entries; lowest order tried first. |
| `scheduling_policy` | ref → AwsBatchSchedulingPolicy | No | Fair-share ordering within the queue. |
| `job_state_time_limit_actions` | list | No | `CANCEL` jobs stuck in RUNNABLE past 600-86400s. |

## Outputs

| Field | Description |
|-------|-------------|
| `job_queue_arn` | The submission handle (SubmitJob) and EventBridge Batch target. |
| `job_queue_name` | The queue's name (from `metadata.name`). |

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBatchJobQueue
metadata:
  name: etl-queue
  org: my-org
spec:
  region: us-west-2
  priority: 10
  computeEnvironmentOrder:
    - order: 1
      computeEnvironment:
        valueFrom:
          kind: AwsBatchComputeEnvironment
          name: etl-spot
          fieldPath: status.outputs.compute_environment_arn
    - order: 2
      computeEnvironment:
        valueFrom:
          kind: AwsBatchComputeEnvironment
          name: etl-ondemand
          fieldPath: status.outputs.compute_environment_arn
  jobStateTimeLimitActions:
    - action: CANCEL
      maxTimeSeconds: 3600
      reason: MISCONFIGURATION:JOB_RESOURCE_REQUIREMENT
      state: RUNNABLE
```

See [docs/README.md](docs/README.md) for routing semantics and the fairness model.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
