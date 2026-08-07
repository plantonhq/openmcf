# AWS ECS Task Definition

Declares the immutable, versioned blueprint of an ECS workload — the containers a task runs, their images, ports, environment and secrets, health checks, sizing, volumes, and the IAM identities the task assumes. A task definition is the composition anchor of ECS compute, the way a launch template anchors EC2 fleets: it has its own lifecycle and is referenced from the places that RUN it — an ECS service for steady-state workloads, EventBridge scheduled tasks, and one-off RunTask calls. Revisions are immutable in AWS: every change registers a NEW revision of the family, so a service referencing the `task_definition_arn` output picks up each new revision on its next deployment — "change the image tag, the service rolls" falls out of the composition.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ECS Task Definition** -- a new revision of the family (named from the resource metadata), with its containers, sizing, volumes, and IAM wiring
- **CloudWatch Log Group** -- one group named `/ecs/<family>` (30-day retention unless overridden) when the default logging is left on and no existing group is referenced; each container streams under its own name
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the task definition

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Execution Role** -- an AwsIamRole trusting `ecs-tasks.amazonaws.com` with `AmazonECSTaskExecutionRolePolicy`, referenced by its `role_arn` output. Required in practice: the default CloudWatch wiring, private ECR images, and secret injection all need it.
- **Task Role** (optional) -- an AwsIamRole scoped to the AWS APIs the application itself calls, referenced by its `role_arn` output.
- **Log Group** (optional) -- an AwsCloudwatchLogGroup referenced by its name output, when several task families should share one group.

### AWS Account

- **ECS permissions** -- the credentials used by the Provider Connection must have `ecs:RegisterTaskDefinition`, `ecs:DeregisterTaskDefinition`, and `ecs:DescribeTaskDefinition`, plus `iam:PassRole` on the execution and task roles and `logs:CreateLogGroup`/`logs:PutRetentionPolicy` for the default log wiring.
- **Secret access** -- when containers inject secrets, the execution role must be able to read those Secrets Manager secrets / SSM parameters; the values never appear in the task definition.
- **Image availability** -- images must be pullable from the task's region (ECR via the execution role; other private registries via a repository-credentials secret).

## Deploy

### Console

Open the deployment store, find **AWS ECS Task Definition**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web App** preset in the [Presets](#presets) tab for a single-container Fargate service blueprint.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEcsTaskDefinition
metadata:
  name: api
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  requiresCompatibilities:
    - FARGATE
  cpu: 512
  memory: 1024
  networkMode: awsvpc
  executionRole:
    valueFrom:
      kind: AwsIamRole
      name: ecs-execution
      fieldPath: status.outputs.role_arn
  containers:
    - name: api
      image: 123456789012.dkr.ecr.us-west-2.amazonaws.com/api:1.4.2
      portMappings:
        - containerPort: 8080
          protocol: tcp
          name: http
      environment:
        LOG_LEVEL: info
      secrets:
        DB_PASSWORD: arn:aws:secretsmanager:us-west-2:123456789012:secret:db-pass
      healthCheck:
        command:
          - CMD-SHELL
          - curl -f http://localhost:8080/healthz || exit 1
```

```shell
planton apply -f task-definition.yaml
```

This registers revision 1 of the `api` family with a health-checked container, a CloudWatch log group, and a secret injected at task start. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a task definition. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Revisions, not edits** -- every change to the spec registers a NEW immutable revision; nothing mutates a live one. Consumers referencing the revisioned ARN output roll forward automatically on their next deployment.

**Fargate sizing is a pairing table** -- task-level `cpu` selects one of seven fixed sizes (256–16384 units) and constrains the valid `memory` window (e.g. 512 CPU pairs with 1–4 GiB). AWS rejects any combination outside the table at registration. On EC2, task-level sizing is optional and per-container values drive bin-packing.

**Two IAM roles by design** -- the `execution_role` is the ECS agent's setup identity (pull images, resolve secrets, write logs); the `task_role` is what the application assumes at runtime. Never merge them: a compromise of one must not grant the other's access.

**Secrets are references** -- `containers[].secrets` maps a variable name to a Secrets Manager / SSM ARN the agent resolves at task start. Plain `environment` values are visible to anyone who can describe the task definition — anything sensitive belongs in `secrets`.

**Essential vs sidecar** -- when an essential container exits, ECS stops the whole task. Leave the application essential (the default) and mark sidecars (log routers, collectors, proxies) non-essential; order startup with `depends_on` (a HEALTHY condition requires the dependency to define a health check).

**One log group per task** -- default logging creates `/ecs/<family>` and streams each container under its own name. Per-container `log_configuration` is the escape hatch to Splunk, Fluentd, or a FireLens sibling router.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Via |
|-------|-----------|-----|
| `executionRole` | AwsIamRole | `status.outputs.role_arn` |
| `taskRole` | AwsIamRole | `status.outputs.role_arn` |
| `logging.logGroup` | AwsCloudwatchLogGroup | `status.outputs.log_group_name` |
| `volumes[].efs.fileSystemId` | AwsElasticFileSystem | `status.outputs.file_system_id` |
| `volumes[].efs.accessPointId` | AwsEfsAccessPoint | `status.outputs.access_point_id` |

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `task_definition_arn` | The revisioned ARN | An AwsEcsService's `taskDefinition` reference — each new revision rolls the service |
| `arn_without_revision` | The family-level ARN | Consumers that always want the latest ACTIVE revision |
| `family` | The family name | Cross-referencing with `aws ecs describe-task-definition` |
| `revision` | The registered revision number | Deployment verification and rollback bookkeeping |
| `log_group_name` | The auto-created CloudWatch group's name | Log queries and dashboards |
| `log_group_arn` | The auto-created CloudWatch group's ARN | IAM policies scoping log access |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web application** -- one container, one HTTP port, health-checked, default CloudWatch logging. Start from the **Web App** preset.

**App with an OpenTelemetry sidecar** -- the application container plus an OTel collector sidecar, ordered with `depends_on` so telemetry is receiving before the app emits. Start from the **App with OTel Sidecar** preset.

**ARM64 worker** -- a queue-consumer on Graviton (~20% cheaper per vCPU on Fargate), no ports at all. Start from the **ARM64 Worker** preset.

## Works With

- **AwsEcsService** -- runs this blueprint as a steady-state service, referencing `task_definition_arn`; the service picks up each new revision on its next deployment.
- **AwsEcsCluster** -- the compute namespace the service places tasks into.
- **AwsIamRole** -- the execution and task roles, referenced by `executionRole` / `taskRole`.
- **AwsCloudwatchLogGroup** -- an existing shared log group, referenced by `logging.logGroup`.
- **AwsElasticFileSystem** / **AwsEfsAccessPoint** -- durable shared volumes, referenced per volume by `volumes[].efs.fileSystemId` and (recommended) `volumes[].efs.accessPointId`.
- **AwsLbTargetGroup** -- receives traffic for the container/port a service exposes from this blueprint.
