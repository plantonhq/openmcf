# AWS Batch Job Queue: Concepts

The job queue is Batch's routing layer: jobs go in, the queue decides which compute environment runs them. This reference covers the ordered-mapping semantics, queue-vs-job priority, and the operational quirks worth knowing before production.

## The Ordered Mapping

`compute_environment_order` is a preference list, not a load-balancing pool. The scheduler places every job on the LOWEST-order environment that has (or can scale to) capacity, and only falls forward when it cannot:

```yaml
computeEnvironmentOrder:
  - order: 1        # tried first: cheap, interruptible
    computeEnvironment: { valueFrom: { kind: AwsBatchComputeEnvironment, name: spot } }
  - order: 2        # overflow: reliable, pricier
    computeEnvironment: { valueFrom: { kind: AwsBatchComputeEnvironment, name: on-demand } }
```

Constraints AWS enforces (and the spec states honestly):

- **Maximum three** environments per queue.
- **One family only** — EC2-based and Fargate-based environments cannot mix on one queue, because job definitions target one platform.
- **VALID state required** — an INVALID environment cannot be associated; the compute environment's `status` output surfaces this.

## Priority Is Between Queues

`priority` matters only when multiple queues feed the SAME compute environment: the higher-priority queue's jobs are dispatched first. It says nothing about ordering within a queue — that is FIFO by default, or share-based when a scheduling policy is attached.

The two-tier pattern: a `priority: 10` interactive queue and a `priority: 1` backfill queue on one environment. Backfill jobs consume idle capacity but never delay interactive work.

## Fair-Share Attachment

Attaching an [AwsBatchSchedulingPolicy](../../awsbatchschedulingpolicy/v1alpha1/docs/README.md) switches the queue from FIFO to fair-share ordering by the share identifiers jobs carry at submission.

The one-way door: **once a queue has a scheduling policy, it can be replaced but never removed.** Going back to FIFO means recreating the queue. Decide fairness needs before the first deploy when possible.

## Stuck-Job Protection

A job whose resource requirements no associated environment can ever satisfy sits in RUNNABLE forever — and everything behind it in a FIFO queue waits. `job_state_time_limit_actions` is the fuse:

- `CANCEL` after N seconds (600-86400) removes the stuck job. It is the only action for compute-environment queues (TERMINATE exists only on the separate SageMaker Training service-environment surface).
- RUNNABLE is the only monitorable state today.
- **The `reason` field is a matcher, not a message.** Despite the name, AWS accepts only its own stuck-job cause identifiers — `CAPACITY:INSUFFICIENT_INSTANCE_CAPACITY`, `MISCONFIGURATION:COMPUTE_ENVIRONMENT_MAX_RESOURCE`, `MISCONFIGURATION:JOB_RESOURCE_REQUIREMENT` — and the action fires only on jobs stuck for that cause. Anything else fails CreateJobQueue with "Invalid job reason"; the spec pins the verified set so a manifest fails at validation instead. Add one entry per cause the fuse should cover.

## Deletion Semantics

Deleting a queue is disable-then-delete: submissions stop, running jobs finish, then the queue is removed. Both modules rely on the provider's built-in drain wait. Queues must be deleted (or repointed) before the compute environments they reference can be destroyed.

## Composition

| This kind references | Via | Output consumed |
|----------------------|-----|-----------------|
| AwsBatchComputeEnvironment | `compute_environment_order[].compute_environment` | `compute_environment_arn` |
| AwsBatchSchedulingPolicy | `scheduling_policy` | `scheduling_policy_arn` |

| Consumed by | Via | Output referenced |
|-------------|-----|-------------------|
| AwsEventBridgeRule Batch targets | target `arn` | `job_queue_arn` |
| SubmitJob callers (data plane) | name or ARN | `job_queue_name` / `job_queue_arn` |
