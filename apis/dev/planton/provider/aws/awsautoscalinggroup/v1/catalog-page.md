# AWS Auto Scaling Group

Deploys an EC2 Auto Scaling group: the fleet manager that launches
instances from an `AwsLaunchTemplate`, keeps them healthy across subnets,
registers them into `AwsLbTargetGroup` nodes, and scales them with
policies, schedules, and forecasts -- with scaling policies, scheduled
actions, lifecycle hooks, and notifications managed together.

## What Gets Created

When you deploy an AwsAutoScalingGroup resource, Planton provisions:

- **Auto Scaling group** — an `aws_autoscaling_group` /
  `autoscaling.Group` named from `metadata.name`, with capacity bounds,
  subnet spread, launch template (or mixed-instances policy), health
  model, instance refresh, warm pool, and maintenance policy
- **Scaling policies** — one `aws_autoscaling_policy` per entry in
  `scalingPolicies` (target tracking, step, simple, or predictive)
- **Scheduled actions** — one `aws_autoscaling_schedule` per entry in
  `scheduledActions`
- **Lifecycle hooks** — one `aws_autoscaling_lifecycle_hook` per entry in
  `lifecycleHooks` (managed standalone so each stays updatable)
- **Notifications** — an `aws_autoscaling_notification` wiring lifecycle
  events to the referenced SNS topic

Identity tags are applied through the group's native tag mechanism with
`propagate_at_launch`, so every launched instance carries them.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **Subnets** (`AwsSubnet`) in at least two availability zones for real fault tolerance.
- **A launch template** (`AwsLaunchTemplate`) carrying an AMI — groups require their template to be fully formed.
- **Target groups** (`AwsLbTargetGroup`) if the fleet serves load-balanced traffic.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAutoScalingGroup
metadata:
  name: web
spec:
  region: us-west-2
  subnets:
    - valueFrom:
        kind: AwsSubnet
        name: private-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: private-b
        fieldPath: status.outputs.subnet_id
  launchTemplate:
    launchTemplateId:
      valueFrom:
        kind: AwsLaunchTemplate
        name: web
        fieldPath: status.outputs.launch_template_id
  minSize: 2
  maxSize: 10
```

```shell
planton apply -f autoscaling-group.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region. Must match the subnets, template, and target groups. | Required; non-empty |
| `subnets` | `string[] \| valueFrom` | Subnets capacity is placed in. Spread across ≥2 AZs. Defaults to referencing `AwsSubnet` `subnet_id` outputs. | Required; ≥1 entry |
| `launchTemplate` XOR `mixedInstancesPolicy` | `object` | What to launch: a single template reference, or an On-Demand/Spot blend with per-type overrides. | Exactly one must be set |
| `maxSize` | `int` | Capacity ceiling. | ≥ `minSize` |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `minSize` | `int` | `0` | Capacity floor. 0 allows scale-to-zero. |
| `desiredCapacity` | `int` | starts at `minSize` | Actively maintained capacity. Leave unset and let policies own the number. |
| `desiredCapacityType` | `string` | `units` | What min/max/desired count: `units`, `vcpu`, or `memory-mib`. |
| `capacityRebalance` | `bool` | `false` | Proactively replace at-risk Spot instances. Recommended for Spot fleets. |
| `defaultCooldownSeconds` | `int` | `300` | Cooldown between simple-scaling activities. |
| `defaultInstanceWarmupSeconds` | `int` | — | Warmup used by target tracking, refresh, and rebalancing. |
| `healthCheckType` | `string` | `EC2` | `EC2` (status checks) or `ELB` (also trust target-group health — pair with `targetGroups`). |
| `healthCheckGracePeriodSeconds` | `int` | `300` | Boot-and-warm time before health checks can mark unhealthy. |
| `targetGroups` | `string[] \| valueFrom` | `[]` | Target groups the fleet serves; instances register on launch and drain on termination. |
| `terminationPolicies` | `string[]` | `Default` | Scale-in victim ordering: `OldestInstance`, `OldestLaunchTemplate`, `AllocationStrategy`, a Lambda ARN, etc. |
| `maxInstanceLifetimeSeconds` | `int` | disabled | Replace every instance after 1 day–1 year: continuous fleet hygiene. |
| `protectFromScaleIn` | `bool` | `false` | New instances start protected from scale-in. |
| `placementGroup` | `string` | — | Placement group name (literal — no placement-group kind yet). |
| `serviceLinkedRoleArn` | `string` | account default | Custom service-linked role ARN (literal — AWS-owned roles). |
| `enabledMetrics` | `string[]` | `[]` | Free group-level CloudWatch metrics (`GroupInServiceInstances`, ...). |
| `suspendedProcesses` | `string[]` | `[]` | Freeze processes (`Launch`, `Terminate`, `AZRebalance`, ...) for maintenance. |
| `instanceRefresh` | `object` | — | Rolling replacement on template change: health bounds, surge, checkpoints, alarm watch, auto-rollback. |
| `warmPool` | `object` | — | Pre-initialized capacity: `Stopped`/`Running`/`Hibernated`, sizes, reuse on scale-in. |
| `instanceMaintenancePolicy` | `object` | — | Group-wide min/max-healthy bounds for every replacement operation. |
| `capacityDistributionStrategy` | `string` | `balanced-best-effort` | Zone balance behavior: `balanced-only` or `balanced-best-effort`. |
| `forceDelete` | `bool` | `false` | Delete without draining (wedged-group escape hatch). |
| `waitForCapacityTimeout` | `string` | `10m` | How long the IaC engine waits for capacity; `"0"` skips. |
| `scalingPolicies` | `object[]` | `[]` | Target tracking / step / simple / predictive policies (discriminated by `policyType`). |
| `scheduledActions` | `object[]` | `[]` | Cron or one-shot capacity changes with time zones. |
| `lifecycleHooks` | `object[]` | `[]` | Launch/terminate pause points with SNS/SQS delivery and IAM role references. |
| `notifications` | `object` | — | SNS topic reference + lifecycle event types. |

## Examples

### Spot-majority mixed fleet with an On-Demand base

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAutoScalingGroup
metadata:
  name: workers
spec:
  region: us-west-2
  subnets:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  mixedInstancesPolicy:
    launchTemplate:
      launchTemplateId:
        valueFrom:
          kind: AwsLaunchTemplate
          name: workers
          fieldPath: status.outputs.launch_template_id
    overrides:
      - instanceType: m6i.large
      - instanceType: m5.large
      - instanceType: m5a.large
    instancesDistribution:
      onDemandBaseCapacity: 2
      onDemandPercentageAboveBaseCapacity: 0
      spotAllocationStrategy: price-capacity-optimized
  minSize: 2
  maxSize: 20
  capacityRebalance: true
```

### Web service behind a target group with target tracking and rolling refresh

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAutoScalingGroup
metadata:
  name: api
spec:
  region: us-west-2
  subnets:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  launchTemplate:
    launchTemplateId:
      valueFrom:
        kind: AwsLaunchTemplate
        name: api
        fieldPath: status.outputs.launch_template_id
  minSize: 2
  maxSize: 12
  healthCheckType: ELB
  healthCheckGracePeriodSeconds: 120
  targetGroups:
    - valueFrom:
        kind: AwsLbTargetGroup
        name: api
        fieldPath: status.outputs.target_group_arn
  instanceRefresh:
    strategy: Rolling
    preferences:
      minHealthyPercentage: 100
      maxHealthyPercentage: 110
      autoRollback: true
  scalingPolicies:
    - name: cpu-target
      policyType: TargetTrackingScaling
      targetTracking:
        targetValue: 60
        predefinedMetricType: ASGAverageCPUUtilization
```

### Business-hours schedule with a warm pool

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAutoScalingGroup
metadata:
  name: batch-workers
spec:
  region: us-west-2
  subnets:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  launchTemplate:
    launchTemplateId:
      valueFrom:
        kind: AwsLaunchTemplate
        name: batch-workers
        fieldPath: status.outputs.launch_template_id
  minSize: 0
  maxSize: 20
  scheduledActions:
    - name: business-hours-up
      recurrence: "0 8 * * MON-FRI"
      timeZone: America/New_York
      minSize: 4
    - name: overnight-down
      recurrence: "0 20 * * *"
      timeZone: America/New_York
      minSize: 0
      desiredCapacity: 0
  warmPool:
    poolState: Stopped
    minSize: 2
    reuseOnScaleIn: true
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `autoscaling_group_name` | Group name — the `AutoScalingGroupName` CloudWatch dimension and the handle ECS capacity providers take |
| `autoscaling_group_arn` | Group ARN, for IAM policies and EventBridge rules scoped to this group |

## Related Components

- [AwsLaunchTemplate](/docs/catalog/aws/awslaunchtemplate) — the blueprint this group launches from
- [AwsLbTargetGroup](/docs/catalog/aws/awslbtargetgroup) — where the fleet's traffic comes from
- [AwsSubnet](/docs/catalog/aws/awssubnet) — the zones capacity spreads across
- [AwsCloudwatchAlarm](/docs/catalog/aws/awscloudwatchalarm) — alarms watched during instance refresh; triggers for step policies
- [AwsSnsTopic](/docs/catalog/aws/awssnstopic) — lifecycle-hook and notification delivery
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — the role lifecycle hooks publish with
