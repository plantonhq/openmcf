# AwsEcsCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsEcsClusterSpec defines an ECS cluster: the logical boundary that
groups services and tasks, decides where their containers run (Fargate,
EC2 capacity providers, or a blend), and carries cluster-wide posture --
Container Insights observability, ECS Exec auditing, Fargate storage
encryption, and the Service Connect default namespace.

Capacity is a spectrum. A serverless cluster associates the AWS-managed
FARGATE / FARGATE_SPOT providers and never thinks about instances. An
EC2-backed cluster defines ec2_capacity_providers -- each wrapping a
referenced AwsAutoScalingGroup whose fleet ECS scales up and down through
managed scaling -- and services blend across all of it by provider name
in their capacity_provider_strategy. The cluster itself is free; only
the tasks and instances it schedules cost money.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEcsCluster
metadata:
  name: ecs-cluster-demo
spec:
  region: us-west-2
  containerInsights: enhanced
  capacityProviders:
    - FARGATE
    - FARGATE_SPOT
  defaultCapacityProviderStrategy:
    - capacityProvider: FARGATE
      base: 1
      weight: 1
    - capacityProvider: FARGATE_SPOT
      weight: 4
  executeCommandConfiguration:
    logging: DEFAULT
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.containerInsights` | `string` |  |  |  |
| `spec.capacityProviders` | `[]string` |  |  |  |
| `spec.ec2CapacityProviders` | `[]AwsEcsClusterEc2CapacityProvider` |  |  |  |
| `spec.ec2CapacityProviders[].name` | `string` | yes |  |  |
| `spec.ec2CapacityProviders[].autoScalingGroupArn` | `string \| valueFrom` | yes |  | AwsAutoScalingGroup (`status.outputs.autoscaling_group_arn`) |
| `spec.ec2CapacityProviders[].managedScaling` | `AwsEcsClusterManagedScaling` |  |  |  |
| `spec.ec2CapacityProviders[].managedScaling.status` | `string` |  |  |  |
| `spec.ec2CapacityProviders[].managedScaling.targetCapacity` | `int32` |  |  |  |
| `spec.ec2CapacityProviders[].managedScaling.minimumScalingStepSize` | `int32` |  |  |  |
| `spec.ec2CapacityProviders[].managedScaling.maximumScalingStepSize` | `int32` |  |  |  |
| `spec.ec2CapacityProviders[].managedScaling.instanceWarmupPeriodSeconds` | `int32` |  |  |  |
| `spec.ec2CapacityProviders[].managedTerminationProtection` | `string` |  |  |  |
| `spec.ec2CapacityProviders[].managedDraining` | `string` |  |  |  |
| `spec.defaultCapacityProviderStrategy` | `[]AwsEcsClusterCapacityProviderStrategy` |  |  |  |
| `spec.defaultCapacityProviderStrategy[].capacityProvider` | `string` | yes |  |  |
| `spec.defaultCapacityProviderStrategy[].base` | `int32` |  |  |  |
| `spec.defaultCapacityProviderStrategy[].weight` | `int32` |  |  |  |
| `spec.executeCommandConfiguration` | `AwsEcsClusterExecuteCommandConfiguration` |  |  |  |
| `spec.executeCommandConfiguration.logging` | `string` |  |  |  |
| `spec.executeCommandConfiguration.logConfiguration` | `AwsEcsClusterExecuteCommandLogConfiguration` |  |  |  |
| `spec.executeCommandConfiguration.logConfiguration.cloudWatchLogGroupName` | `string` |  |  |  |
| `spec.executeCommandConfiguration.logConfiguration.cloudWatchEncryptionEnabled` | `bool` |  |  |  |
| `spec.executeCommandConfiguration.logConfiguration.s3BucketName` | `string` |  |  |  |
| `spec.executeCommandConfiguration.logConfiguration.s3KeyPrefix` | `string` |  |  |  |
| `spec.executeCommandConfiguration.logConfiguration.s3BucketEncryptionEnabled` | `bool` |  |  |  |
| `spec.executeCommandConfiguration.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.managedStorageConfiguration` | `AwsEcsClusterManagedStorageConfiguration` |  |  |  |
| `spec.managedStorageConfiguration.fargateEphemeralStorageKmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.managedStorageConfiguration.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.serviceConnectNamespaceArn` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.containerInsights

`string`

CloudWatch Container Insights for the cluster:
"enabled" -- metrics and logs at the cluster/service/task level.
"enhanced" -- adds container-level observability with automatic
  dashboards (recommended for production; higher CloudWatch cost).
"disabled" -- no Insights telemetry.
Unset keeps the account's default setting. Updatable in place.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["enabled","enhanced","disabled"]}}

### spec.capacityProviders

`[]string`

The AWS-managed serverless capacity providers to associate:
"FARGATE" and/or "FARGATE_SPOT". These are built into every account
-- associating them here is what lets services in this cluster name
them in a capacity_provider_strategy. EC2 capacity is defined
separately in ec2_capacity_providers; both sets associate onto the
cluster together.

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["FARGATE","FARGATE_SPOT"]}}}}

### spec.ec2CapacityProviders

`[]AwsEcsClusterEc2CapacityProvider`

EC2 capacity providers, each wrapping a referenced auto-scaling
group. ECS's managed scaling drives the group's desired count from
task demand -- you size the ASG's bounds, ECS turns instances on and
off. Each entry materializes as its own capacity provider resource
(keyed by name, so adding or removing one never disturbs the others)
and is automatically associated with the cluster alongside the
Fargate built-ins. Services reference entries by name in their
capacity_provider_strategy.

- rule: capacity provider names may not start with 'aws', 'ecs', or 'fargate' (reserved by AWS)
- rule: managed_termination_protection must be 'ENABLED' or 'DISABLED' when set
- rule: managed_draining must be 'ENABLED' or 'DISABLED' when set

### spec.ec2CapacityProviders[].name

`string` · required

The capacity provider name -- what services put in their
capacity_provider_strategy. 1-255 characters: letters, digits,
hyphens, underscores; must not start with "aws", "ecs", or "fargate"
(AWS reserves those prefixes).

- rule: {"required":true,"string":{"pattern":"^[a-zA-Z0-9_-]{1,255}$"}}

### spec.ec2CapacityProviders[].autoScalingGroupArn

`string | valueFrom` · required

The auto-scaling group that provides the instances. Reference an
AwsAutoScalingGroup's autoscaling_group_arn output or pass a literal
ARN. The group's launch template decides the instance shape (use an
ECS-optimized AMI whose agent joins this cluster via user data);
ECS's managed scaling then drives the group's desired capacity
between the group's own min/max bounds. ForceNew: changing the group
replaces the provider.

- references: AwsAutoScalingGroup (`status.outputs.autoscaling_group_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsAutoScalingGroup, name: <that resource's name>, fieldPath: status.outputs.autoscaling_group_arn}} -- a bare string does not parse

### spec.ec2CapacityProviders[].managedScaling

`AwsEcsClusterManagedScaling`

ECS-managed scaling of the auto-scaling group. Leave unset to keep
AWS's defaults (managed scaling enabled, target capacity 100); set
it to tune headroom and scaling step bounds.

- rule: status must be 'ENABLED' or 'DISABLED' when set
- rule: target_capacity must be between 1 and 100 when set
- rule: minimum_scaling_step_size must be between 1 and 10000 when set
- rule: maximum_scaling_step_size must be between 1 and 10000 when set
- rule: instance_warmup_period_seconds must be between 0 and 10000
- rule: maximum_scaling_step_size must be greater than or equal to minimum_scaling_step_size when both are set

### spec.ec2CapacityProviders[].managedScaling.status

`string`

Managed scaling on or off: "ENABLED" (AWS default -- ECS sizes the
group) or "DISABLED" (you size the group yourself; ECS only places
tasks on what exists).

### spec.ec2CapacityProviders[].managedScaling.targetCapacity

`int32`

Utilization target for the group, 1-100 percent. 100 (AWS default)
runs instances fully packed; a lower value (e.g. 80) keeps headroom
so new tasks place without waiting for an instance launch.

### spec.ec2CapacityProviders[].managedScaling.minimumScalingStepSize

`int32`

Smallest scale-out step, 1-10000 instances. AWS default: 1.

### spec.ec2CapacityProviders[].managedScaling.maximumScalingStepSize

`int32`

Largest scale-out step, 1-10000 instances. AWS default: 10000.

### spec.ec2CapacityProviders[].managedScaling.instanceWarmupPeriodSeconds

`int32`

Seconds a newly launched instance warms up before counting toward
capacity metrics, 0-10000. AWS default: 300.

### spec.ec2CapacityProviders[].managedTerminationProtection

`string`

Protect instances running non-daemon tasks from scale-in
termination: "ENABLED" or "DISABLED". Enabling it requires the
auto-scaling group itself to enable new-instance scale-in protection
(protect_from_scale_in on the group) -- AWS rejects the provider
otherwise. The safe default for task-dense clusters.

### spec.ec2CapacityProviders[].managedDraining

`string`

Gracefully drain tasks off instances the group is terminating:
"ENABLED" (AWS default) or "DISABLED". Draining is what makes
scale-in and instance refresh invisible to services.

### spec.defaultCapacityProviderStrategy

`[]AwsEcsClusterCapacityProviderStrategy`

The cluster's default capacity provider strategy -- what ECS uses
when a service or run-task does not declare its own strategy. Name
any associated provider: the Fargate built-ins or an
ec2_capacity_providers entry. Example: FARGATE base 1 / weight 1 +
FARGATE_SPOT weight 4 keeps one guaranteed On-Demand task and runs
~80% of scaled capacity on Spot.

### spec.defaultCapacityProviderStrategy[].capacityProvider

`string` · required

The capacity provider: "FARGATE", "FARGATE_SPOT", or the name of an
ec2_capacity_providers entry.

- rule: {"required":true}

### spec.defaultCapacityProviderStrategy[].base

`int32`

Minimum number of tasks guaranteed on this provider before weights
apply. Only one entry of the strategy may set a non-zero base.

- rule: {"int32":{"lte":100000,"gte":0}}

### spec.defaultCapacityProviderStrategy[].weight

`int32`

Relative share of tasks beyond the bases. Example: FARGATE weight 1 +
FARGATE_SPOT weight 4 runs ~80% of scaled tasks on Spot.

- rule: {"int32":{"lte":1000,"gte":0}}

### spec.executeCommandConfiguration

`AwsEcsClusterExecuteCommandConfiguration`

ECS Exec auditing for the cluster: where interactive exec sessions
(`aws ecs execute-command`) are logged and how session traffic is
encrypted. Without this block, exec sessions still work when a
service enables them -- they are simply not centrally audited.

- rule: logging must be 'DEFAULT', 'OVERRIDE', or 'NONE' when set
- rule: logging 'OVERRIDE' requires log_configuration with at least one destination
- rule: log_configuration only applies when logging is 'OVERRIDE'

### spec.executeCommandConfiguration.logging

`string`

Log destination behavior for exec sessions:
"DEFAULT" (AWS default) -- sessions log to the task's own awslogs
  configuration.
"OVERRIDE" -- sessions log to the destinations in log_configuration.
"NONE" -- exec works but sessions are not logged (avoid outside
  sandboxes; unaudited interactive access defeats compliance).

### spec.executeCommandConfiguration.logConfiguration

`AwsEcsClusterExecuteCommandLogConfiguration`

Custom destinations for exec session logs. Only used (and required)
when logging is "OVERRIDE".

- rule: provide cloud_watch_log_group_name and/or s3_bucket_name

### spec.executeCommandConfiguration.logConfiguration.cloudWatchLogGroupName

`string`

The CloudWatch log group session logs are sent to. The group must
already exist -- ECS does not create it.

### spec.executeCommandConfiguration.logConfiguration.cloudWatchEncryptionEnabled

`bool`

Require the CloudWatch log group to be KMS-encrypted; the send fails
if it is not. Pair with an encrypted log group for compliance
postures.

### spec.executeCommandConfiguration.logConfiguration.s3BucketName

`string`

The S3 bucket session logs are written to.

### spec.executeCommandConfiguration.logConfiguration.s3KeyPrefix

`string`

Key prefix for session log objects within the bucket.

### spec.executeCommandConfiguration.logConfiguration.s3BucketEncryptionEnabled

`bool`

Require the S3 bucket to be encrypted; the write fails if it is not.

### spec.executeCommandConfiguration.kmsKeyId

`string | valueFrom`

A KMS key to encrypt the exec session traffic between client and
container. Reference an AwsKmsKey's key_arn output or pass a literal
key ARN/ID. Unset uses TLS without customer-managed encryption.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.managedStorageConfiguration

`AwsEcsClusterManagedStorageConfiguration`

Customer-managed KMS encryption for Fargate ephemeral task storage
and managed storage -- the compliance posture for regulated
workloads. Unset uses AWS-owned keys (data is still encrypted).

### spec.managedStorageConfiguration.fargateEphemeralStorageKmsKeyId

`string | valueFrom`

The KMS key encrypting Fargate ephemeral task storage. Reference an
AwsKmsKey's key_arn output or pass a literal key ARN. The key policy
must grant the Fargate service principal decrypt/generate rights --
AWS rejects the cluster configuration otherwise.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.managedStorageConfiguration.kmsKeyId

`string | valueFrom`

The KMS key for other ECS-managed storage. Reference an AwsKmsKey's
key_arn output or pass a literal key ARN.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.serviceConnectNamespaceArn

`string`

The AWS Cloud Map namespace (by ARN) that Service Connect uses by
default for services in this cluster. Services can override it;
setting it here is what lets a whole environment share one service
mesh namespace without per-service wiring.

## Validation Rules

- `default_strategy_names_associated_providers`: every default_capacity_provider_strategy entry must name an associated provider -- a capacity_providers built-in (FARGATE/FARGATE_SPOT) or an ec2_capacity_providers entry
- `default_strategy_single_base`: only one default_capacity_provider_strategy entry may set a non-zero base
- `ec2_capacity_provider_names_unique`: ec2_capacity_providers entries must have unique names

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsEcsCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_name` | `string` | The cluster name (mirrors metadata.name). What the AWS CLI and the ECS agent's ECS_CLUSTER setting address. |
| `status.outputs.cluster_arn` | `string` | The ARN of the cluster. The join key: an AwsEcsService's cluster_arn references this output. |
| `status.outputs.capacity_provider_names` | `[]string` | Every capacity provider associated with the cluster -- the Fargate built-ins plus the names of the folded EC2 capacity providers. The vocabulary services can use in a capacity_provider_strategy. |
| `status.outputs.capacity_provider_arns` | `[]string` | The ARNs of the EC2 capacity providers this cluster defines (empty for Fargate-only clusters), in the order declared in the spec. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.ec2CapacityProviders[].autoScalingGroupArn` | AwsAutoScalingGroup | `status.outputs.autoscaling_group_arn` |
| `spec.executeCommandConfiguration.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.managedStorageConfiguration.fargateEphemeralStorageKmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.managedStorageConfiguration.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsEcsService | `spec.clusterArn` | `status.outputs.cluster_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
