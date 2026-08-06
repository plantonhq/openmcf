# AwsAutoScalingGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsAutoScalingGroupSpec defines an EC2 Auto Scaling group: the fleet
manager that keeps a set of instances launched from a launch template at
the desired size, replaces unhealthy members, spreads capacity across
subnets, and scales in response to policies and schedules.

The group is deliberately a pure ORCHESTRATOR: what to launch lives in
the referenced AwsLaunchTemplate (AMI, instance type, storage, IAM
identity, metadata posture), and where traffic comes from lives in the
referenced AwsLbTargetGroup nodes. This group decides how many, where,
and when -- capacity bounds, subnet spread, purchase-option mix, health
model, scaling behavior, and instance lifecycle.

Scaling policies, scheduled actions, lifecycle hooks, and notifications
are AWS sub-resources of exactly one group, are referenced by nothing
else, and have no meaning in isolation -- so they are folded into this
spec rather than modeled as standalone kinds. Both IaC modules manage
each entry as its own provider resource, so adding or removing one is an
in-place update, never a group replacement.

The group name comes from metadata.name (AWS limit: 255 characters).
The name is create-only in AWS; everything else updates in place.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsAutoScalingGroup
metadata:
  name: web-demo
spec:
  region: us-west-2
  subnets:
    - value: subnet-0123456789abcdef0
    - value: subnet-0123456789abcdef1
  launchTemplate:
    launchTemplateId:
      value: lt-0123456789abcdef0
  minSize: 1
  maxSize: 4
  desiredCapacity: 2
  # The web-service shape: ELB health checks against a target group,
  # rolling instance refresh with surge, and a CPU target-tracking policy.
  # Exercises the nested launch-template, refresh, and policy objects so
  # the fixture proves the full variable contract, not just the scalars.
  healthCheckType: ELB
  healthCheckGracePeriodSeconds: 120
  targetGroups:
    - value: arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/api/50dc6c495c0c9188
  terminationPolicies:
    - OldestLaunchTemplate
    - Default
  instanceRefresh:
    strategy: Rolling
    preferences:
      minHealthyPercentage: 90
      maxHealthyPercentage: 110
      autoRollback: true
  scalingPolicies:
    - name: cpu-target
      policyType: TargetTrackingScaling
      targetTracking:
        targetValue: 60
        predefinedMetricType: ASGAverageCPUUtilization
  enabledMetrics:
    - GroupMinSize
    - GroupMaxSize
    - GroupDesiredCapacity
    - GroupInServiceInstances
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.subnets` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.launchTemplate` | `AwsAutoScalingGroupLaunchTemplateRef` |  |  |  |
| `spec.launchTemplate.launchTemplateId` | `string \| valueFrom` | yes |  | AwsLaunchTemplate (`status.outputs.launch_template_id`) |
| `spec.launchTemplate.version` | `string` |  |  |  |
| `spec.mixedInstancesPolicy` | `AwsAutoScalingGroupMixedInstancesPolicy` |  |  |  |
| `spec.mixedInstancesPolicy.launchTemplate` | `AwsAutoScalingGroupLaunchTemplateRef` | yes |  |  |
| `spec.mixedInstancesPolicy.launchTemplate.launchTemplateId` | `string \| valueFrom` | yes |  | AwsLaunchTemplate (`status.outputs.launch_template_id`) |
| `spec.mixedInstancesPolicy.launchTemplate.version` | `string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides` | `[]AwsAutoScalingGroupMixedInstancesOverride` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceType` | `string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].weightedCapacity` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].launchTemplate` | `AwsAutoScalingGroupLaunchTemplateRef` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].launchTemplate.launchTemplateId` | `string \| valueFrom` | yes |  | AwsLaunchTemplate (`status.outputs.launch_template_id`) |
| `spec.mixedInstancesPolicy.overrides[].launchTemplate.version` | `string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements` | `AwsAutoScalingGroupInstanceRequirements` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryMib` | `AwsAutoScalingGroupIntRange` | yes |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryMib.min` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryMib.max` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.vcpuCount` | `AwsAutoScalingGroupIntRange` | yes |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.vcpuCount.min` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.vcpuCount.max` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.allowedInstanceTypes` | `[]string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.excludedInstanceTypes` | `[]string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.instanceGenerations` | `[]string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.cpuManufacturers` | `[]string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.bareMetal` | `string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.burstablePerformance` | `string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.requireHibernateSupport` | `bool` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.spotMaxPricePercentageOverLowestPrice` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.maxSpotPriceAsPercentageOfOptimalOnDemandPrice` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.onDemandMaxPricePercentageOverLowestPrice` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.localStorage` | `string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.localStorageTypes` | `[]string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.totalLocalStorageGb` | `AwsAutoScalingGroupDoubleRange` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.totalLocalStorageGb.min` | `double` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.totalLocalStorageGb.max` | `double` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryGibPerVcpu` | `AwsAutoScalingGroupDoubleRange` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryGibPerVcpu.min` | `double` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryGibPerVcpu.max` | `double` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkInterfaceCount` | `AwsAutoScalingGroupIntRange` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkInterfaceCount.min` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkInterfaceCount.max` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkBandwidthGbps` | `AwsAutoScalingGroupDoubleRange` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkBandwidthGbps.min` | `double` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkBandwidthGbps.max` | `double` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.baselineEbsBandwidthMbps` | `AwsAutoScalingGroupIntRange` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.baselineEbsBandwidthMbps.min` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.baselineEbsBandwidthMbps.max` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorCount` | `AwsAutoScalingGroupIntRange` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorCount.min` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorCount.max` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorManufacturers` | `[]string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorNames` | `[]string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorTypes` | `[]string` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorTotalMemoryMib` | `AwsAutoScalingGroupIntRange` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorTotalMemoryMib.min` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorTotalMemoryMib.max` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.instancesDistribution` | `AwsAutoScalingGroupInstancesDistribution` |  |  |  |
| `spec.mixedInstancesPolicy.instancesDistribution.onDemandAllocationStrategy` | `string` |  |  |  |
| `spec.mixedInstancesPolicy.instancesDistribution.onDemandBaseCapacity` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.instancesDistribution.onDemandPercentageAboveBaseCapacity` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.instancesDistribution.spotAllocationStrategy` | `string` |  |  |  |
| `spec.mixedInstancesPolicy.instancesDistribution.spotInstancePools` | `int32` |  |  |  |
| `spec.mixedInstancesPolicy.instancesDistribution.spotMaxPrice` | `string` |  |  |  |
| `spec.minSize` | `int32` |  |  |  |
| `spec.maxSize` | `int32` |  |  |  |
| `spec.desiredCapacity` | `int32` |  |  |  |
| `spec.desiredCapacityType` | `string` |  |  |  |
| `spec.capacityRebalance` | `bool` |  |  |  |
| `spec.defaultCooldownSeconds` | `int32` |  |  |  |
| `spec.defaultInstanceWarmupSeconds` | `int32` |  |  |  |
| `spec.healthCheckType` | `string` |  |  |  |
| `spec.healthCheckGracePeriodSeconds` | `int32` |  |  |  |
| `spec.targetGroups` | `[]string \| valueFrom` |  |  | AwsLbTargetGroup (`status.outputs.target_group_arn`) |
| `spec.terminationPolicies` | `[]string` |  |  |  |
| `spec.maxInstanceLifetimeSeconds` | `int32` |  |  |  |
| `spec.protectFromScaleIn` | `bool` |  |  |  |
| `spec.placementGroup` | `string` |  |  |  |
| `spec.serviceLinkedRoleArn` | `string` |  |  |  |
| `spec.enabledMetrics` | `[]string` |  |  |  |
| `spec.suspendedProcesses` | `[]string` |  |  |  |
| `spec.instanceRefresh` | `AwsAutoScalingGroupInstanceRefresh` |  |  |  |
| `spec.instanceRefresh.strategy` | `string` | yes |  |  |
| `spec.instanceRefresh.triggers` | `[]string` |  |  |  |
| `spec.instanceRefresh.preferences` | `AwsAutoScalingGroupInstanceRefreshPreferences` |  |  |  |
| `spec.instanceRefresh.preferences.minHealthyPercentage` | `int32` |  |  |  |
| `spec.instanceRefresh.preferences.maxHealthyPercentage` | `int32` |  |  |  |
| `spec.instanceRefresh.preferences.instanceWarmupSeconds` | `int32` |  |  |  |
| `spec.instanceRefresh.preferences.checkpointPercentages` | `[]int32` |  |  |  |
| `spec.instanceRefresh.preferences.checkpointDelaySeconds` | `int32` |  |  |  |
| `spec.instanceRefresh.preferences.autoRollback` | `bool` |  |  |  |
| `spec.instanceRefresh.preferences.alarms` | `[]string \| valueFrom` |  |  | AwsCloudwatchAlarm (`status.outputs.alarm_name`) |
| `spec.instanceRefresh.preferences.scaleInProtectedInstances` | `string` |  |  |  |
| `spec.instanceRefresh.preferences.standbyInstances` | `string` |  |  |  |
| `spec.instanceRefresh.preferences.skipMatching` | `bool` |  |  |  |
| `spec.warmPool` | `AwsAutoScalingGroupWarmPool` |  |  |  |
| `spec.warmPool.poolState` | `string` |  |  |  |
| `spec.warmPool.minSize` | `int32` |  |  |  |
| `spec.warmPool.maxGroupPreparedCapacity` | `int32` |  |  |  |
| `spec.warmPool.reuseOnScaleIn` | `bool` |  |  |  |
| `spec.instanceMaintenancePolicy` | `AwsAutoScalingGroupInstanceMaintenancePolicy` |  |  |  |
| `spec.instanceMaintenancePolicy.minHealthyPercentage` | `int32` |  |  |  |
| `spec.instanceMaintenancePolicy.maxHealthyPercentage` | `int32` |  |  |  |
| `spec.capacityDistributionStrategy` | `string` |  |  |  |
| `spec.forceDelete` | `bool` |  |  |  |
| `spec.waitForCapacityTimeout` | `string` |  |  |  |
| `spec.scalingPolicies` | `[]AwsAutoScalingGroupScalingPolicy` |  |  |  |
| `spec.scalingPolicies[].name` | `string` | yes |  |  |
| `spec.scalingPolicies[].policyType` | `string` | yes |  |  |
| `spec.scalingPolicies[].estimatedInstanceWarmupSeconds` | `int32` |  |  |  |
| `spec.scalingPolicies[].targetTracking` | `AwsAutoScalingGroupTargetTrackingConfig` |  |  |  |
| `spec.scalingPolicies[].targetTracking.targetValue` | `double` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.predefinedMetricType` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.resourceLabel` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric` | `AwsAutoScalingGroupCustomizedMetric` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metricName` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.namespace` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.statistic` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.unit` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.dimensions` | `[]AwsAutoScalingGroupMetricDimension` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.dimensions[].name` | `string` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.dimensions[].value` | `string` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.periodSeconds` | `int32` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics` | `[]AwsAutoScalingGroupMetricDataQuery` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].id` | `string` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].expression` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat` | `AwsAutoScalingGroupMetricStat` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.metricName` | `string` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.namespace` | `string` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.stat` | `string` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.unit` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.dimensions` | `[]AwsAutoScalingGroupMetricDimension` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.dimensions[].name` | `string` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.dimensions[].value` | `string` | yes |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.periodSeconds` | `int32` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].label` | `string` |  |  |  |
| `spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].returnData` | `bool` |  |  |  |
| `spec.scalingPolicies[].targetTracking.disableScaleIn` | `bool` |  |  |  |
| `spec.scalingPolicies[].stepScaling` | `AwsAutoScalingGroupStepScalingConfig` |  |  |  |
| `spec.scalingPolicies[].stepScaling.adjustmentType` | `string` | yes |  |  |
| `spec.scalingPolicies[].stepScaling.metricAggregationType` | `string` |  |  |  |
| `spec.scalingPolicies[].stepScaling.minAdjustmentMagnitude` | `int32` |  |  |  |
| `spec.scalingPolicies[].stepScaling.stepAdjustments` | `[]AwsAutoScalingGroupStepAdjustment` | yes |  |  |
| `spec.scalingPolicies[].stepScaling.stepAdjustments[].scalingAdjustment` | `int32` |  |  |  |
| `spec.scalingPolicies[].stepScaling.stepAdjustments[].metricIntervalLowerBound` | `string` |  |  |  |
| `spec.scalingPolicies[].stepScaling.stepAdjustments[].metricIntervalUpperBound` | `string` |  |  |  |
| `spec.scalingPolicies[].simpleScaling` | `AwsAutoScalingGroupSimpleScalingConfig` |  |  |  |
| `spec.scalingPolicies[].simpleScaling.adjustmentType` | `string` | yes |  |  |
| `spec.scalingPolicies[].simpleScaling.scalingAdjustment` | `int32` |  |  |  |
| `spec.scalingPolicies[].simpleScaling.cooldownSeconds` | `int32` |  |  |  |
| `spec.scalingPolicies[].simpleScaling.minAdjustmentMagnitude` | `int32` |  |  |  |
| `spec.scalingPolicies[].predictiveScaling` | `AwsAutoScalingGroupPredictiveScalingConfig` |  |  |  |
| `spec.scalingPolicies[].predictiveScaling.targetValue` | `double` | yes |  |  |
| `spec.scalingPolicies[].predictiveScaling.predefinedMetricPairType` | `string` | yes |  |  |
| `spec.scalingPolicies[].predictiveScaling.resourceLabel` | `string` |  |  |  |
| `spec.scalingPolicies[].predictiveScaling.mode` | `string` |  |  |  |
| `spec.scalingPolicies[].predictiveScaling.schedulingBufferTimeSeconds` | `int32` |  |  |  |
| `spec.scalingPolicies[].predictiveScaling.maxCapacityBreachBehavior` | `string` |  |  |  |
| `spec.scalingPolicies[].predictiveScaling.maxCapacityBuffer` | `int32` |  |  |  |
| `spec.scheduledActions` | `[]AwsAutoScalingGroupScheduledAction` |  |  |  |
| `spec.scheduledActions[].name` | `string` | yes |  |  |
| `spec.scheduledActions[].recurrence` | `string` |  |  |  |
| `spec.scheduledActions[].timeZone` | `string` |  |  |  |
| `spec.scheduledActions[].startTime` | `string` |  |  |  |
| `spec.scheduledActions[].endTime` | `string` |  |  |  |
| `spec.scheduledActions[].minSize` | `int32` |  |  |  |
| `spec.scheduledActions[].maxSize` | `int32` |  |  |  |
| `spec.scheduledActions[].desiredCapacity` | `int32` |  |  |  |
| `spec.lifecycleHooks` | `[]AwsAutoScalingGroupLifecycleHook` |  |  |  |
| `spec.lifecycleHooks[].name` | `string` | yes |  |  |
| `spec.lifecycleHooks[].lifecycleTransition` | `string` | yes |  |  |
| `spec.lifecycleHooks[].defaultResult` | `string` |  |  |  |
| `spec.lifecycleHooks[].heartbeatTimeoutSeconds` | `int32` |  |  |  |
| `spec.lifecycleHooks[].notificationTargetArn` | `string \| valueFrom` |  |  | AwsSnsTopic (`status.outputs.topic_arn`) |
| `spec.lifecycleHooks[].roleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.lifecycleHooks[].notificationMetadata` | `string` |  |  |  |
| `spec.notifications` | `AwsAutoScalingGroupNotifications` |  |  |  |
| `spec.notifications.topic` | `string \| valueFrom` | yes |  | AwsSnsTopic (`status.outputs.topic_arn`) |
| `spec.notifications.eventTypes` | `[]string` | yes |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the group is created in. Must match the region of the
subnets, launch template, and any target groups it references.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.subnets

`[]string | valueFrom` · required

The subnets capacity is placed in (AWS calls this the VPC zone
identifier). Spread across at least two availability zones for real
fault tolerance -- the group automatically rebalances instances across
the zones these subnets cover. Reference AwsSubnet subnet_id outputs
or pass literal subnet IDs.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.launchTemplate

`AwsAutoScalingGroupLaunchTemplateRef`

The launch template every instance launches from. Exactly one of
launch_template or mixed_instances_policy must be set -- use this one
for a single-type, single-purchase-option fleet.

### spec.launchTemplate.launchTemplateId

`string | valueFrom` · required

The launch template. Reference an AwsLaunchTemplate's
launch_template_id output or pass a literal template ID ("lt-...").

- references: AwsLaunchTemplate (`status.outputs.launch_template_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLaunchTemplate, name: <that resource's name>, fieldPath: status.outputs.launch_template_id}} -- a bare string does not parse

### spec.launchTemplate.version

`string`

Which template version to launch: "$Default" (follow the template's
default version -- the AWS default and the setup that lets a template
update roll the fleet), "$Latest" (always the newest version, even
one not yet promoted), or a numeric version for a hard pin.

### spec.mixedInstancesPolicy

`AwsAutoScalingGroupMixedInstancesPolicy`

Blend instance types and purchase options in one group: a base of
On-Demand capacity plus a Spot overflow, drawn from several instance
types (or attribute-based requirements) for pool diversity. Exactly
one of launch_template or mixed_instances_policy must be set.

### spec.mixedInstancesPolicy.launchTemplate

`AwsAutoScalingGroupLaunchTemplateRef` · required

The base launch template the overrides specialize. Required.

- rule: {"required":true}

### spec.mixedInstancesPolicy.launchTemplate.launchTemplateId

`string | valueFrom` · required

The launch template. Reference an AwsLaunchTemplate's
launch_template_id output or pass a literal template ID ("lt-...").

- references: AwsLaunchTemplate (`status.outputs.launch_template_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLaunchTemplate, name: <that resource's name>, fieldPath: status.outputs.launch_template_id}} -- a bare string does not parse

### spec.mixedInstancesPolicy.launchTemplate.version

`string`

Which template version to launch: "$Default" (follow the template's
default version -- the AWS default and the setup that lets a template
update roll the fleet), "$Latest" (always the newest version, even
one not yet promoted), or a numeric version for a hard pin.

### spec.mixedInstancesPolicy.overrides

`[]AwsAutoScalingGroupMixedInstancesOverride`

Instance-type overrides. Each entry widens the pool set: an explicit
type, a different template, a capacity weight, or attribute-based
requirements. With no overrides the group uses only the base
template's type.

- rule: instance_type and instance_requirements are mutually exclusive on one override
- rule: weighted_capacity must be between 1 and 999 when set

### spec.mixedInstancesPolicy.overrides[].instanceType

`string`

An explicit instance type for this override (e.g. "m5.large").
Mutually exclusive with instance_requirements.

### spec.mixedInstancesPolicy.overrides[].weightedCapacity

`int32`

How many capacity units an instance of this type fulfills, 1-999.
Weights let heterogeneous sizes count fairly: an m5.2xlarge at
weight 4 next to an m5.large at weight 1 keeps "desired = 8"
meaningful. 0 leaves the weight unset (every instance counts as 1).

### spec.mixedInstancesPolicy.overrides[].launchTemplate

`AwsAutoScalingGroupLaunchTemplateRef`

Launch this override from a different template (e.g. an arm64-AMI
template for Graviton types next to the x86 base).

### spec.mixedInstancesPolicy.overrides[].launchTemplate.launchTemplateId

`string | valueFrom` · required

The launch template. Reference an AwsLaunchTemplate's
launch_template_id output or pass a literal template ID ("lt-...").

- references: AwsLaunchTemplate (`status.outputs.launch_template_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLaunchTemplate, name: <that resource's name>, fieldPath: status.outputs.launch_template_id}} -- a bare string does not parse

### spec.mixedInstancesPolicy.overrides[].launchTemplate.version

`string`

Which template version to launch: "$Default" (follow the template's
default version -- the AWS default and the setup that lets a template
update roll the fleet), "$Latest" (always the newest version, even
one not yet promoted), or a numeric version for a hard pin.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements

`AwsAutoScalingGroupInstanceRequirements`

Attribute-based selection for this override -- one entry that
resolves to many pools. Mutually exclusive with instance_type.

- rule: allowed_instance_types and excluded_instance_types are mutually exclusive
- rule: spot_max_price_percentage_over_lowest_price and max_spot_price_as_percentage_of_optimal_on_demand_price are mutually exclusive
- rule: bare_metal must be 'included', 'excluded', or 'required' when set
- rule: burstable_performance must be 'included', 'excluded', or 'required' when set
- rule: local_storage must be 'included', 'excluded', or 'required' when set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryMib

`AwsAutoScalingGroupIntRange` · required

Required. Memory per instance, in MiB. min is required; leave max
unset (0) for no upper bound.

- rule: {"required":true}
- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryMib.min

`int32`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryMib.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.vcpuCount

`AwsAutoScalingGroupIntRange` · required

Required. vCPUs per instance. min is required; leave max unset (0)
for no upper bound.

- rule: {"required":true}
- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.vcpuCount.min

`int32`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.vcpuCount.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.allowedInstanceTypes

`[]string`

Allow-list of instance types or families, with wildcards
("m5.large", "m5.*", "c*"). At most 400 entries. Mutually exclusive
with excluded_instance_types.

- rule: {"repeated":{"maxItems":"400"}}

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.excludedInstanceTypes

`[]string`

Deny-list of instance types or families, with wildcards. At most
400 entries. Mutually exclusive with allowed_instance_types.

- rule: {"repeated":{"maxItems":"400"}}

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.instanceGenerations

`[]string`

Instance generations to include: "current" and/or "previous".

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.cpuManufacturers

`[]string`

CPU manufacturers to include: "intel", "amd",
"amazon-web-services" (Graviton), "apple".

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.bareMetal

`string`

Bare-metal eligibility: "included", "excluded" (AWS default), or
"required".

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.burstablePerformance

`string`

Burstable (T-family) eligibility: "included", "excluded" (AWS
default), or "required".

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.requireHibernateSupport

`bool`

Only instance types that support hibernation.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.spotMaxPricePercentageOverLowestPrice

`int32`

Spot price protection: exclude types whose Spot price exceeds the
identified lowest-priced type's Spot price by more than this
percentage. Mutually exclusive with
max_spot_price_as_percentage_of_optimal_on_demand_price.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.maxSpotPriceAsPercentageOfOptimalOnDemandPrice

`int32`

Spot price protection anchored to On-Demand: exclude types whose
Spot price exceeds this percentage of the optimal type's On-Demand
price. Mutually exclusive with
spot_max_price_percentage_over_lowest_price.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.onDemandMaxPricePercentageOverLowestPrice

`int32`

On-Demand price protection: exclude types whose On-Demand price
exceeds the identified lowest-priced type's by more than this
percentage. AWS default: 20.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.localStorage

`string`

Instance-store (local disk) eligibility: "included" (AWS default),
"excluded", or "required".

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.localStorageTypes

`[]string`

Local storage technologies when instance-store is in play: "hdd"
and/or "ssd".

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.totalLocalStorageGb

`AwsAutoScalingGroupDoubleRange`

Total local (instance-store) storage, in GB.

- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.totalLocalStorageGb.min

`double`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.totalLocalStorageGb.max

`double`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryGibPerVcpu

`AwsAutoScalingGroupDoubleRange`

Memory-to-vCPU ratio, in GiB per vCPU.

- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryGibPerVcpu.min

`double`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.memoryGibPerVcpu.max

`double`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkInterfaceCount

`AwsAutoScalingGroupIntRange`

Number of network interfaces the type must support.

- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkInterfaceCount.min

`int32`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkInterfaceCount.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkBandwidthGbps

`AwsAutoScalingGroupDoubleRange`

Network bandwidth, in Gbps.

- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkBandwidthGbps.min

`double`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.networkBandwidthGbps.max

`double`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.baselineEbsBandwidthMbps

`AwsAutoScalingGroupIntRange`

Baseline (non-burst) EBS bandwidth, in Mbps.

- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.baselineEbsBandwidthMbps.min

`int32`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.baselineEbsBandwidthMbps.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorCount

`AwsAutoScalingGroupIntRange`

Number of accelerators (GPUs, FPGAs, inference chips).

- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorCount.min

`int32`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorCount.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorManufacturers

`[]string`

Accelerator manufacturers: "nvidia", "amd", "amazon-web-services",
"xilinx", "habana".

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorNames

`[]string`

Specific accelerator models (e.g. "a100", "v100", "t4",
"inferentia").

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorTypes

`[]string`

Accelerator categories: "gpu", "fpga", "inference".

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorTotalMemoryMib

`AwsAutoScalingGroupIntRange`

Total accelerator memory, in MiB.

- rule: max must be greater than or equal to min when both are set

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorTotalMemoryMib.min

`int32`

Lower bound, inclusive.

### spec.mixedInstancesPolicy.overrides[].instanceRequirements.acceleratorTotalMemoryMib.max

`int32`

Upper bound, inclusive. 0 means no upper bound.

### spec.mixedInstancesPolicy.instancesDistribution

`AwsAutoScalingGroupInstancesDistribution`

How capacity splits between On-Demand and Spot, and how each side
picks pools.

- rule: on_demand_allocation_strategy must be 'lowest-price' or 'prioritized' when set
- rule: spot_allocation_strategy must be one of: lowest-price, diversified, capacity-optimized, capacity-optimized-prioritized, price-capacity-optimized
- rule: on_demand_percentage_above_base_capacity must be between 0 and 100
- rule: spot_instance_pools only applies when spot_allocation_strategy is 'lowest-price'

### spec.mixedInstancesPolicy.instancesDistribution.onDemandAllocationStrategy

`string`

How On-Demand capacity picks pools: "lowest-price" or "prioritized"
(the override list order is a preference ranking). AWS default:
"lowest-price".

### spec.mixedInstancesPolicy.instancesDistribution.onDemandBaseCapacity

`int32`

Instances of guaranteed On-Demand capacity before the percentage
split applies -- the "always-on core" of the fleet.

### spec.mixedInstancesPolicy.instancesDistribution.onDemandPercentageAboveBaseCapacity

`int32` · optional (explicit presence)

Percentage of capacity ABOVE the base that is On-Demand, 0-100.
AWS default: 100 (all On-Demand). Explicit 0 means all-Spot above
the base -- the aggressive cost posture -- which is why this field
is optional: 0 must be distinguishable from unset.

### spec.mixedInstancesPolicy.instancesDistribution.spotAllocationStrategy

`string`

How Spot capacity picks pools: "price-capacity-optimized" (the AWS
recommendation -- weighs price AND interruption risk),
"capacity-optimized", "capacity-optimized-prioritized",
"lowest-price" (cheapest but interruption-prone), or "diversified".

### spec.mixedInstancesPolicy.instancesDistribution.spotInstancePools

`int32`

Number of Spot pools to spread across. Only valid with the
"lowest-price" strategy. AWS default: 2.

### spec.mixedInstancesPolicy.instancesDistribution.spotMaxPrice

`string`

Maximum Spot price per instance-hour, as a decimal string. AWS
default (unset): the On-Demand price -- the AWS recommendation.

### spec.minSize

`int32`

The floor the group never shrinks below. 0 is valid -- a group that
scales to zero when idle.

- rule: {"int32":{"gte":0}}

### spec.maxSize

`int32`

The ceiling the group never grows above. Scaling policies and
instance refresh honor it strictly (unless a policy explicitly allows
a predictive-scaling buffer).

- rule: {"int32":{"gte":0}}

### spec.desiredCapacity

`int32`

The capacity the group actively maintains. Leave 0 to start at
min_size and let scaling policies take over -- the declarative-fleet
default, since a literal desired count here fights the autoscaler on
every apply.

- rule: {"int32":{"gte":0}}

### spec.desiredCapacityType

`string`

What min/max/desired count: "units" (instances, the default), "vcpu",
or "memory-mib". The vCPU/memory units only make sense with
attribute-based instance requirements, where instance sizes vary.

### spec.capacityRebalance

`bool`

Proactively replace Spot instances that AWS signals as
at-elevated-risk of interruption, before the two-minute notice.
Recommended for every Spot-bearing group.

### spec.defaultCooldownSeconds

`int32`

Seconds between scaling activities initiated by simple scaling
policies. 0 keeps the AWS default (300). Step and target-tracking
policies ignore this and use instance warmup instead.

### spec.defaultInstanceWarmupSeconds

`int32`

Seconds a newly launched instance is expected to take before its
metrics are representative. Used by target tracking, instance
refresh, and rebalancing as the default warmup. Setting it (even to
0) noticeably improves scaling accuracy for fast-booting services.

### spec.healthCheckType

`string`

How instance health is judged:
- "EC2" (AWS default): instance status checks only.
- "ELB": additionally trust the load balancer's target health checks
  -- an instance failing its target group health check is replaced.
  The right choice whenever target_groups is set; without it a
  wedged-but-running process is never replaced.

### spec.healthCheckGracePeriodSeconds

`int32`

Seconds after launch before health checks can mark an instance
unhealthy -- boot-and-warm time. 0 keeps the provider default (300).
Ignored until the instance reaches the InService state.

### spec.targetGroups

`[]string | valueFrom`

Target groups whose traffic this group serves. Instances are
registered on launch and deregistered (drained) on termination.
Reference AwsLbTargetGroup target_group_arn outputs or pass literal
ARNs. Pair with health_check_type = "ELB".

- references: AwsLbTargetGroup (`status.outputs.target_group_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLbTargetGroup, name: <that resource's name>, fieldPath: status.outputs.target_group_arn}} -- a bare string does not parse

### spec.terminationPolicies

`[]string`

Which instances are terminated first on scale-in, evaluated in order:
"Default", "OldestInstance", "NewestInstance",
"OldestLaunchTemplate", "OldestLaunchConfiguration",
"ClosestToNextInstanceHour", "AllocationStrategy", or the ARN of a
custom termination Lambda. "OldestLaunchTemplate" pairs naturally
with template-version rollouts; "AllocationStrategy" keeps a mixed
fleet on its preferred pools.

### spec.maxInstanceLifetimeSeconds

`int32`

Maximum seconds any instance lives before being replaced, 86400 (1
day) to 31536000 (1 year); 0 disables. Continuous fleet hygiene:
guarantees patched AMIs and clean processes without a manual rotate.

### spec.protectFromScaleIn

`bool`

Protect instances from scale-in by default (scaling policies cannot
pick them; explicit terminations still work). For fleets whose
members hold long-lived work -- pair with lifecycle hooks for
graceful drain.

### spec.placementGroup

`string`

The placement group launched instances join (cluster/spread/
partition). A literal name -- placement groups have no Planton kind
yet.

### spec.serviceLinkedRoleArn

`string`

The service-linked IAM role the Auto Scaling service itself assumes.
A literal ARN: service-linked roles are created and owned by AWS
(not user IAM roles), and the account default
AWSServiceRoleForAutoScaling is used when unset -- which is almost
always right. Set only for a custom-suffix role (e.g. per-team KMS
grants).

### spec.enabledMetrics

`[]string`

CloudWatch group-level metrics to enable (e.g. "GroupMinSize",
"GroupMaxSize", "GroupDesiredCapacity", "GroupInServiceInstances",
"GroupPendingInstances", "GroupTerminatingInstances",
"GroupTotalInstances", and the warm-pool variants). Free of charge --
enabling them is almost always worth it for fleet observability.

### spec.suspendedProcesses

`[]string`

Auto Scaling processes to suspend, for maintenance windows or
incident response: "Launch", "Terminate", "AddToLoadBalancer",
"AlarmNotification", "AZRebalance", "HealthCheck", "InstanceRefresh",
"ReplaceUnhealthy", "ScheduledActions". Suspending "Launch" and
"Terminate" freezes the fleet entirely.

### spec.instanceRefresh

`AwsAutoScalingGroupInstanceRefresh`

Rolling replacement of instances when the launch template (or other
watched attributes) changes -- the mechanism that turns a template
update into a zero-downtime fleet rollout.

- rule: strategy must be 'Rolling'

### spec.instanceRefresh.strategy

`string` · required

The refresh strategy. "Rolling" is the only strategy AWS currently
supports. Required.

- rule: {"required":true}

### spec.instanceRefresh.triggers

`[]string`

Additional attribute changes that trigger a refresh beyond the
launch template (e.g. "tag"). Leave empty for template-only
triggers.

### spec.instanceRefresh.preferences

`AwsAutoScalingGroupInstanceRefreshPreferences`

Fine-grained rollout behavior. Unset keeps AWS defaults (90% min
healthy, no surge, no rollback).

- rule: min_healthy_percentage must be between 0 and 100
- rule: max_healthy_percentage must be between 100 and 200 when set
- rule: each checkpoint percentage must be between 1 and 100
- rule: scale_in_protected_instances must be 'Refresh', 'Ignore', or 'Wait' when set
- rule: standby_instances must be 'Terminate', 'Ignore', or 'Wait' when set

### spec.instanceRefresh.preferences.minHealthyPercentage

`int32` · optional (explicit presence)

Percentage of desired capacity that must stay InService during the
refresh, 0-100. AWS default: 90. Lower = faster, riskier waves --
explicit 0 replaces the whole fleet at once, which is why this field
is optional: 0 must be distinguishable from unset.

### spec.instanceRefresh.preferences.maxHealthyPercentage

`int32`

Upper bound on capacity during the refresh as a percentage of
desired, 100-200. Values above 100 let the refresh SURGE (launch
before terminate) -- with 110/100 min/max the fleet never dips below
full strength. 0 keeps the AWS default (100, no surge).

### spec.instanceRefresh.preferences.instanceWarmupSeconds

`int32`

Seconds a fresh instance warms before counting toward min-healthy.
0 keeps the group's default_instance_warmup (or health check grace
period).

### spec.instanceRefresh.preferences.checkpointPercentages

`[]int32`

Percentage milestones (ascending, each 1-100) where the refresh
pauses for checkpoint_delay_seconds -- a staged canary rollout:
[10, 50, 100] proves 10% before committing half the fleet.

### spec.instanceRefresh.preferences.checkpointDelaySeconds

`int32`

Seconds to wait at each checkpoint before the next wave. AWS
default: 3600 (1 hour).

### spec.instanceRefresh.preferences.autoRollback

`bool`

Roll the fleet back to its previous configuration if the refresh
fails (or a watch alarm fires). The safety net that makes
template-driven rollouts trustworthy.

### spec.instanceRefresh.preferences.alarms

`[]string | valueFrom`

CloudWatch alarms watched during the refresh: any alarm firing
fails the refresh (and rolls back when auto_rollback is set).
Reference AwsCloudwatchAlarm alarm_name outputs or pass literal
alarm names.

- references: AwsCloudwatchAlarm (`status.outputs.alarm_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCloudwatchAlarm, name: <that resource's name>, fieldPath: status.outputs.alarm_name}} -- a bare string does not parse

### spec.instanceRefresh.preferences.scaleInProtectedInstances

`string`

What happens to instances protected from scale-in: "Ignore" (AWS
default -- leave them on the old config), "Refresh" (replace them
too), or "Wait" (block until protection is removed).

### spec.instanceRefresh.preferences.standbyInstances

`string`

What happens to Standby instances: "Ignore" (AWS default),
"Terminate", or "Wait".

### spec.instanceRefresh.preferences.skipMatching

`bool`

Skip instances that already match the target configuration instead
of replacing everything -- resumes interrupted rollouts cheaply.

### spec.warmPool

`AwsAutoScalingGroupWarmPool`

A pool of pre-initialized (stopped, running, or hibernated)
instances that dramatically cuts scale-out latency for slow-booting
workloads.

- rule: pool_state must be 'Stopped', 'Running', or 'Hibernated' when set
- rule: max_group_prepared_capacity must be 0 or greater

### spec.warmPool.poolState

`string`

The state pooled instances wait in: "Stopped" (AWS default --
near-zero compute cost, seconds to start), "Running" (instant but
full price), or "Hibernated" (RAM restored from disk -- fast JVM/
cache warmup without running cost).

### spec.warmPool.minSize

`int32`

Minimum number of instances always kept in the pool.

### spec.warmPool.maxGroupPreparedCapacity

`int32` · optional (explicit presence)

Ceiling on pool size. Unset keeps the AWS default: the gap between
the group's max_size and desired capacity. Explicit 0 is
meaningful (no prepared capacity beyond min_size), which is why
this field is optional.

### spec.warmPool.reuseOnScaleIn

`bool`

Return scaled-in instances to the pool instead of terminating them
-- reuse the warm boot instead of paying for it again.

### spec.instanceMaintenancePolicy

`AwsAutoScalingGroupInstanceMaintenancePolicy`

Group-wide health bounds for REPLACEMENT operations (instance
refresh, health replacement): the percentage of capacity that must
stay in service and the surge allowed above desired.

### spec.instanceMaintenancePolicy.minHealthyPercentage

`int32`

Percentage of desired capacity that must stay InService during
replacements, 0-100.

- rule: {"int32":{"lte":100,"gte":0}}

### spec.instanceMaintenancePolicy.maxHealthyPercentage

`int32`

Upper bound on capacity during replacements as a percentage of
desired, 100-200. Set min 100 / max 110 for launch-before-terminate
on every replacement.

- rule: {"int32":{"lte":200,"gte":100}}

### spec.capacityDistributionStrategy

`string`

How capacity distributes across availability zones:
"balanced-best-effort" (AWS default -- launch in another zone when
one is impaired) or "balanced-only" (strict balance; launches wait
for the impaired zone).

### spec.forceDelete

`bool`

Delete the group without waiting for instances to terminate
gracefully. Reach for it only when tearing down a wedged group --
instances are orphaned mid-flight.

### spec.waitForCapacityTimeout

`string`

How long the IaC engine waits for the group to reach its capacity on
create/update, as a duration string (e.g. "10m"). "0" skips the wait
entirely. Unset keeps the provider default (10m). An
engine-behavior knob (both engines honor it identically), not an AWS
API field.

### spec.scalingPolicies

`[]AwsAutoScalingGroupScalingPolicy`

Scaling policies attached to the group. Target tracking is the right
default for most services; step/simple react to specific CloudWatch
alarms; predictive scaling pre-provisions for forecast load.

- rule: policy_type must be 'TargetTrackingScaling', 'StepScaling', 'SimpleScaling', or 'PredictiveScaling'
- rule: target_tracking must be set exactly when policy_type is 'TargetTrackingScaling'
- rule: step_scaling must be set exactly when policy_type is 'StepScaling'
- rule: simple_scaling must be set exactly when policy_type is 'SimpleScaling'
- rule: predictive_scaling must be set exactly when policy_type is 'PredictiveScaling'

### spec.scalingPolicies[].name

`string` · required

Policy name, unique within the group. Required.

- rule: {"required":true}

### spec.scalingPolicies[].policyType

`string` · required

The policy engine. Required.
- "TargetTrackingScaling": hold a metric at a target value -- the
  right default for services (CPU at 60%, requests-per-target).
- "StepScaling": react to a CloudWatch alarm with stepped
  adjustments.
- "SimpleScaling": the legacy single-step react-and-cooldown model.
- "PredictiveScaling": forecast load and pre-provision capacity
  (daily/weekly patterns).

- rule: {"required":true}

### spec.scalingPolicies[].estimatedInstanceWarmupSeconds

`int32`

Seconds a new instance warms before its metrics count. 0 keeps the
group default. Target tracking and step scaling only.

### spec.scalingPolicies[].targetTracking

`AwsAutoScalingGroupTargetTrackingConfig`

Configuration for "TargetTrackingScaling".

- rule: predefined_metric_type must be one of: ASGAverageCPUUtilization, ASGAverageNetworkIn, ASGAverageNetworkOut, ALBRequestCountPerTarget
- rule: exactly one of predefined_metric_type or customized_metric must be set
- rule: resource_label only applies when predefined_metric_type is 'ALBRequestCountPerTarget'

### spec.scalingPolicies[].targetTracking.targetValue

`double` · required

The value to hold the metric at (e.g. 60.0 for 60% CPU). Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.predefinedMetricType

`string`

Track a predefined group metric: "ASGAverageCPUUtilization",
"ASGAverageNetworkIn", "ASGAverageNetworkOut", or
"ALBRequestCountPerTarget". Mutually exclusive with
customized_metric.

### spec.scalingPolicies[].targetTracking.resourceLabel

`string`

Identifies the ALB target group when predefined_metric_type is
"ALBRequestCountPerTarget", in the form
"app/<lb-name>/<lb-id>/targetgroup/<tg-name>/<tg-id>" (the load
balancer's arn_suffix + "/" + the target group's arn_suffix).

### spec.scalingPolicies[].targetTracking.customizedMetric

`AwsAutoScalingGroupCustomizedMetric`

Track a custom CloudWatch metric instead of a predefined one.
Mutually exclusive with predefined_metric_type.

- rule: use either the single-metric fields (metric_name/namespace/statistic) or metrics, not both
- rule: statistic must be one of: Average, Minimum, Maximum, SampleCount, Sum
- rule: period_seconds must be 10, 30, or 60 when set

### spec.scalingPolicies[].targetTracking.customizedMetric.metricName

`string`

The metric name (single-metric form). Mutually exclusive with
metrics.

### spec.scalingPolicies[].targetTracking.customizedMetric.namespace

`string`

The metric namespace (e.g. "MyApp/Queue").

### spec.scalingPolicies[].targetTracking.customizedMetric.statistic

`string`

The statistic: "Average", "Minimum", "Maximum", "SampleCount", or
"Sum".

### spec.scalingPolicies[].targetTracking.customizedMetric.unit

`string`

The metric unit (e.g. "Percent", "Count").

### spec.scalingPolicies[].targetTracking.customizedMetric.dimensions

`[]AwsAutoScalingGroupMetricDimension`

Dimensions identifying the metric stream.

### spec.scalingPolicies[].targetTracking.customizedMetric.dimensions[].name

`string` · required

Dimension name (e.g. "QueueName"). Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.customizedMetric.dimensions[].value

`string` · required

Dimension value (e.g. "orders"). Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.customizedMetric.periodSeconds

`int32`

Metric granularity in seconds: 10, 30, or 60. High-resolution
metrics (10/30) let target tracking react in seconds.

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics

`[]AwsAutoScalingGroupMetricDataQuery`

Metric-math form: a set of query expressions combined into the
tracked value (e.g. backlog-per-instance = queue depth / instance
count). Mutually exclusive with the single-metric fields.

- rule: exactly one of expression or metric_stat must be set

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].id

`string` · required

Short identifier, unique within the query set, referenced by
expressions (e.g. "m1", "e1"). Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].expression

`string`

A metric-math expression over other query ids (e.g. "m1 / m2").
Mutually exclusive with metric_stat.

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat

`AwsAutoScalingGroupMetricStat`

A raw metric to fetch. Mutually exclusive with expression.

- rule: period_seconds must be 10, 30, or 60 when set

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.metricName

`string` · required

The metric name. Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.namespace

`string` · required

The metric namespace. Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.stat

`string` · required

The statistic to fetch (e.g. "Average", "Sum"). Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.unit

`string`

The metric unit.

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.dimensions

`[]AwsAutoScalingGroupMetricDimension`

Dimensions identifying the metric stream.

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.dimensions[].name

`string` · required

Dimension name (e.g. "QueueName"). Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.dimensions[].value

`string` · required

Dimension value (e.g. "orders"). Required.

- rule: {"required":true}

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].metricStat.periodSeconds

`int32`

Granularity in seconds: 10, 30, or 60.

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].label

`string`

Human-readable label for the query.

### spec.scalingPolicies[].targetTracking.customizedMetric.metrics[].returnData

`bool` · optional (explicit presence)

Whether this entry is the value target tracking consumes. Exactly
one entry in the set should return data (AWS default: true --
explicitly set false on intermediate entries).

### spec.scalingPolicies[].targetTracking.disableScaleIn

`bool`

Never scale IN from this policy -- it only adds capacity. For
pairing a conservative scale-out tracker with a separate, slower
scale-in mechanism.

### spec.scalingPolicies[].stepScaling

`AwsAutoScalingGroupStepScalingConfig`

Configuration for "StepScaling".

- rule: adjustment_type must be 'ChangeInCapacity', 'ExactCapacity', or 'PercentChangeInCapacity'
- rule: metric_aggregation_type must be 'Minimum', 'Maximum', or 'Average' when set
- rule: a step's scaling_adjustment of 0 is only meaningful with adjustment_type 'ExactCapacity'

### spec.scalingPolicies[].stepScaling.adjustmentType

`string` · required

How the adjustment numbers are interpreted: "ChangeInCapacity" (add/
remove N instances), "ExactCapacity" (set capacity to N), or
"PercentChangeInCapacity" (grow/shrink by N%). Required.

- rule: {"required":true}

### spec.scalingPolicies[].stepScaling.metricAggregationType

`string`

How the metric is aggregated across the breach evaluation:
"Average" (AWS default), "Minimum", or "Maximum".

### spec.scalingPolicies[].stepScaling.minAdjustmentMagnitude

`int32`

With "PercentChangeInCapacity", the minimum number of instances any
single step changes -- keeps percentage scaling meaningful on small
fleets.

### spec.scalingPolicies[].stepScaling.stepAdjustments

`[]AwsAutoScalingGroupStepAdjustment` · required

The steps, keyed by breach distance. At least one. Bounds are
decimal strings relative to the alarm threshold; an empty bound is
open-ended (negative infinity for the first lower bound, positive
infinity for the last upper bound).

- rule: {"repeated":{"minItems":"1"}}

### spec.scalingPolicies[].stepScaling.stepAdjustments[].scalingAdjustment

`int32`

The capacity change this step applies (interpreted per the policy's
adjustment_type; negative shrinks; 0 is only meaningful with
"ExactCapacity").

### spec.scalingPolicies[].stepScaling.stepAdjustments[].metricIntervalLowerBound

`string`

Lower bound of the breach range this step covers, relative to the
alarm threshold, as a decimal string (e.g. "0", "10.5"). Empty =
negative infinity.

### spec.scalingPolicies[].stepScaling.stepAdjustments[].metricIntervalUpperBound

`string`

Upper bound of the breach range, relative to the alarm threshold,
as a decimal string. Empty = positive infinity.

### spec.scalingPolicies[].simpleScaling

`AwsAutoScalingGroupSimpleScalingConfig`

Configuration for "SimpleScaling".

- rule: adjustment_type must be 'ChangeInCapacity', 'ExactCapacity', or 'PercentChangeInCapacity'
- rule: scaling_adjustment of 0 is only meaningful with adjustment_type 'ExactCapacity'

### spec.scalingPolicies[].simpleScaling.adjustmentType

`string` · required

How scaling_adjustment is interpreted: "ChangeInCapacity",
"ExactCapacity", or "PercentChangeInCapacity". Required.

- rule: {"required":true}

### spec.scalingPolicies[].simpleScaling.scalingAdjustment

`int32`

The capacity change per breach (negative shrinks; 0 is only
meaningful with "ExactCapacity" -- scale to zero).

### spec.scalingPolicies[].simpleScaling.cooldownSeconds

`int32`

Seconds after a scaling activity before this policy may fire again.
0 keeps the group's default_cooldown.

### spec.scalingPolicies[].simpleScaling.minAdjustmentMagnitude

`int32`

With "PercentChangeInCapacity", the minimum number of instances any
adjustment changes.

### spec.scalingPolicies[].predictiveScaling

`AwsAutoScalingGroupPredictiveScalingConfig`

Configuration for "PredictiveScaling".

- rule: predefined_metric_pair_type must be one of: ASGCPUUtilization, ASGNetworkIn, ASGNetworkOut, ALBRequestCount
- rule: mode must be 'ForecastOnly' or 'ForecastAndScale' when set
- rule: max_capacity_breach_behavior must be 'HonorMaxCapacity' or 'IncreaseMaxCapacity' when set
- rule: max_capacity_buffer only applies when max_capacity_breach_behavior is 'IncreaseMaxCapacity'
- rule: max_capacity_buffer must be between 0 and 100

### spec.scalingPolicies[].predictiveScaling.targetValue

`double` · required

The value to hold the scaling metric at (e.g. 60.0 for 60% CPU).
Required.

- rule: {"required":true}

### spec.scalingPolicies[].predictiveScaling.predefinedMetricPairType

`string` · required

The load/scaling metric pair to forecast: "ASGCPUUtilization",
"ASGNetworkIn", "ASGNetworkOut", or "ALBRequestCount". Required.

- rule: {"required":true}

### spec.scalingPolicies[].predictiveScaling.resourceLabel

`string`

Identifies the ALB target group when the metric pair is
"ALBRequestCount" (load balancer arn_suffix + "/" + target group
arn_suffix).

### spec.scalingPolicies[].predictiveScaling.mode

`string`

"ForecastOnly" (AWS default -- observe the forecast before trusting
it) or "ForecastAndScale" (act on it).

### spec.scalingPolicies[].predictiveScaling.schedulingBufferTimeSeconds

`int32`

Seconds ahead of the forecasted need that instances launch --
boot-and-warm lead time.

### spec.scalingPolicies[].predictiveScaling.maxCapacityBreachBehavior

`string`

What happens when the forecast exceeds max_size:
"HonorMaxCapacity" (AWS default) or "IncreaseMaxCapacity" (grow
max_size by max_capacity_buffer percent).

### spec.scalingPolicies[].predictiveScaling.maxCapacityBuffer

`int32`

Percentage buffer above forecasted capacity when
max_capacity_breach_behavior is "IncreaseMaxCapacity", 0-100.

### spec.scheduledActions

`[]AwsAutoScalingGroupScheduledAction`

Time-based capacity changes (cron or one-shot): business-hours
scale-up, overnight scale-down, batch-window pre-provisioning.

- rule: a scheduled action must set at least one of min_size, max_size, or desired_capacity
- rule: a scheduled action needs a recurrence (recurring) or a start_time (one-shot)

### spec.scheduledActions[].name

`string` · required

Action name, unique within the group. Required.

- rule: {"required":true}

### spec.scheduledActions[].recurrence

`string`

Cron expression in UTC (or time_zone), e.g. "0 8 * * MON-FRI" for
business-hours scale-up. Leave empty for a one-shot action at
start_time.

### spec.scheduledActions[].timeZone

`string`

IANA time zone for the recurrence (e.g. "America/New_York"). Unset
= UTC.

### spec.scheduledActions[].startTime

`string`

First (or only) trigger time, RFC3339 UTC (e.g.
"2026-08-01T08:00:00Z").

### spec.scheduledActions[].endTime

`string`

Last trigger time for a recurring action, RFC3339 UTC.

### spec.scheduledActions[].minSize

`int32` · optional (explicit presence)

New min_size when the action fires. Absent = leave unchanged
(which is why these are optional: 0 is a meaningful new value).

### spec.scheduledActions[].maxSize

`int32` · optional (explicit presence)

New max_size when the action fires. Absent = leave unchanged.

### spec.scheduledActions[].desiredCapacity

`int32` · optional (explicit presence)

New desired capacity when the action fires. Absent = leave
unchanged.

### spec.lifecycleHooks

`[]AwsAutoScalingGroupLifecycleHook`

Pause points in the instance lifecycle: run custom logic (warm a
cache, drain work, pull logs) while an instance waits in a
launching or terminating state.

- rule: lifecycle_transition must be 'autoscaling:EC2_INSTANCE_LAUNCHING' or 'autoscaling:EC2_INSTANCE_TERMINATING'
- rule: default_result must be 'ABANDON' or 'CONTINUE' when set
- rule: heartbeat_timeout_seconds must be between 30 and 7200 when set

### spec.lifecycleHooks[].name

`string` · required

Hook name, unique within the group. Required.

- rule: {"required":true}

### spec.lifecycleHooks[].lifecycleTransition

`string` · required

The transition the hook pauses:
"autoscaling:EC2_INSTANCE_LAUNCHING" or
"autoscaling:EC2_INSTANCE_TERMINATING". Required.

- rule: {"required":true}

### spec.lifecycleHooks[].defaultResult

`string`

What happens when the heartbeat times out without a completion
signal: "ABANDON" (AWS default for launch hooks -- roll the
instance back) or "CONTINUE" (proceed anyway).

### spec.lifecycleHooks[].heartbeatTimeoutSeconds

`int32`

Seconds the instance waits in the transition state before
default_result applies, 30-7200. AWS default: 3600.

### spec.lifecycleHooks[].notificationTargetArn

`string | valueFrom`

Where the pause notification is delivered: an SNS topic or SQS
queue ARN. Reference an AwsSnsTopic's topic_arn output or pass a
literal ARN. Unset relies on EventBridge rules watching lifecycle
events (the modern pattern).

- references: AwsSnsTopic (`status.outputs.topic_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSnsTopic, name: <that resource's name>, fieldPath: status.outputs.topic_arn}} -- a bare string does not parse

### spec.lifecycleHooks[].roleArn

`string | valueFrom`

The IAM role Auto Scaling assumes to publish to the notification
target. Required by AWS when notification_target_arn is set.
Reference an AwsIamRole's role_arn output or pass a literal ARN.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.lifecycleHooks[].notificationMetadata

`string`

Free-form JSON delivered with every notification -- routing context
for the consumer.

### spec.notifications

`AwsAutoScalingGroupNotifications`

SNS notifications for fleet lifecycle events -- the simplest way to
observe launches, terminations, and their failures.

- rule: each event type must be one of: autoscaling:EC2_INSTANCE_LAUNCH, autoscaling:EC2_INSTANCE_LAUNCH_ERROR, autoscaling:EC2_INSTANCE_TERMINATE, autoscaling:EC2_INSTANCE_TERMINATE_ERROR

### spec.notifications.topic

`string | valueFrom` · required

The SNS topic events are published to. Reference an AwsSnsTopic's
topic_arn output or pass a literal ARN. Required.

- references: AwsSnsTopic (`status.outputs.topic_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSnsTopic, name: <that resource's name>, fieldPath: status.outputs.topic_arn}} -- a bare string does not parse

### spec.notifications.eventTypes

`[]string` · required

The event types to publish. At least one of:
"autoscaling:EC2_INSTANCE_LAUNCH",
"autoscaling:EC2_INSTANCE_LAUNCH_ERROR",
"autoscaling:EC2_INSTANCE_TERMINATE",
"autoscaling:EC2_INSTANCE_TERMINATE_ERROR".

- rule: {"required":true,"repeated":{"minItems":"1"}}

## Validation Rules

- `launch_template_xor_mixed_instances`: exactly one of launch_template or mixed_instances_policy must be set
- `max_size_gte_min_size`: max_size must be greater than or equal to min_size
- `desired_within_bounds`: desired_capacity must be between min_size and max_size when set
- `desired_capacity_type_valid`: desired_capacity_type must be 'units', 'vcpu', or 'memory-mib' when set
- `health_check_type_valid`: health_check_type must be 'EC2' or 'ELB' when set
- `max_instance_lifetime_range`: max_instance_lifetime_seconds must be 0 (disabled) or between 86400 (1 day) and 31536000 (1 year)
- `capacity_distribution_strategy_valid`: capacity_distribution_strategy must be 'balanced-only' or 'balanced-best-effort' when set
- `suspended_processes_valid`: each suspended process must be one of: Launch, Terminate, AddToLoadBalancer, AlarmNotification, AZRebalance, HealthCheck, InstanceRefresh, ReplaceUnhealthy, ScheduledActions

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsAutoScalingGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.autoscaling_group_name` | `string` | The name of the auto-scaling group (metadata.name). The handle the AWS CLI, CloudWatch dimensions (AutoScalingGroupName), and ECS capacity providers reference. |
| `status.outputs.autoscaling_group_arn` | `string` | The ARN of the auto-scaling group, for IAM policies and EventBridge rules scoped to this group. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.subnets` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.launchTemplate.launchTemplateId` | AwsLaunchTemplate | `status.outputs.launch_template_id` |
| `spec.mixedInstancesPolicy.launchTemplate.launchTemplateId` | AwsLaunchTemplate | `status.outputs.launch_template_id` |
| `spec.mixedInstancesPolicy.overrides[].launchTemplate.launchTemplateId` | AwsLaunchTemplate | `status.outputs.launch_template_id` |
| `spec.targetGroups` | AwsLbTargetGroup | `status.outputs.target_group_arn` |
| `spec.instanceRefresh.preferences.alarms` | AwsCloudwatchAlarm | `status.outputs.alarm_name` |
| `spec.lifecycleHooks[].notificationTargetArn` | AwsSnsTopic | `status.outputs.topic_arn` |
| `spec.lifecycleHooks[].roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.notifications.topic` | AwsSnsTopic | `status.outputs.topic_arn` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsEcsCluster | `spec.ec2CapacityProviders[].autoScalingGroupArn` | `status.outputs.autoscaling_group_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
