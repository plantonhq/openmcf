# AWS Batch Job Definition

Deploys an AWS Batch job definition: the versioned container blueprint jobs are submitted from — image, command, sizing, IAM identities, retries, and timeout. The compute environment provides capacity, the queue routes, and the job definition describes WHAT runs. It models single-container ECS-based jobs for EC2 and Fargate — the shape nearly every Batch workload uses. It integrates with Planton's Provider Connections for credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Batch Job Definition (a new revision)** -- registered under `metadata.name`. Revisions are immutable in AWS: every spec change registers a NEW revision rather than mutating the old one, and by default the previous revision is deregistered so exactly one stays ACTIVE
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied automatically for tracking and governance

Because the `job_definition_arn` output carries the revision, an EventBridge rule that references it picks up each new revision on its next deployment — "change the image tag, the schedule runs the new code" falls out of the composition.

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **A container image** the compute can pull: ECR images need pull permissions on the execution role (Fargate) or instance role (EC2); private non-ECR registries use a Secrets Manager credentials secret.
- **An execution role** (Fargate only, required) -- lets the agent pull the image, resolve secrets, and write logs. Reference an AwsIamRole or provide the ARN.
- **A job role** (recommended) -- the identity the workload's code runs as. Keep it minimal; it is the job's blast radius.
- **EFS file systems / access points** (only when mounting volumes) -- reference AwsElasticFileSystem / AwsEfsAccessPoint Cloud Resources.

## Deploy

### Console

Open the deployment store, find **AWS Batch Job Definition**, and click **Deploy**. The creation wizard walks you through the platform choice (which reshapes the flow — the Linux step exists only on EC2, the Fargate runtime step only on Fargate), the image, sizing with the structural Fargate pairing table, identities, configuration, storage, and the execution posture. Start from the **Fargate Container Job** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsBatchJobDefinition
metadata:
  name: etl-job
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  platformCapabilities:
    - FARGATE
  container:
    image: 123456789012.dkr.ecr.us-west-2.amazonaws.com/etl:1.4.2
    command: ["python", "run.py", "Ref::dataset"]
    vcpus: 1
    memoryMib: 2048
    executionRole:
      valueFrom:
        kind: AwsIamRole
        name: batch-execution
        fieldPath: status.outputs.role_arn
    jobRole:
      valueFrom:
        kind: AwsIamRole
        name: etl-job-role
        fieldPath: status.outputs.role_arn
  parameters:
    dataset: s3://acme-data/default.csv
  retryStrategy:
    attempts: 3
    evaluateOnExit:
      - action: RETRY
        onStatusReason: "Host EC2*"
  timeout:
    attemptDurationSeconds: 3600
```

```shell
planton apply -f batch-job-definition.yaml
```

This registers a Fargate container job with a parameterized command, the recommended Spot-reclaim retry posture (retry infrastructure failures, exit on application failures), and a one-hour per-attempt limit. A Stack Job tracks the provisioning and streams progress in real time.

### InfraChart

Jobs are submitted at runtime against a queue + this definition; within an InfraPipeline the definition typically references the IAM roles and EFS resources deployed beside it, exactly as in the CLI example's `valueFrom` blocks.

## Key Configuration

These are the most important decisions when configuring a job definition. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Platform** -- `FARGATE` requires the execution role and offers the fixed vCPU/memory pairing table, the Graviton (ARM64) runtime, and ephemeral storage up to 200 GiB; `EC2` unlocks GPUs, ulimits, privileged mode, devices, tmpfs, and swap. The definition's platform must match the queue's compute-environment family.

**Sizing** -- Fargate accepts only specific vCPU/memory pairs (0.25 vCPU → 512-2048 MiB, up to 16 vCPU → 32-120 GiB); the console offers only legal pairs. EC2 sizes are free-form — size against the compute environment's instance types or jobs sit in RUNNABLE forever.

**The two-role split** -- the EXECUTION role sets the job up (image pull, secret resolution, logs); the JOB role is what your code runs as. Deliberately never one role.

**Secrets vs environment** -- `environment` is plain text, visible to anyone who can describe the definition. `secrets` maps env names to Secrets Manager / SSM ARNs resolved by the agent at job start — the values never enter the definition.

**Storage** -- named `volumes` (EFS-backed — mount through an access point to pin POSIX identity — or EC2 host paths) placed by `mountPoints`. Pair a read-only root filesystem with one writable volume for the hardening default.

**Retry strategy** -- ordered `evaluateOnExit` conditions; first match decides RETRY or EXIT. The canonical posture: RETRY on status reason `Host EC2*` (Spot reclaims), EXIT on everything else.

**Revision posture** -- `deregisterOnNewRevision` defaults to true (exactly one ACTIVE revision). Disable it only when out-of-band SubmitJob calls pin old revisions.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AwsIamRole** (job role) | `container.jobRole` | `status.outputs.role_arn` |
| **AwsIamRole** (execution role; required on Fargate) | `container.executionRole` | `status.outputs.role_arn` |
| **AwsElasticFileSystem** (per EFS volume) | `container.volumes[].efs.fileSystemId` | `status.outputs.file_system_id` |
| **AwsEfsAccessPoint** (optional per EFS volume) | `container.volumes[].efs.accessPointId` | `status.outputs.access_point_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `job_definition_arn` | Revisioned ARN — rolls with each new revision | EventBridge Batch targets, SubmitJob calls |
| `arn_without_revision` | Revision-free ARN ("latest active") | IAM policies covering every revision |
| `job_definition_name` | Definition name | CLI commands, monitoring dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Fargate container job** -- the zero-management default: a parameterized container with the Spot-reclaim retry posture. Start from the **Fargate Container Job** preset.

**EC2 GPU job** -- GPU-reserved training or inference on EC2 compute with the NVIDIA image family. Start from the **EC2 GPU Job** preset.

**Scheduled batch** -- an EventBridge rule targeting the queue + this definition's revisioned ARN: pushing a new image tag and redeploying rolls the schedule onto the new code automatically.

## Works With

- [**AWS Batch Job Queue**](/cloud-catalog/aws-batch-job-queue) -- where jobs from this definition are submitted at runtime
- [**AWS Batch Compute Environment**](/cloud-catalog/aws-batch-compute-environment) -- the capacity the jobs run on (its family must match the platform here)
- [**AWS Batch Scheduling Policy**](/cloud-catalog/aws-batch-scheduling-policy) -- `schedulingPriority` orders this definition's jobs within their share on fair-share queues
- [**AWS IAM Role**](/cloud-catalog/aws-iam-role) -- the job and execution identities
- [**AWS Elastic File System**](/cloud-catalog/aws-elastic-file-system) -- durable shared volumes
- [**AWS EFS Access Point**](/cloud-catalog/aws-efs-access-point) -- POSIX-pinned volume mounts
