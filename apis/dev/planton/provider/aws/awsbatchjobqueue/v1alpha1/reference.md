# AwsBatchJobQueue

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsBatchJobQueueSpec defines an AWS Batch job queue: the place jobs are
submitted to, and the mapping that decides WHICH compute environment runs
them.

A queue is the routing layer of the Batch graph. It maps onto up to three
compute environments in preference order, which is what makes the
canonical Batch cost pattern possible: try the Spot environment first and
overflow to On-Demand (or the reverse for latency-sensitive queues). It is
also how a compute environment is replaced with zero queue downtime --
associate the new environment, then drain and remove the old one.

Jobs themselves are submitted against the queue at runtime (SubmitJob)
using an AwsBatchJobDefinition; the queue holds them until a mapped
compute environment has capacity.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBatchJobQueue
metadata:
  name: test-batch-queue
  id: test-batch-queue
  org: test-org
  env: dev
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: test
    pulumi.planton.dev/project: test
    pulumi.planton.dev/stack.name: dev.AwsBatchJobQueue.test-batch-queue
spec:
  region: us-west-2
  priority: 10
  computeEnvironmentOrder:
    - order: 1
      computeEnvironment:
        value: arn:aws:batch:us-west-2:123456789012:compute-environment/test-batch
  jobStateTimeLimitActions:
    - action: CANCEL
      maxTimeSeconds: 3600
      reason: MISCONFIGURATION:JOB_RESOURCE_REQUIREMENT
      state: RUNNABLE
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.state` | `string` |  | `ENABLED` |  |
| `spec.priority` | `int32` |  |  |  |
| `spec.computeEnvironmentOrder` | `[]AwsBatchJobQueueComputeEnvironmentOrder` | yes |  |  |
| `spec.computeEnvironmentOrder[].order` | `int32` |  |  |  |
| `spec.computeEnvironmentOrder[].computeEnvironment` | `string \| valueFrom` | yes |  | AwsBatchComputeEnvironment (`status.outputs.compute_environment_arn`) |
| `spec.schedulingPolicy` | `string \| valueFrom` |  |  | AwsBatchSchedulingPolicy (`status.outputs.scheduling_policy_arn`) |
| `spec.jobStateTimeLimitActions` | `[]AwsBatchJobStateTimeLimitAction` |  |  |  |
| `spec.jobStateTimeLimitActions[].action` | `string` | yes |  |  |
| `spec.jobStateTimeLimitActions[].maxTimeSeconds` | `int32` | yes |  |  |
| `spec.jobStateTimeLimitActions[].reason` | `string` | yes |  |  |
| `spec.jobStateTimeLimitActions[].state` | `string` | yes |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the job queue is created. Must match the region of
the compute environments it maps onto.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.state

`string` · optional (explicit presence)

Whether the queue accepts newly submitted jobs. When DISABLED, new
submissions are rejected but jobs already in the queue finish normally
-- the drain switch for maintenance or decommissioning.

- default: `ENABLED`
- rule: {"string":{"in":["ENABLED","DISABLED"]}}

### spec.priority

`int32`

The queue's priority relative to OTHER queues that share a compute
environment: when two queues compete for the same capacity, the one
with the higher value is scheduled first. 0 (the default when omitted)
is the lowest priority. Priority is between queues, not between jobs --
ordering within one queue comes from submission order or the
fair-share scheduling policy.

- rule: {"int32":{"gte":0}}

### spec.computeEnvironmentOrder

`[]AwsBatchJobQueueComputeEnvironmentOrder` · required

The compute environments this queue dispatches jobs to, in preference
order (lowest `order` value is tried first). AWS allows at most three,
all of one family: EC2-based (EC2/SPOT) or Fargate-based
(FARGATE/FARGATE_SPOT), never mixed. Every referenced environment must
be in the VALID state before it can be associated.

- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"3"}}

### spec.computeEnvironmentOrder[].order

`int32`

The preference position of this environment: the scheduler tries the
LOWEST order first and moves to higher values when capacity runs out.
Values must be unique within the queue.

- rule: {"int32":{"gte":1}}

### spec.computeEnvironmentOrder[].computeEnvironment

`string | valueFrom` · required

The compute environment. Reference an AwsBatchComputeEnvironment's
compute_environment_arn output or pass a literal ARN.

- references: AwsBatchComputeEnvironment (`status.outputs.compute_environment_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsBatchComputeEnvironment, name: <that resource's name>, fieldPath: status.outputs.compute_environment_arn}} -- a bare string does not parse

### spec.schedulingPolicy

`string | valueFrom`

The fair-share scheduling policy that orders jobs WITHIN this queue by
share identifier instead of first-in-first-out. Reference an
AwsBatchSchedulingPolicy's scheduling_policy_arn output or pass a
literal ARN. AWS QUIRK: once a queue has a scheduling policy it can be
replaced with another, but never removed -- removing fair-share
scheduling requires recreating the queue.

- references: AwsBatchSchedulingPolicy (`status.outputs.scheduling_policy_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsBatchSchedulingPolicy, name: <that resource's name>, fieldPath: status.outputs.scheduling_policy_arn}} -- a bare string does not parse

### spec.jobStateTimeLimitActions

`[]AwsBatchJobStateTimeLimitAction`

Automatic actions for jobs stuck at the head of the queue in a given
state past a time threshold -- the guard against a mis-sized job
blocking the whole queue in RUNNABLE forever.

### spec.jobStateTimeLimitActions[].action

`string` · required

The action to take. "CANCEL" is the only action AWS supports for
queues mapped to compute environments (the TERMINATE action exists
only for SageMaker Training service-environment queues -- a separate
Batch surface this catalog does not model).

- rule: {"required":true,"string":{"const":"CANCEL"}}

### spec.jobStateTimeLimitActions[].maxTimeSeconds

`int32` · required

Seconds a job may remain in the monitored state before the action
fires. Range: 600-86400 (10 minutes to 24 hours).

- rule: {"required":true,"int32":{"lte":86400,"gte":600}}

### spec.jobStateTimeLimitActions[].reason

`string` · required

WHICH stuck jobs the action applies to -- despite the name, this is a
MATCHER against AWS's own stuck-in-RUNNABLE causes, not a free-text log
message (CreateJobQueue rejects anything else as "Invalid job reason"):
  CAPACITY:INSUFFICIENT_INSTANCE_CAPACITY        -- no capacity in any
      mapped compute environment could place the job.
  MISCONFIGURATION:COMPUTE_ENVIRONMENT_MAX_RESOURCE -- the job asks for
      more than the environment's max_vcpus can ever offer.
  MISCONFIGURATION:JOB_RESOURCE_REQUIREMENT      -- the job's resource
      requirements cannot be satisfied by the mapped environments'
      instance types.
Add one entry per cause you want the fuse to cover.

- rule: {"required":true,"string":{"in":["CAPACITY:INSUFFICIENT_INSTANCE_CAPACITY","MISCONFIGURATION:COMPUTE_ENVIRONMENT_MAX_RESOURCE","MISCONFIGURATION:JOB_RESOURCE_REQUIREMENT"]}}

### spec.jobStateTimeLimitActions[].state

`string` · required

The job state to monitor. "RUNNABLE" is the only state AWS supports
today -- it catches jobs the scheduler can never place (mis-sized
resource requirements, exhausted capacity).

- rule: {"required":true,"string":{"const":"RUNNABLE"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsBatchJobQueue, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.job_queue_arn` | `string` | The Amazon Resource Name (ARN) of the job queue -- the handle jobs are submitted against (SubmitJob) and the target EventBridge Batch targets point at. |
| `status.outputs.job_queue_name` | `string` | The job queue's name (derived from metadata.name). SubmitJob accepts the name or the ARN interchangeably within the same account/region. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.computeEnvironmentOrder[].computeEnvironment` | AwsBatchComputeEnvironment | `status.outputs.compute_environment_arn` |
| `spec.schedulingPolicy` | AwsBatchSchedulingPolicy | `status.outputs.scheduling_policy_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
