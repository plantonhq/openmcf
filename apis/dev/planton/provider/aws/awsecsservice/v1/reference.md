# AwsEcsService

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsEcsServiceSpec defines an ECS service: the scheduler that keeps a
desired number of copies of a task definition running in a cluster,
replaces unhealthy tasks, rolls deployments, and wires tasks into load
balancers and service discovery.

The service composes onto its neighbors instead of bundling them: the
task definition (WHAT runs) is a referenced AwsEcsTaskDefinition, the
cluster (WHERE it runs) is a referenced AwsEcsCluster, and traffic
arrives through referenced AwsLbTargetGroup nodes that AwsLbListener /
AwsLbListenerRule nodes route into. Because the task-definition
reference resolves to a revision-carrying ARN, registering a new
revision (a new image tag) rolls the service on its next deployment --
the deploy pipeline is the resource graph itself.

Capacity is a spectrum, not a mode: name a launch_type for the simple
case, or blend FARGATE and FARGATE_SPOT (or EC2 capacity providers) with
capacity_provider_strategy for cost-optimized fleets. Deployments are
guarded in layers -- circuit breaker (stop a failing rollout), CloudWatch
alarms (roll back on regression), and native blue/green with bake time
and lifecycle hooks for the most sensitive services.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEcsService
metadata:
  name: awsecsservice-demo
spec:
  region: us-west-2
  clusterArn:
    value: arn:aws:ecs:us-west-2:123456789012:cluster/awsecscluster-demo
  taskDefinition:
    value: arn:aws:ecs:us-west-2:123456789012:task-definition/awsecstaskdefinition-demo:1
  desiredCount: 1
  network:
    subnets:
      - value: subnet-0123456789abcdef0
      - value: subnet-0fedcba9876543210
  deploymentCircuitBreaker:
    enable: true
    rollback: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.clusterArn` | `string \| valueFrom` | yes |  | AwsEcsCluster (`status.outputs.cluster_arn`) |
| `spec.taskDefinition` | `string \| valueFrom` | yes |  | AwsEcsTaskDefinition (`status.outputs.task_definition_arn`) |
| `spec.desiredCount` | `int32` |  | `1` |  |
| `spec.launchType` | `string` |  | `FARGATE` |  |
| `spec.capacityProviderStrategy` | `[]AwsEcsServiceCapacityProviderStrategy` |  |  |  |
| `spec.capacityProviderStrategy[].capacityProvider` | `string` | yes |  |  |
| `spec.capacityProviderStrategy[].base` | `int32` |  |  |  |
| `spec.capacityProviderStrategy[].weight` | `int32` |  |  |  |
| `spec.platformVersion` | `string` |  |  |  |
| `spec.schedulingStrategy` | `string` |  |  |  |
| `spec.network` | `AwsEcsServiceNetwork` |  |  |  |
| `spec.network.subnets` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.network.securityGroups` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.network.assignPublicIp` | `bool` |  |  |  |
| `spec.loadBalancers` | `[]AwsEcsServiceLoadBalancer` |  |  |  |
| `spec.loadBalancers[].targetGroupArn` | `string \| valueFrom` | yes |  | AwsLbTargetGroup (`status.outputs.target_group_arn`) |
| `spec.loadBalancers[].containerName` | `string` | yes |  |  |
| `spec.loadBalancers[].containerPort` | `int32` |  |  |  |
| `spec.loadBalancers[].advancedConfiguration` | `AwsEcsServiceLoadBalancerAdvancedConfiguration` |  |  |  |
| `spec.loadBalancers[].advancedConfiguration.alternateTargetGroupArn` | `string \| valueFrom` | yes |  | AwsLbTargetGroup (`status.outputs.target_group_arn`) |
| `spec.loadBalancers[].advancedConfiguration.productionListenerRule` | `string \| valueFrom` | yes |  | AwsLbListenerRule (`status.outputs.rule_arn`) |
| `spec.loadBalancers[].advancedConfiguration.testListenerRule` | `string \| valueFrom` |  |  | AwsLbListenerRule (`status.outputs.rule_arn`) |
| `spec.loadBalancers[].advancedConfiguration.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.healthCheckGracePeriodSeconds` | `int32` |  |  |  |
| `spec.deploymentMaximumPercent` | `int32` |  | `200` |  |
| `spec.deploymentMinimumHealthyPercent` | `int32` |  | `100` |  |
| `spec.deploymentCircuitBreaker` | `AwsEcsServiceDeploymentCircuitBreaker` |  |  |  |
| `spec.deploymentCircuitBreaker.enable` | `bool` |  |  |  |
| `spec.deploymentCircuitBreaker.rollback` | `bool` |  |  |  |
| `spec.alarms` | `AwsEcsServiceDeploymentAlarms` |  |  |  |
| `spec.alarms.alarmNames` | `[]string \| valueFrom` | yes |  | AwsCloudwatchAlarm (`status.outputs.alarm_name`) |
| `spec.alarms.enable` | `bool` |  |  |  |
| `spec.alarms.rollback` | `bool` |  |  |  |
| `spec.deploymentConfiguration` | `AwsEcsServiceDeploymentConfiguration` |  |  |  |
| `spec.deploymentConfiguration.strategy` | `string` |  |  |  |
| `spec.deploymentConfiguration.bakeTimeInMinutes` | `int32` |  |  |  |
| `spec.deploymentConfiguration.canaryConfiguration` | `AwsEcsServiceDeploymentCanary` |  |  |  |
| `spec.deploymentConfiguration.canaryConfiguration.canaryPercent` | `double` |  |  |  |
| `spec.deploymentConfiguration.canaryConfiguration.canaryBakeTimeInMinutes` | `int32` |  |  |  |
| `spec.deploymentConfiguration.linearConfiguration` | `AwsEcsServiceDeploymentLinear` |  |  |  |
| `spec.deploymentConfiguration.linearConfiguration.stepPercent` | `double` |  |  |  |
| `spec.deploymentConfiguration.linearConfiguration.stepBakeTimeInMinutes` | `int32` |  |  |  |
| `spec.deploymentConfiguration.lifecycleHooks` | `[]AwsEcsServiceDeploymentLifecycleHook` |  |  |  |
| `spec.deploymentConfiguration.lifecycleHooks[].hookTargetArn` | `string` | yes |  |  |
| `spec.deploymentConfiguration.lifecycleHooks[].roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.deploymentConfiguration.lifecycleHooks[].lifecycleStages` | `[]string` | yes |  |  |
| `spec.deploymentConfiguration.lifecycleHooks[].hookDetails` | `string` |  |  |  |
| `spec.deploymentController` | `string` |  |  |  |
| `spec.serviceConnect` | `AwsEcsServiceServiceConnect` |  |  |  |
| `spec.serviceConnect.enabled` | `bool` |  |  |  |
| `spec.serviceConnect.namespace` | `string` |  |  |  |
| `spec.serviceConnect.services` | `[]AwsEcsServiceServiceConnectService` |  |  |  |
| `spec.serviceConnect.services[].portName` | `string` | yes |  |  |
| `spec.serviceConnect.services[].discoveryName` | `string` |  |  |  |
| `spec.serviceConnect.services[].clientAlias` | `AwsEcsServiceServiceConnectClientAlias` |  |  |  |
| `spec.serviceConnect.services[].clientAlias.port` | `int32` |  |  |  |
| `spec.serviceConnect.services[].clientAlias.dnsName` | `string` |  |  |  |
| `spec.serviceConnect.services[].ingressPortOverride` | `int32` |  |  |  |
| `spec.serviceConnect.services[].timeout` | `AwsEcsServiceServiceConnectTimeout` |  |  |  |
| `spec.serviceConnect.services[].timeout.idleTimeoutSeconds` | `int32` |  |  |  |
| `spec.serviceConnect.services[].timeout.perRequestTimeoutSeconds` | `int32` |  |  |  |
| `spec.serviceConnect.services[].tls` | `AwsEcsServiceServiceConnectTls` |  |  |  |
| `spec.serviceConnect.services[].tls.awsPcaAuthorityArn` | `string` | yes |  |  |
| `spec.serviceConnect.services[].tls.kmsKey` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.serviceConnect.services[].tls.roleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.serviceConnect.logConfiguration` | `AwsEcsServiceLogConfiguration` |  |  |  |
| `spec.serviceConnect.logConfiguration.logDriver` | `string` | yes |  |  |
| `spec.serviceConnect.logConfiguration.options` | `map<string, string>` |  |  |  |
| `spec.serviceConnect.logConfiguration.secretOptions` | `map<string, string>` |  |  |  |
| `spec.serviceRegistries` | `AwsEcsServiceServiceRegistries` |  |  |  |
| `spec.serviceRegistries.registryArn` | `string` | yes |  |  |
| `spec.serviceRegistries.containerName` | `string` |  |  |  |
| `spec.serviceRegistries.containerPort` | `int32` |  |  |  |
| `spec.serviceRegistries.port` | `int32` |  |  |  |
| `spec.volumeConfiguration` | `AwsEcsServiceVolumeConfiguration` |  |  |  |
| `spec.volumeConfiguration.name` | `string` | yes |  |  |
| `spec.volumeConfiguration.managedEbsVolume` | `AwsEcsServiceManagedEbsVolume` | yes |  |  |
| `spec.volumeConfiguration.managedEbsVolume.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.volumeConfiguration.managedEbsVolume.sizeInGb` | `int32` |  |  |  |
| `spec.volumeConfiguration.managedEbsVolume.volumeType` | `string` |  |  |  |
| `spec.volumeConfiguration.managedEbsVolume.iops` | `int32` |  |  |  |
| `spec.volumeConfiguration.managedEbsVolume.throughput` | `int32` |  |  |  |
| `spec.volumeConfiguration.managedEbsVolume.encrypted` | `bool` |  |  |  |
| `spec.volumeConfiguration.managedEbsVolume.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.volumeConfiguration.managedEbsVolume.snapshotId` | `string` |  |  |  |
| `spec.volumeConfiguration.managedEbsVolume.fileSystemType` | `string` |  |  |  |
| `spec.orderedPlacementStrategy` | `[]AwsEcsServicePlacementStrategy` |  |  |  |
| `spec.orderedPlacementStrategy[].type` | `string` |  |  |  |
| `spec.orderedPlacementStrategy[].field` | `string` |  |  |  |
| `spec.placementConstraints` | `[]AwsEcsServicePlacementConstraint` |  |  |  |
| `spec.placementConstraints[].type` | `string` |  |  |  |
| `spec.placementConstraints[].expression` | `string` |  |  |  |
| `spec.availabilityZoneRebalancing` | `string` |  |  |  |
| `spec.propagateTags` | `string` |  |  |  |
| `spec.enableEcsManagedTags` | `bool` |  |  |  |
| `spec.enableExecuteCommand` | `bool` |  |  |  |
| `spec.forceDelete` | `bool` |  |  |  |
| `spec.autoscaling` | `AwsEcsServiceAutoscaling` |  |  |  |
| `spec.autoscaling.minTasks` | `int32` |  |  |  |
| `spec.autoscaling.maxTasks` | `int32` |  |  |  |
| `spec.autoscaling.cpu` | `AwsEcsServiceAutoscalingTarget` |  |  |  |
| `spec.autoscaling.cpu.targetPercent` | `int32` |  |  |  |
| `spec.autoscaling.cpu.scaleInCooldownSeconds` | `int32` |  | `300` |  |
| `spec.autoscaling.cpu.scaleOutCooldownSeconds` | `int32` |  | `60` |  |
| `spec.autoscaling.cpu.disableScaleIn` | `bool` |  |  |  |
| `spec.autoscaling.memory` | `AwsEcsServiceAutoscalingTarget` |  |  |  |
| `spec.autoscaling.memory.targetPercent` | `int32` |  |  |  |
| `spec.autoscaling.memory.scaleInCooldownSeconds` | `int32` |  | `300` |  |
| `spec.autoscaling.memory.scaleOutCooldownSeconds` | `int32` |  | `60` |  |
| `spec.autoscaling.memory.disableScaleIn` | `bool` |  |  |  |
| `spec.autoscaling.requestsPerTarget` | `AwsEcsServiceAutoscalingRequestCountTarget` |  |  |  |
| `spec.autoscaling.requestsPerTarget.targetRequestsPerTarget` | `double` |  |  |  |
| `spec.autoscaling.requestsPerTarget.loadBalancerArnSuffix` | `string \| valueFrom` | yes |  | AwsAlb (`status.outputs.arn_suffix`) |
| `spec.autoscaling.requestsPerTarget.targetGroupArnSuffix` | `string \| valueFrom` | yes |  | AwsLbTargetGroup (`status.outputs.arn_suffix`) |
| `spec.autoscaling.requestsPerTarget.scaleInCooldownSeconds` | `int32` |  | `300` |  |
| `spec.autoscaling.requestsPerTarget.scaleOutCooldownSeconds` | `int32` |  | `60` |  |
| `spec.autoscaling.requestsPerTarget.disableScaleIn` | `bool` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the service is created in. Must match the cluster's and
task definition's region.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.clusterArn

`string | valueFrom` · required

The ECS cluster the service runs in. Reference an AwsEcsCluster's
cluster_arn output or pass a literal cluster ARN.

- references: AwsEcsCluster (`status.outputs.cluster_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsEcsCluster, name: <that resource's name>, fieldPath: status.outputs.cluster_arn}} -- a bare string does not parse

### spec.taskDefinition

`string | valueFrom` · required

The task definition revision the service runs. Reference an
AwsEcsTaskDefinition's task_definition_arn output (the recommended
wiring -- each new revision changes the output and rolls the service)
or pass a literal "family:revision" / full ARN.

- references: AwsEcsTaskDefinition (`status.outputs.task_definition_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsEcsTaskDefinition, name: <that resource's name>, fieldPath: status.outputs.task_definition_arn}} -- a bare string does not parse

### spec.desiredCount

`int32` · optional (explicit presence)

The number of task copies to keep running. Explicit 0 deploys the
service with nothing running (the wiring exists; scale up later).
When autoscaling is configured, this only seeds the initial count --
both modules then leave the live count to the scaler (and to
operators) rather than fighting it on every apply.

- default: `1`
- rule: {"int32":{"gte":0}}

### spec.launchType

`string`

The launch type: "FARGATE" (serverless -- the default), "EC2"
(container instances you manage), or "EXTERNAL" (ECS Anywhere).
Mutually exclusive with capacity_provider_strategy: name a type OR
blend providers, not both.

- default: `FARGATE`

### spec.capacityProviderStrategy

`[]AwsEcsServiceCapacityProviderStrategy`

Blend capacity across providers instead of naming one launch type --
the cost-optimization lever. Example: FARGATE base 1 / weight 1 +
FARGATE_SPOT weight 4 keeps one guaranteed on-demand task and runs
~80% of scaled capacity on Spot. EC2 clusters list their
auto-scaling-group-backed capacity providers by name.

### spec.capacityProviderStrategy[].capacityProvider

`string` · required

The capacity provider: "FARGATE", "FARGATE_SPOT", or the name of an
EC2 capacity provider attached to the cluster.

- rule: {"required":true}

### spec.capacityProviderStrategy[].base

`int32`

Minimum number of tasks guaranteed on this provider before weights
apply. Only one entry of the strategy may set a non-zero base.

- rule: {"int32":{"lte":100000,"gte":0}}

### spec.capacityProviderStrategy[].weight

`int32`

Relative share of tasks beyond the bases. Example: weight 1 on
FARGATE + weight 4 on FARGATE_SPOT scales 1:4 on-demand:Spot.

- rule: {"int32":{"lte":1000,"gte":0}}

### spec.platformVersion

`string`

The Fargate platform version (e.g. "1.4.0" or "LATEST"). Fargate
only; leave unset to track LATEST.

### spec.schedulingStrategy

`string`

Task placement across the cluster: "REPLICA" (the default -- run
desired_count copies wherever they fit) or "DAEMON" (exactly one task
per container instance -- EC2 only, for host agents like log shippers
and monitors; desired_count and autoscaling do not apply).

### spec.network

`AwsEcsServiceNetwork`

VPC networking for the task ENIs. Required for tasks whose definition
uses "awsvpc" networking -- which is every Fargate task and the modern
EC2 posture. Omit only for EC2 bridge/host-mode task definitions.

### spec.network.subnets

`[]string | valueFrom` · required

The subnets task ENIs are placed in -- private subnets for production
services, at least two AZs for availability. Reference AwsSubnet
subnet_id outputs or pass literal subnet IDs.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.network.securityGroups

`[]string | valueFrom`

Security groups applied to each task ENI. Reference AwsSecurityGroup
security_group_id outputs or pass literal group IDs. Unset falls back
to the VPC's default security group -- not a production posture.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.network.assignPublicIp

`bool`

Assign each task ENI a public IPv4 address. Only for tasks in public
subnets that must reach the internet without a NAT gateway; keep
false for private-subnet services.

### spec.loadBalancers

`[]AwsEcsServiceLoadBalancer`

Load balancer wiring: which container/port registers into which
target group. The target group is a first-class AwsLbTargetGroup the
listener (rule) routes into -- the service only registers task IPs
there. Multiple entries register multiple container ports (e.g. an
app port behind the public listener and a metrics port behind an
internal one). AWS requires each referenced target group to already
be associated with a load balancer listener at service creation.

### spec.loadBalancers[].targetGroupArn

`string | valueFrom` · required

The target group tasks register into. Reference an AwsLbTargetGroup's
target_group_arn output (the group must use target type "ip" for
awsvpc tasks) or pass a literal ARN.

- references: AwsLbTargetGroup (`status.outputs.target_group_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLbTargetGroup, name: <that resource's name>, fieldPath: status.outputs.target_group_arn}} -- a bare string does not parse

### spec.loadBalancers[].containerName

`string` · required

The container (by task-definition container name) that receives the
traffic.

- rule: {"required":true}

### spec.loadBalancers[].containerPort

`int32`

The container port that receives the traffic.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.loadBalancers[].advancedConfiguration

`AwsEcsServiceLoadBalancerAdvancedConfiguration`

Blue/green target-group pair for the BLUE_GREEN deployment strategy:
ECS shifts the production listener rule between this entry's target
group and the alternate as deployments bake. Requires
deployment_configuration.strategy = "BLUE_GREEN".

### spec.loadBalancers[].advancedConfiguration.alternateTargetGroupArn

`string | valueFrom` · required

The second target group of the blue/green pair. Reference an
AwsLbTargetGroup's target_group_arn output.

- references: AwsLbTargetGroup (`status.outputs.target_group_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLbTargetGroup, name: <that resource's name>, fieldPath: status.outputs.target_group_arn}} -- a bare string does not parse

### spec.loadBalancers[].advancedConfiguration.productionListenerRule

`string | valueFrom` · required

The production listener rule whose forward action ECS swaps between
the two target groups. Reference an AwsLbListenerRule's rule_arn
output.

- references: AwsLbListenerRule (`status.outputs.rule_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLbListenerRule, name: <that resource's name>, fieldPath: status.outputs.rule_arn}} -- a bare string does not parse

### spec.loadBalancers[].advancedConfiguration.testListenerRule

`string | valueFrom`

An optional test listener rule pointed at the green tasks before the
production swap -- how smoke traffic reaches the new version during
the bake.

- references: AwsLbListenerRule (`status.outputs.rule_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLbListenerRule, name: <that resource's name>, fieldPath: status.outputs.rule_arn}} -- a bare string does not parse

### spec.loadBalancers[].advancedConfiguration.roleArn

`string | valueFrom` · required

The IAM role ECS assumes to modify the listener rules during the
swap. Reference an AwsIamRole's role_arn output.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.healthCheckGracePeriodSeconds

`int32` · optional (explicit presence)

Seconds ECS ignores load-balancer health-check failures after a task
starts, so slow-booting apps are not killed mid-startup. Only valid
with load_balancers (no platform default, since most services set it
only when fronting traffic). Recommended: 60-120 for typical apps.

- rule: {"int32":{"gte":0}}

### spec.deploymentMaximumPercent

`int32` · optional (explicit presence)

Upper bound on running tasks during a deployment, as a percentage of
desired_count. 200 (the AWS default) starts a full replacement set
before draining the old one; 100 forces in-place replacement (needed
when capacity is tight).

- default: `200`

### spec.deploymentMinimumHealthyPercent

`int32` · optional (explicit presence)

Lower bound on healthy running tasks during a deployment, as a
percentage of desired_count. 100 (the AWS default) never dips below
desired capacity; lower values trade headroom for faster rollouts.

- default: `100`

### spec.deploymentCircuitBreaker

`AwsEcsServiceDeploymentCircuitBreaker`

The deployment circuit breaker: stop a rollout whose tasks keep
failing to reach steady state, and optionally roll back to the last
healthy deployment. The zero-configuration deployment guard -- enable
both for every production service.

- rule: rollback requires the circuit breaker to be enabled

### spec.deploymentCircuitBreaker.enable

`bool`

Enable the circuit breaker.

### spec.deploymentCircuitBreaker.rollback

`bool`

Roll back to the last steady-state deployment when the breaker
trips, instead of leaving the service stuck in a failing rollout.

### spec.alarms

`AwsEcsServiceDeploymentAlarms`

CloudWatch alarms that gate deployments: while a deployment is in
progress, if any referenced alarm fires, ECS marks the deployment
failed and (optionally) rolls back. Catches regressions the circuit
breaker cannot see -- error rates, latency, business metrics.

- rule: rollback requires alarm-gated deployments to be enabled

### spec.alarms.alarmNames

`[]string | valueFrom` · required

The alarms to watch, by NAME (the CloudWatch alarm API keys on names,
not ARNs). Reference AwsCloudwatchAlarm alarm_name outputs or pass
literal alarm names.

- references: AwsCloudwatchAlarm (`status.outputs.alarm_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCloudwatchAlarm, name: <that resource's name>, fieldPath: status.outputs.alarm_name}} -- a bare string does not parse

### spec.alarms.enable

`bool`

Enable alarm-gated deployments.

### spec.alarms.rollback

`bool`

Roll back to the last steady-state deployment when an alarm fires
mid-deployment.

### spec.deploymentConfiguration

`AwsEcsServiceDeploymentConfiguration`

Advanced deployment behavior: the ROLLING/BLUE_GREEN strategy choice,
bake time, canary/linear traffic shifting, and deployment lifecycle
hooks. Leave unset for plain rolling deployments.

- rule: strategy must be 'ROLLING' or 'BLUE_GREEN' when set
- rule: canary_configuration and linear_configuration are mutually exclusive traffic-shifting styles
- rule: bake time, canary/linear traffic shifting, and lifecycle hooks require strategy 'BLUE_GREEN'

### spec.deploymentConfiguration.strategy

`string`

"ROLLING" or "BLUE_GREEN". BLUE_GREEN requires load_balancers entries
with advanced_configuration (the target-group pair to swap between).

### spec.deploymentConfiguration.bakeTimeInMinutes

`int32` · optional (explicit presence)

Minutes the new (green) version serves production traffic before the
old one is drained -- the window in which alarms or a manual check
can still trigger an instant rollback. 0-1440.

- rule: {"int32":{"lte":1440,"gte":0}}

### spec.deploymentConfiguration.canaryConfiguration

`AwsEcsServiceDeploymentCanary`

Canary traffic shifting: send a fixed percentage to green, bake, then
shift the rest. Mutually exclusive with linear_configuration.

### spec.deploymentConfiguration.canaryConfiguration.canaryPercent

`double`

The percentage of traffic shifted in the first step (0.1-100).

- rule: {"double":{"lte":100,"gte":0.1}}

### spec.deploymentConfiguration.canaryConfiguration.canaryBakeTimeInMinutes

`int32` · optional (explicit presence)

Minutes the canary bakes before the remaining traffic shifts. 0-1440.

- rule: {"int32":{"lte":1440,"gte":0}}

### spec.deploymentConfiguration.linearConfiguration

`AwsEcsServiceDeploymentLinear`

Linear traffic shifting: shift in equal percentage steps with a bake
between steps. Mutually exclusive with canary_configuration.

### spec.deploymentConfiguration.linearConfiguration.stepPercent

`double`

The percentage shifted per step (3-100).

- rule: {"double":{"lte":100,"gte":3}}

### spec.deploymentConfiguration.linearConfiguration.stepBakeTimeInMinutes

`int32` · optional (explicit presence)

Minutes each step bakes before the next. 0-1440.

- rule: {"int32":{"lte":1440,"gte":0}}

### spec.deploymentConfiguration.lifecycleHooks

`[]AwsEcsServiceDeploymentLifecycleHook`

Lambda hooks invoked at chosen stages of the deployment lifecycle --
run integration tests against the green stack before traffic shifts,
notify, or veto.

### spec.deploymentConfiguration.lifecycleHooks[].hookTargetArn

`string` · required

The Lambda function ECS invokes. A literal function ARN.

- rule: {"required":true}

### spec.deploymentConfiguration.lifecycleHooks[].roleArn

`string | valueFrom` · required

The IAM role ECS assumes to invoke the hook. Reference an AwsIamRole's
role_arn output.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.deploymentConfiguration.lifecycleHooks[].lifecycleStages

`[]string` · required

The deployment stages the hook fires at.

- rule: {"repeated":{"minItems":"1","unique":true,"items":{"string":{"in":["RECONCILE_SERVICE","PRE_SCALE_UP","POST_SCALE_UP","TEST_TRAFFIC_SHIFT","POST_TEST_TRAFFIC_SHIFT","PRODUCTION_TRAFFIC_SHIFT","POST_PRODUCTION_TRAFFIC_SHIFT"]}}}}

### spec.deploymentConfiguration.lifecycleHooks[].hookDetails

`string`

Free-form JSON passed to the hook invocation -- deployment context
the function needs (environment names, test-suite selectors).

### spec.deploymentController

`string`

Who orchestrates deployments: "ECS" (the default -- rolling and native
blue/green), "CODE_DEPLOY" (AWS CodeDeploy drives blue/green), or
"EXTERNAL" (a third-party controller owns deployments).

### spec.serviceConnect

`AwsEcsServiceServiceConnect`

Service Connect: ECS-managed service-to-service networking on top of
AWS Cloud Map -- sidecar-free discovery ("call http://orders") with
per-request telemetry, retries, and optional TLS between services.

- rule: declaring exposed services requires Service Connect to be enabled

### spec.serviceConnect.enabled

`bool`

Enable Service Connect for this service. A service with enabled =
true and no services entries is a CLIENT: it can call other Service
Connect services in the namespace but exposes nothing itself.

### spec.serviceConnect.namespace

`string`

The Cloud Map namespace (name or ARN) the mesh lives in. Unset falls
back to the cluster's default Service Connect namespace. A literal
value -- Planton has no Cloud Map kind yet.

### spec.serviceConnect.services

`[]AwsEcsServiceServiceConnectService`

The ports this service EXPOSES to the mesh. Each entry publishes one
named task-definition port under a discovery name siblings call.

### spec.serviceConnect.services[].portName

`string` · required

The name of a port mapping in the task definition (the port_mappings
entry must set name). This is the join key between the service and
its task definition.

- rule: {"required":true}

### spec.serviceConnect.services[].discoveryName

`string`

The discovery name published to the namespace. Default: port_name.

### spec.serviceConnect.services[].clientAlias

`AwsEcsServiceServiceConnectClientAlias`

The port and DNS name CLIENTS use to call this service.

### spec.serviceConnect.services[].clientAlias.port

`int32`

The port clients connect to (often the same as the container port).

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.serviceConnect.services[].clientAlias.dnsName

`string`

The DNS name clients use (e.g. "orders" or "orders.internal").
Default: the discovery name.

### spec.serviceConnect.services[].ingressPortOverride

`int32` · optional (explicit presence)

Override the ingress port the proxy listens on for this service --
for interposing network appliances; rarely needed.

- rule: {"int32":{"lte":65535,"gte":0}}

### spec.serviceConnect.services[].timeout

`AwsEcsServiceServiceConnectTimeout`

Proxy timeouts for calls TO this service.

### spec.serviceConnect.services[].timeout.idleTimeoutSeconds

`int32`

Seconds an idle connection stays open. 0 disables the idle timeout.

- rule: {"int32":{"gte":0}}

### spec.serviceConnect.services[].timeout.perRequestTimeoutSeconds

`int32`

Seconds a single request may take end to end. 0 disables the
per-request timeout.

- rule: {"int32":{"gte":0}}

### spec.serviceConnect.services[].tls

`AwsEcsServiceServiceConnectTls`

TLS between mesh services, with certificates issued by a Private CA.

### spec.serviceConnect.services[].tls.awsPcaAuthorityArn

`string` · required

The AWS Private Certificate Authority that issues the service's
certificates. A literal PCA ARN.

- rule: {"required":true}

### spec.serviceConnect.services[].tls.kmsKey

`string | valueFrom`

The KMS key that protects the private key material. Reference an
AwsKmsKey's key_arn output or pass a literal key ARN.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.serviceConnect.services[].tls.roleArn

`string | valueFrom`

The IAM role ECS assumes to request certificates. Reference an
AwsIamRole's role_arn output.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.serviceConnect.logConfiguration

`AwsEcsServiceLogConfiguration`

Where the Service Connect proxy's own logs go (the injected agent's
logs, not the application's).

- rule: log_driver must be one of: awslogs, awsfirelens, splunk, fluentd, gelf, syslog, journald, json-file

### spec.serviceConnect.logConfiguration.logDriver

`string` · required

The log driver: "awslogs" (CloudWatch), "awsfirelens", "splunk",
"fluentd", "gelf", "syslog", "journald", or "json-file".

- rule: {"required":true}

### spec.serviceConnect.logConfiguration.options

`map<string, string>`

Driver-specific options (e.g. awslogs-group / awslogs-region /
awslogs-stream-prefix).

### spec.serviceConnect.logConfiguration.secretOptions

`map<string, string>`

Driver options whose values come from Secrets Manager / SSM (name ->
ARN), resolved by the ECS agent at task start.

### spec.serviceRegistries

`AwsEcsServiceServiceRegistries`

Legacy AWS Cloud Map service discovery registration (DNS-based). For
new meshes prefer service_connect; use this when other consumers
already resolve the Cloud Map DNS name.

### spec.serviceRegistries.registryArn

`string` · required

The Cloud Map service to register into. A literal registry ARN --
Planton has no Cloud Map kind yet.

- rule: {"required":true}

### spec.serviceRegistries.containerName

`string`

For SRV records on bridge/host-mode tasks: the container name (from
the task definition) whose address is published.

### spec.serviceRegistries.containerPort

`int32` · optional (explicit presence)

For SRV records on bridge/host-mode tasks: the container port
published alongside container_name.

- rule: {"int32":{"lte":65536,"gte":0}}

### spec.serviceRegistries.port

`int32` · optional (explicit presence)

For SRV records on awsvpc tasks: the port published with the task IP.

- rule: {"int32":{"lte":65536,"gte":0}}

### spec.volumeConfiguration

`AwsEcsServiceVolumeConfiguration`

A per-deployment managed EBS volume attached to each task, configured
at deployment time against a volume declared (by name) in the task
definition. How ECS tasks get real block storage beyond ephemeral
scratch space.

### spec.volumeConfiguration.name

`string` · required

The volume name -- must match a volume name in the task definition.

- rule: {"required":true}

### spec.volumeConfiguration.managedEbsVolume

`AwsEcsServiceManagedEbsVolume` · required

The EBS volume ECS creates and attaches per task.

- rule: {"required":true}
- rule: set size_in_gb, or snapshot_id (whose snapshot defines the size)
- rule: volume_type must be one of: gp2, gp3, io1, io2, st1, sc1, standard
- rule: throughput only applies to gp3 volumes (125-1000 MiB/s)
- rule: iops only applies to gp3, io1, and io2 volumes
- rule: file_system_type must be 'xfs', 'ext4', 'ext3', or 'ntfs' when set

### spec.volumeConfiguration.managedEbsVolume.roleArn

`string | valueFrom` · required

The IAM role ECS assumes to create, attach, and delete the volumes
(needs the AmazonECSInfrastructureRolePolicyForVolumes managed
policy). Reference an AwsIamRole's role_arn output.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.volumeConfiguration.managedEbsVolume.sizeInGb

`int32`

Volume size in GiB. Required unless snapshot_id is set (the snapshot
then defines the minimum size).

### spec.volumeConfiguration.managedEbsVolume.volumeType

`string`

Volume type: "gp3" (the sensible default), "gp2", "io1", "io2",
"st1", "sc1", or "standard".

### spec.volumeConfiguration.managedEbsVolume.iops

`int32`

Provisioned IOPS -- required for io1/io2, optional for gp3.

### spec.volumeConfiguration.managedEbsVolume.throughput

`int32`

Throughput in MiB/s, gp3 only (125-1000).

### spec.volumeConfiguration.managedEbsVolume.encrypted

`bool` · optional (explicit presence)

Encrypt the volume at rest. AWS default: true (and account-level
EBS-encryption-by-default may enforce it regardless). Optional so an
explicit false is distinguishable from unset.

### spec.volumeConfiguration.managedEbsVolume.kmsKeyId

`string | valueFrom`

The KMS key for encryption. Reference an AwsKmsKey's key_arn output
or pass a literal key ARN; unset with encryption uses the AWS-managed
aws/ebs key.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.volumeConfiguration.managedEbsVolume.snapshotId

`string`

Create each volume from this snapshot instead of empty.

### spec.volumeConfiguration.managedEbsVolume.fileSystemType

`string`

The filesystem ECS formats the volume with: "xfs" (default), "ext4",
"ext3", or "ntfs" (Windows tasks).

### spec.orderedPlacementStrategy

`[]AwsEcsServicePlacementStrategy`

Task placement strategies, applied in order (EC2 launch type only --
Fargate places tasks itself). Example: spread across AZs, then
binpack on memory.

- rule: {"repeated":{"maxItems":"5"}}

### spec.orderedPlacementStrategy[].type

`string`

"spread" (distribute across the field's values), "binpack" (fill the
least-remaining field first -- densest packing), or "random".

- rule: {"string":{"in":["binpack","random","spread"]}}

### spec.orderedPlacementStrategy[].field

`string`

What to spread over or binpack on: "attribute:ecs.availability-zone",
"instanceId" for spread; "cpu" or "memory" for binpack. Unused for
random.

### spec.placementConstraints

`[]AwsEcsServicePlacementConstraint`

Task placement constraints (EC2 launch type only). Example: memberOf
"attribute:ecs.instance-type =~ m5.*".

- rule: {"repeated":{"maxItems":"10"}}

### spec.placementConstraints[].type

`string`

"memberOf" (instances matching a cluster query expression) or
"distinctInstance" (never co-locate two of this service's tasks).

- rule: {"string":{"in":["distinctInstance","memberOf"]}}

### spec.placementConstraints[].expression

`string`

The cluster query for memberOf, e.g.
"attribute:ecs.instance-type =~ m5.*". Unused for distinctInstance.

### spec.availabilityZoneRebalancing

`string`

Automatic redistribution of tasks across availability zones when AZs
become unbalanced (after an AZ event, tasks pile into the surviving
zones and stay there): "ENABLED" or "DISABLED". Unset lets AWS decide
-- new services default to ENABLED where supported.

### spec.propagateTags

`string`

Propagate tags to tasks from "SERVICE" or "TASK_DEFINITION" ("NONE"
disables propagation). Task-level tags are how per-task cost
allocation works.

### spec.enableEcsManagedTags

`bool`

Let ECS add its managed cluster/service tags to tasks for cost and
usage attribution.

### spec.enableExecuteCommand

`bool`

Enable ECS Exec: interactive shells into running containers through
SSM ("kubectl exec" for ECS). The task role needs the SSM messages
permissions; the cluster's execute_command_configuration governs
session audit logging.

### spec.forceDelete

`bool`

Force-delete the service even while it still has running tasks --
destroy skips the scale-to-zero-first dance. Appropriate for
ephemeral environments; leave false where a stuck deletion should be
investigated instead.

### spec.autoscaling

`AwsEcsServiceAutoscaling`

Target-tracking autoscaling of desired_count via Application Auto
Scaling. Folded into the service because the scaler's identity IS
this service (one scalable target per service); the CloudWatch alarms
the policies create are managed by AWS. Not applicable to DAEMON
services.

- rule: max_tasks must be greater than or equal to min_tasks
- rule: configure at least one tracking policy (cpu, memory, or requests_per_target)

### spec.autoscaling.minTasks

`int32`

The floor. 0 allows scale-to-zero (with a metric that can get there,
e.g. a custom metric -- CPU tracking never reaches zero tasks).

- rule: {"int32":{"gte":0}}

### spec.autoscaling.maxTasks

`int32`

The ceiling -- also the cost guardrail.

- rule: {"int32":{"gte":1}}

### spec.autoscaling.cpu

`AwsEcsServiceAutoscalingTarget`

Track average CPU utilization across tasks (e.g. target 70 scales to
hold CPU near 70%). The bread-and-butter policy for compute-bound
services.

### spec.autoscaling.cpu.targetPercent

`int32`

The utilization percentage to hold (1-100). 70-75 is the usual
production sweet spot: headroom for spikes without paying for idle.

- rule: {"int32":{"lte":100,"gte":1}}

### spec.autoscaling.cpu.scaleInCooldownSeconds

`int32` · optional (explicit presence)

Seconds to wait after a scale-in before another may follow. Longer
cooldowns damp flapping. AWS default: 300.

- default: `300`

### spec.autoscaling.cpu.scaleOutCooldownSeconds

`int32` · optional (explicit presence)

Seconds to wait after a scale-out before another may follow. Keep it
short -- under-capacity hurts more than a brief overshoot. AWS
default: 60.

- default: `60`

### spec.autoscaling.cpu.disableScaleIn

`bool`

Only ever scale out on this policy; never remove capacity. For
services where a human decides when to scale in.

### spec.autoscaling.memory

`AwsEcsServiceAutoscalingTarget`

Track average memory utilization across tasks. Note memory rarely
shrinks under reduced load -- prefer CPU or request tracking for
scale-in behavior.

### spec.autoscaling.memory.targetPercent

`int32`

The utilization percentage to hold (1-100). 70-75 is the usual
production sweet spot: headroom for spikes without paying for idle.

- rule: {"int32":{"lte":100,"gte":1}}

### spec.autoscaling.memory.scaleInCooldownSeconds

`int32` · optional (explicit presence)

Seconds to wait after a scale-in before another may follow. Longer
cooldowns damp flapping. AWS default: 300.

- default: `300`

### spec.autoscaling.memory.scaleOutCooldownSeconds

`int32` · optional (explicit presence)

Seconds to wait after a scale-out before another may follow. Keep it
short -- under-capacity hurts more than a brief overshoot. AWS
default: 60.

- default: `60`

### spec.autoscaling.memory.disableScaleIn

`bool`

Only ever scale out on this policy; never remove capacity. For
services where a human decides when to scale in.

### spec.autoscaling.requestsPerTarget

`AwsEcsServiceAutoscalingRequestCountTarget`

Track requests-per-target on the load balancer -- the most direct
signal for request-serving services, reacting before CPU climbs.

### spec.autoscaling.requestsPerTarget.targetRequestsPerTarget

`double`

Requests per target per minute to hold (e.g. 1000). Derive it from a
load test: the RPS one task handles at your latency budget, times 60.

- rule: {"double":{"gt":0}}

### spec.autoscaling.requestsPerTarget.loadBalancerArnSuffix

`string | valueFrom` · required

The load balancer's ARN suffix (e.g. "app/my-alb/50dc6c495c0c9188")
-- the CloudWatch LoadBalancer dimension. Reference an AwsAlb's
arn_suffix output.

- references: AwsAlb (`status.outputs.arn_suffix`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsAlb, name: <that resource's name>, fieldPath: status.outputs.arn_suffix}} -- a bare string does not parse

### spec.autoscaling.requestsPerTarget.targetGroupArnSuffix

`string | valueFrom` · required

The target group's ARN suffix (e.g. "targetgroup/api/943f017f100becff")
-- the CloudWatch TargetGroup dimension. Reference an AwsLbTargetGroup's
arn_suffix output; use the group this service registers into.

- references: AwsLbTargetGroup (`status.outputs.arn_suffix`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLbTargetGroup, name: <that resource's name>, fieldPath: status.outputs.arn_suffix}} -- a bare string does not parse

### spec.autoscaling.requestsPerTarget.scaleInCooldownSeconds

`int32` · optional (explicit presence)

Scale-in cooldown, seconds. AWS default: 300.

- default: `300`

### spec.autoscaling.requestsPerTarget.scaleOutCooldownSeconds

`int32` · optional (explicit presence)

Scale-out cooldown, seconds. AWS default: 60.

- default: `60`

### spec.autoscaling.requestsPerTarget.disableScaleIn

`bool`

Only ever scale out on this policy; never remove capacity.

## Validation Rules

- `launch_type_valid`: launch_type must be 'FARGATE', 'EC2', or 'EXTERNAL' when set
- `launch_type_xor_capacity_providers`: launch_type and capacity_provider_strategy are mutually exclusive -- name one launch type, or blend capacity providers, not both
- `scheduling_strategy_valid`: scheduling_strategy must be 'REPLICA' or 'DAEMON' when set
- `daemon_requires_ec2`: DAEMON scheduling requires the EC2 launch type -- Fargate has no container instances to place one task on each
- `daemon_has_no_autoscaling`: DAEMON services scale with the instance fleet -- autoscaling of desired_count does not apply
- `grace_period_requires_load_balancer`: health_check_grace_period_seconds only applies when load_balancers are configured
- `platform_version_is_fargate_only`: platform_version applies to Fargate tasks -- clear it for the EC2 or EXTERNAL launch types
- `az_rebalancing_valid`: availability_zone_rebalancing must be 'ENABLED' or 'DISABLED' when set (unset lets AWS decide)
- `propagate_tags_valid`: propagate_tags must be 'SERVICE', 'TASK_DEFINITION', or 'NONE' when set
- `deployment_controller_valid`: deployment_controller must be 'ECS', 'CODE_DEPLOY', or 'EXTERNAL' when set
- `placement_is_ec2_only`: placement strategies and constraints apply to the EC2 launch type -- Fargate places tasks itself

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsEcsService, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.serviceArn` | `string` | The ARN of the service (e.g. "arn:aws:ecs:us-west-2:123456789012: service/my-cluster/api"). The primary handle for IAM policies, audit tooling, and imports; it encodes both the cluster and service names. |
| `status.outputs.serviceName` | `string` | The service's name (metadata.name), the ECS API's join key together with the cluster. |
| `status.outputs.clusterArn` | `string` | The ARN of the cluster the service runs in -- resolved from the cluster_arn reference, republished so downstream consumers can join on it without re-resolving the reference chain. |
| `status.outputs.taskDefinitionArn` | `string` | The full task definition ARN (family:revision) this deployment of the service is running -- the resolved value of the task_definition reference at deploy time. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.clusterArn` | AwsEcsCluster | `status.outputs.cluster_arn` |
| `spec.taskDefinition` | AwsEcsTaskDefinition | `status.outputs.task_definition_arn` |
| `spec.network.subnets` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.network.securityGroups` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.loadBalancers[].targetGroupArn` | AwsLbTargetGroup | `status.outputs.target_group_arn` |
| `spec.loadBalancers[].advancedConfiguration.alternateTargetGroupArn` | AwsLbTargetGroup | `status.outputs.target_group_arn` |
| `spec.loadBalancers[].advancedConfiguration.productionListenerRule` | AwsLbListenerRule | `status.outputs.rule_arn` |
| `spec.loadBalancers[].advancedConfiguration.testListenerRule` | AwsLbListenerRule | `status.outputs.rule_arn` |
| `spec.loadBalancers[].advancedConfiguration.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.alarms.alarmNames` | AwsCloudwatchAlarm | `status.outputs.alarm_name` |
| `spec.deploymentConfiguration.lifecycleHooks[].roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.serviceConnect.services[].tls.kmsKey` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.serviceConnect.services[].tls.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.volumeConfiguration.managedEbsVolume.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.volumeConfiguration.managedEbsVolume.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.autoscaling.requestsPerTarget.loadBalancerArnSuffix` | AwsAlb | `status.outputs.arn_suffix` |
| `spec.autoscaling.requestsPerTarget.targetGroupArnSuffix` | AwsLbTargetGroup | `status.outputs.arn_suffix` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
