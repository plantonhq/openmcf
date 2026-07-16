# AWS ECS Service

Runs an ECS service: the scheduler that keeps a desired number of copies
of a referenced `AwsEcsTaskDefinition` running in a referenced
`AwsEcsCluster`, replaces unhealthy tasks, rolls deployments behind
circuit breakers, alarms, and native blue/green, and registers task IPs
into referenced `AwsLbTargetGroup` nodes.

## What Gets Created

When you deploy an AwsEcsService resource, Planton provisions:

- **ECS service** — an `aws_ecs_service` / `ecs.Service` scheduling the
  referenced task-definition revision with your capacity, networking,
  load-balancer, deployment, and Service Connect configuration
- **Autoscaling** (when configured) — an Application Auto Scaling target
  plus one target-tracking policy per metric (CPU, memory,
  requests-per-target), each with AWS-managed alarms

The cluster, task definition, subnets, security groups, and target
groups are referenced resources -- the service never creates or mutates
them.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **An ECS cluster** (`AwsEcsCluster`) and a **task definition** (`AwsEcsTaskDefinition`).
- **Subnets** (`AwsSubnet`) for awsvpc task networking — private subnets for production.
- **A target group** (`AwsLbTargetGroup`) already associated with a listener, when the service fronts traffic.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEcsService
metadata:
  name: api
spec:
  region: us-west-2
  clusterArn:
    valueFrom: { kind: AwsEcsCluster, name: prod, fieldPath: status.outputs.cluster_arn }
  taskDefinition:
    valueFrom: { kind: AwsEcsTaskDefinition, name: api, fieldPath: status.outputs.task_definition_arn }
  network:
    subnets:
      - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
      - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
```

```shell
planton apply -f service.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the cluster's and task definition's. | Required; non-empty |
| `clusterArn` | `string \| valueFrom` | The cluster the service runs in. | Required |
| `taskDefinition` | `string \| valueFrom` | The task-definition revision to run; referencing the output makes new revisions roll the service. | Required |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `desiredCount` | `int32` | `1` | Task copies to run; explicit `0` deploys the wiring with nothing running. The autoscaler owns the live count once configured. |
| `launchType` | `string` | `FARGATE` | `FARGATE`, `EC2`, or `EXTERNAL`. Mutually exclusive with `capacityProviderStrategy`. |
| `capacityProviderStrategy` | `object[]` | `[]` | Base/weight blend across `FARGATE`, `FARGATE_SPOT`, or EC2 capacity providers. |
| `network` | `object` | — | Subnets, security groups, and public-IP assignment for task ENIs (required for awsvpc task definitions). |
| `loadBalancers` | `object[]` | `[]` | Target-group registrations: which container/port registers into which `AwsLbTargetGroup`; `advancedConfiguration` adds the blue/green pair. |
| `healthCheckGracePeriodSeconds` | `int32` | `60` | Startup window during which LB health-check failures are ignored. |
| `deploymentMaximumPercent` / `deploymentMinimumHealthyPercent` | `int32` | `200` / `100` | Rolling-deployment capacity bounds. |
| `deploymentCircuitBreaker` | `object` | — | Stop failing rollouts; optionally roll back. |
| `alarms` | `object` | — | CloudWatch alarms (by reference) that fail/roll back an in-flight deployment. |
| `deploymentConfiguration` | `object` | ROLLING | `BLUE_GREEN` strategy with bake time, canary/linear shifting, and Lambda lifecycle hooks. |
| `deploymentController` | `string` | `ECS` | `ECS`, `CODE_DEPLOY`, or `EXTERNAL`. |
| `serviceConnect` | `object` | — | Service Connect mesh: namespace, exposed ports, timeouts, TLS, proxy logs. |
| `serviceRegistries` | `object` | — | Legacy Cloud Map DNS registration. |
| `volumeConfiguration` | `object` | — | A managed EBS volume per task, configured at deployment time. |
| `orderedPlacementStrategy` / `placementConstraints` | `object[]` | `[]` | EC2 task placement (spread, binpack, memberOf). |
| `availabilityZoneRebalancing` | `string` | AWS decides | `ENABLED` / `DISABLED` automatic AZ redistribution. |
| `propagateTags` / `enableEcsManagedTags` | — | — | Tag propagation to tasks for cost allocation. |
| `enableExecuteCommand` | `bool` | `false` | ECS Exec shells into running containers via SSM. |
| `autoscaling` | `object` | — | Target tracking on `cpu`, `memory`, and/or `requestsPerTarget` between `minTasks`/`maxTasks`. |

## Examples

### Cost-optimized Spot blend

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEcsService
metadata:
  name: worker
spec:
  region: us-west-2
  clusterArn:
    valueFrom: { kind: AwsEcsCluster, name: prod, fieldPath: status.outputs.cluster_arn }
  taskDefinition:
    valueFrom: { kind: AwsEcsTaskDefinition, name: worker, fieldPath: status.outputs.task_definition_arn }
  desiredCount: 4
  capacityProviderStrategy:
    - capacityProvider: FARGATE
      base: 1
      weight: 1
    - capacityProvider: FARGATE_SPOT
      weight: 4
  network:
    subnets:
      - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
      - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
```

### Request-count autoscaling behind an ALB

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEcsService
metadata:
  name: api
spec:
  region: us-west-2
  clusterArn:
    valueFrom: { kind: AwsEcsCluster, name: prod, fieldPath: status.outputs.cluster_arn }
  taskDefinition:
    valueFrom: { kind: AwsEcsTaskDefinition, name: api, fieldPath: status.outputs.task_definition_arn }
  network:
    subnets:
      - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
      - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  loadBalancers:
    - targetGroupArn:
        valueFrom: { kind: AwsLbTargetGroup, name: api, fieldPath: status.outputs.target_group_arn }
      containerName: app
      containerPort: 8080
  autoscaling:
    minTasks: 2
    maxTasks: 20
    requestsPerTarget:
      targetRequestsPerTarget: 1000
      loadBalancerArnSuffix:
        valueFrom: { kind: AwsAlb, name: main, fieldPath: status.outputs.arn_suffix }
      targetGroupArnSuffix:
        valueFrom: { kind: AwsLbTargetGroup, name: api, fieldPath: status.outputs.arn_suffix }
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `service_arn` | The service's ARN — encodes both the cluster and service names |
| `service_name` | The service's name, the ECS API's join key with the cluster |
| `cluster_arn` | The cluster ARN, republished for downstream joins |
| `task_definition_arn` | The task-definition revision this deployment runs |

## Related Components

- [AwsEcsTaskDefinition](/docs/catalog/aws/awsecstaskdefinition) — the container blueprint the service runs
- [AwsEcsCluster](/docs/catalog/aws/awsecscluster) — where the service schedules tasks
- [AwsLbTargetGroup](/docs/catalog/aws/awslbtargetgroup) — where tasks register for traffic
- [AwsLbListenerRule](/docs/catalog/aws/awslblistenerrule) — routes traffic into the target group; the blue/green swap point
- [AwsCloudwatchAlarm](/docs/catalog/aws/awscloudwatchalarm) — deployment-gating alarms
