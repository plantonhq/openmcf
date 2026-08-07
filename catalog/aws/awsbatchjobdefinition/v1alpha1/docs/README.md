# AWS Batch Job Definition: Concepts

The job definition is the workload blueprint of AWS Batch. This reference covers revision semantics, the identity split, retry discrimination, and the deliberate modeling boundaries.

## Revision Semantics

AWS never mutates a job definition: every meaningful change (image, command, sizing, retries — anything but tags) registers the NEXT revision of the name. This kind embraces that:

- The `job_definition_arn` output carries the revision (`.../job-definition/etl:7`), so consumers that reference it by output — EventBridge Batch targets above all — roll to the new revision on their next deployment. Composition, not coordination.
- `arn_without_revision` serves the opposite contract: consumers that should always track the name's latest ACTIVE revision (AWS resolves bare names to latest-ACTIVE at submission).
- `deregister_on_new_revision` (default true) deregisters the previous revision on every update, keeping exactly ONE ACTIVE revision — the one this resource manages. Set it false only when out-of-band consumers pin old revisions.
- Destroying the resource deregisters every ACTIVE revision of the name.

## The Two-Identity Split

`execution_role` and `job_role` are deliberately separate, mirroring the ECS task-definition discipline:

- **`execution_role`** — what the AGENT uses to set the job up: pull private images, resolve `secrets`, write CloudWatch logs. Required by AWS for Fargate jobs.
- **`job_role`** — what the CODE runs as: its identity for every AWS API call the workload makes. Omit it for jobs that call no AWS APIs.

Collapsing them into one role hands the workload the agent's pull/secret permissions (and vice versa) — the spec keeps the boundary visible.

## Retry Discrimination

`retry_strategy.attempts` alone re-runs every failure identically. `evaluate_on_exit` (up to 5 ordered conditions, first match wins) discriminates:

```yaml
retryStrategy:
  attempts: 3
  evaluateOnExit:
    - action: RETRY
      onStatusReason: Host EC2*    # Spot interruption -- infrastructure's fault
    - action: EXIT
      onExitCode: "1*"             # application failure -- retrying won't help
```

This is what makes Spot compute safe for pipelines: reclaims retry automatically, real bugs fail fast. Matchers accept only a trailing `*` wildcard, and a failure matching no condition behaves like EXIT.

## Platform Gating

`platform_capabilities` decides which knobs are legal, and the spec enforces AWS's rules at validation time rather than at registration:

- **Fargate forbids**: `gpus`, `privileged`, `ulimits`, `linux_parameters` (and host-path volumes in practice — EFS is the Fargate-safe volume backing).
- **Fargate requires**: `execution_role`; fractional `vcpus` from the Fargate size table (0.25-16, each pairing with a memory range).
- **EC2 forbids**: `fargate_platform_version`, `assign_public_ip`, `ephemeral_storage_gib`, `runtime_platform`.

## Deliberate Modeling Boundaries

Recorded reasons for what this kind does NOT model:

- **`node_properties` (multinode parallel jobs)** — MPI/HPC jobs spanning instances; a distinct per-node-range container model. Long-tail; additive later as its own spec arm.
- **`ecs_properties` (multi-container jobs)** — the newer multi-container variant; single-container covers the overwhelming majority of Batch workloads today.
- **`eks_properties` (Batch-on-EKS pod jobs)** — a Kubernetes pod template model of its own; an EKS-attached compute environment remains usable with definitions registered outside this graph.
- **`enable_execute_command`** — the ECS Exec debugging arm; niche operational surface.
- **S3-files volumes** — a very new volume backing not yet in the provider's documented surface; EFS + host-path cover the real demand.
- **`instance_type` in container properties** — meaningful only for the multinode arm above.

All are additive later without breaking the existing shape.

## Composition

| This kind references | Via | Output consumed |
|----------------------|-----|-----------------|
| AwsIamRole | `container.job_role`, `container.execution_role` | `role_arn` |
| AwsElasticFileSystem | `container.volumes[].efs.file_system_id` | `file_system_id` |
| AwsEfsAccessPoint | `container.volumes[].efs.access_point_id` | `access_point_id` |

| Consumed by | Via | Output referenced |
|-------------|-----|-------------------|
| AwsEventBridgeRule Batch targets | `batch_target.job_definition` | `job_definition_arn` (revision-carrying — a new revision rolls the rule) |
| SubmitJob callers (data plane) | name or ARN | `job_definition_name` / `arn_without_revision` |
