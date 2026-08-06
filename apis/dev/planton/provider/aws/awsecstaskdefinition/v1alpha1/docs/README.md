# AWS ECS Task Definition: The Container Blueprint

## What This Component Is

An ECS task definition is the immutable, versioned template describing
the containers a task runs. AWS registers a new REVISION of the family on
every change -- nearly every attribute of the underlying resource is
create-only -- which makes the task definition a template with its own
lifecycle, exactly like an EC2 launch template. That is why it is a
first-class node rather than a block inside the service: services,
scheduled tasks, and one-off RunTask calls all reference it, revisions
outlive any one consumer, and `skip_destroy` / latest-revision tracking
are template concerns, not service concerns.

The family name comes from `metadata.name`. The `task_definition_arn`
output carries the revision, so a referencing `AwsEcsService` sees a new
value whenever a revision registers and rolls on its next deployment --
the deploy pipeline is the resource graph.

## Structured Containers, Not JSON

The ECS API takes container definitions as one opaque JSON document. This
spec models them as structured, validated messages -- multi-container
tasks with per-container sizing, named ports (the Service Connect join
key), startup ordering, health checks, ulimits, and restart policies --
and both modules render the same deterministic document (sorted maps,
explicit `essential`) so revisions stay stable across engines and
applies.

## Secrets and Logging, Modeled Honestly

- Container `secrets` and log-driver `secretOptions` are name -> ARN
  maps: references to Secrets Manager / SSM Parameter Store the ECS agent
  resolves at task start via the execution role. They are exempted from
  the secret-by-default annotation with an explicit reason -- the values
  are ARNs, never secret material.
- The task-level logging default creates one CloudWatch log group per
  family (`/ecs/<family>`) and streams each container under its own name
  prefix. A referenced `AwsCloudwatchLogGroup` swaps in an existing
  group (which then owns retention); `disabled: true` opts out; a
  per-container `logConfiguration` overrides the default for that
  container only (the FireLens path).
- The two IAM roles stay split by design: the execution role is the
  agent's setup identity, the task role is the application's runtime
  identity. Both are referenced `AwsIamRole` nodes carrying their own
  policies -- this component never modifies a role it merely references.

## Fargate Rules, Enforced Before AWS Sees Them

- Fargate requires awsvpc networking and task-level cpu/memory; both are
  CEL rules (an empty `requiresCompatibilities` means FARGATE, and the
  rules treat it that way).
- At least one container must be essential; `memoryReservation` cannot
  exceed the hard `memory` limit; `stopTimeoutSeconds` respects the
  120-second Fargate cap; ephemeral storage honors the 21-200 GiB band.

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **App Mesh `proxy_configuration`** -- App Mesh is deprecated by AWS in
  favor of ECS Service Connect, which the service spec models fully.
- **`ipc_mode` / `pid_mode`** -- EC2-only namespace sharing for niche
  host-integration workloads.
- **Fault injection (`enable_fault_injection`)** -- chaos-testing
  surface; demand-gated.
- **Docker volume drivers and FSx/S3 volume types** -- EFS and host
  paths cover the real Fargate/EC2 storage patterns; FSx is
  Windows-oriented and demand-gated.
- **`track_latest`** -- the resource graph pins revisions through
  outputs by design; tracking out-of-band ACTIVE revisions would
  reintroduce the drift the composition eliminates.
- **Inference accelerators** -- the underlying Elastic Inference service
  is discontinued by AWS.
- **Per-task-definition custom tags** -- custom user tags are a
  platform-wide concern, not per-component scope.

## Failure Modes Worth Knowing

- The awslogs driver fails at task START (not registration) when the log
  group is missing -- both modules order the group before the task
  definition and keep group creation and reference mutually exclusive.
- A container with `dependsOn: HEALTHY` on a sibling that defines no
  health check never starts; AWS rejects some of these at registration,
  others hang at start. Give the dependency a health check.
- Fargate rejects `privileged` containers and GPU
  `resourceRequirements`; both are EC2-only and documented as such on
  the fields.
