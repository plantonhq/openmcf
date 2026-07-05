# AWS ECS Task Definition

Registers an ECS task definition: the immutable, versioned blueprint of
the containers a task runs -- images, ports, environment, secrets, health
checks, sizing, volumes, and the IAM identities the task assumes. The
composition anchor of ECS compute: an `AwsEcsService` references it by
its revision-carrying ARN, so registering a new revision (a new image
tag) rolls the service through the resource graph.

## What Gets Created

When you deploy an AwsEcsTaskDefinition resource, Planton provisions:

- **Task definition revision** — an `aws_ecs_task_definition` /
  `ecs.TaskDefinition` under the family named by `metadata.name`; every
  spec change registers the next immutable revision
- **CloudWatch log group** (by default) — `/ecs/<family>` with 30-day
  retention; every container without its own log configuration streams
  there under its own name prefix

The IAM roles the task assumes are never modified: attach the required
policies on the referenced `AwsIamRole` nodes themselves.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **An execution role** (`AwsIamRole`) whenever the task pulls private ECR images, injects `secrets`, or writes CloudWatch logs.
- **A task role** (`AwsIamRole`) whenever the application itself calls AWS APIs.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEcsTaskDefinition
metadata:
  name: api
spec:
  region: us-west-2
  cpu: 256
  memory: 512
  containers:
    - name: app
      image: public.ecr.aws/nginx/nginx:stable
      portMappings:
        - containerPort: 80
          name: http
```

```shell
planton apply -f task-definition.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the service's and cluster's. | Required; non-empty |
| `containers` | `object[]` | The containers the task runs; sidecars are additional entries ordered by `dependsOn`. | ≥1; at least one essential |
| `containers[].name` | `string` | Unique name; the service's `containerName` join key and the log stream prefix. | Required |
| `containers[].image` | `string` | Full image reference (`repo:tag` or `repo@digest`). | Required |
| `cpu` / `memory` | `int32` | Task-level sizing (CPU units / MiB). | Required for Fargate; CEL-enforced |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `requiresCompatibilities` | `string[]` | `[FARGATE]` | Launch types validated at registration: `FARGATE`, `EC2`, `EXTERNAL`. |
| `networkMode` | `string` | `awsvpc` | `awsvpc` (required for Fargate), `bridge`, `host`, `none`. |
| `executionRole` | `string \| valueFrom` | — | The ECS agent's setup identity (image pulls, secrets, logs). |
| `taskRole` | `string \| valueFrom` | — | The application's runtime AWS identity. |
| `runtimePlatform.cpuArchitecture` | `string` | `X86_64` | `ARM64` runs on Graviton (~20% cheaper per vCPU). |
| `ephemeralStorageGib` | `int32` | 20 | Scratch storage, 21-200 GiB beyond the free 20. |
| `volumes` | `object[]` | `[]` | Named volumes: EFS (durable, Fargate-supported) or host path (EC2). |
| `logging` | `object` | auto | Task-level CloudWatch default: auto-created `/ecs/<family>` group, or a referenced `AwsCloudwatchLogGroup`; `disabled: true` turns the default off. |
| `skipDestroy` | `bool` | `false` | Keep old revisions registered on destroy. |
| `containers[].essential` | `bool` | `true` | An essential container's exit stops the task; sidecars set `false`. |
| `containers[].environment` / `secrets` | `map` | `{}` | Plain values, and Secrets Manager/SSM ARNs the agent resolves at task start. |
| `containers[].healthCheck` | `object` | — | In-container probe; `dependsOn: HEALTHY` gates siblings on it. |
| `containers[].dependsOn` | `object[]` | `[]` | Startup ordering: `START`, `HEALTHY`, `COMPLETE`, `SUCCESS`. |
| `containers[].restartPolicy` | `object` | — | In-place restarts for a crashing container without cycling the task. |
| `containers[].firelensConfiguration` | `object` | — | Mark the container as a FireLens log router (fluentbit/fluentd). |

## Examples

### Application with secrets and roles

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEcsTaskDefinition
metadata:
  name: api
spec:
  region: us-west-2
  cpu: 512
  memory: 1024
  executionRole:
    valueFrom: { kind: AwsIamRole, name: api-execution, fieldPath: status.outputs.role_arn }
  taskRole:
    valueFrom: { kind: AwsIamRole, name: api-task, fieldPath: status.outputs.role_arn }
  containers:
    - name: app
      image: 123456789012.dkr.ecr.us-west-2.amazonaws.com/api:1.4.2
      portMappings:
        - containerPort: 8080
          name: http
          appProtocol: http
      secrets:
        DATABASE_URL: arn:aws:secretsmanager:us-west-2:123456789012:secret:api-db
      healthCheck:
        command: ["CMD-SHELL", "curl -f http://localhost:8080/healthz || exit 1"]
```

### Graviton worker with an EFS volume

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEcsTaskDefinition
metadata:
  name: media-worker
spec:
  region: us-west-2
  cpu: 1024
  memory: 2048
  runtimePlatform:
    cpuArchitecture: ARM64
  volumes:
    - name: media
      efs:
        fileSystemId: fs-0123456789abcdef0
        accessPointId: fsap-0123456789abcdef0
        iamAuthorization: true
  containers:
    - name: worker
      image: 123456789012.dkr.ecr.us-west-2.amazonaws.com/media-worker:2.1.0
      mountPoints:
        - sourceVolume: media
          containerPath: /var/media
      stopTimeoutSeconds: 120
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `task_definition_arn` | The revision-carrying ARN (family:revision) services reference — a new revision changes it and rolls the referencing service |
| `arn_without_revision` | The family ARN for consumers tracking the latest ACTIVE revision |
| `family` | The family name revisions register under |
| `revision` | The revision number this deployment registered |
| `log_group_name` | The CloudWatch log group the containers log to |
| `log_group_arn` | The ARN of the auto-created log group |

## Related Components

- [AwsEcsService](/docs/catalog/aws/awsecsservice) — schedules this task definition into a cluster
- [AwsEcsCluster](/docs/catalog/aws/awsecscluster) — where the tasks run
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — the execution and task identities
- [AwsCloudwatchLogGroup](/docs/catalog/aws/awscloudwatchloggroup) — an existing log destination to reference instead of the auto-created group
