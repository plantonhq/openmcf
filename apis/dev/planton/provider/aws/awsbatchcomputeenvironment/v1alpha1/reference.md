# AwsBatchComputeEnvironment

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsBatchComputeEnvironmentSpec defines an AWS Batch MANAGED compute
environment: the elastic pool of compute (EC2 On-Demand, EC2 Spot, Fargate,
or Fargate Spot) that AWS Batch scales up and down to run submitted jobs.

A compute environment is one node of the Batch resource graph, not the whole
story: jobs are submitted to an AwsBatchJobQueue (which maps to one or more
compute environments in priority order) using an AwsBatchJobDefinition (the
container blueprint). Keeping those first-class lets one queue span an
On-Demand primary and a Spot overflow environment, and lets a compute
environment be replaced behind a queue without touching the queue.

UNMANAGED compute environments (where you run your own ECS container
instances) are not modeled; the modules always create type MANAGED, which is
why compute_resources is required here even though AWS makes it optional at
the API level.

Update semantics worth knowing before you plan a change: AWS Batch can only
update a compute environment's infrastructure IN PLACE when (a) the
environment uses the Batch service-linked role (leave service_role unset)
AND (b) allocation_strategy is one of BEST_FIT_PROGRESSIVE,
SPOT_CAPACITY_OPTIMIZED, or SPOT_PRICE_CAPACITY_OPTIMIZED. Outside that
envelope, changes to most compute_resources fields (instance types, key
pair, launch template, AMI configuration, security groups, subnets, tags)
REPLACE the whole environment.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBatchComputeEnvironment
metadata:
  name: test-batch
  id: test-batch
  org: test-org
  env: dev
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: test
    pulumi.planton.dev/project: test
    pulumi.planton.dev/stack.name: dev.AwsBatchComputeEnvironment.test-batch
spec:
  region: us-west-2
  computeResources:
    type: FARGATE
    maxVcpus: 256
    subnetIds:
      - value: subnet-0a1b2c3d4e5f00001
      - value: subnet-0a1b2c3d4e5f00002
    securityGroupIds:
      - value: sg-0a1b2c3d4e5f00001
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.state` | `string` |  | `ENABLED` |  |
| `spec.serviceRole` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.computeResources` | `AwsBatchComputeResources` | yes |  |  |
| `spec.computeResources.type` | `string` | yes |  |  |
| `spec.computeResources.maxVcpus` | `int32` | yes |  |  |
| `spec.computeResources.minVcpus` | `int32` |  | `0` |  |
| `spec.computeResources.desiredVcpus` | `int32` |  |  |  |
| `spec.computeResources.subnetIds` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.computeResources.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.computeResources.instanceTypes` | `[]string` |  |  |  |
| `spec.computeResources.allocationStrategy` | `string` |  |  |  |
| `spec.computeResources.instanceRole` | `string \| valueFrom` |  |  | AwsIamInstanceProfile (`status.outputs.instance_profile_arn`) |
| `spec.computeResources.ec2KeyPair` | `string` |  |  |  |
| `spec.computeResources.bidPercentage` | `int32` |  |  |  |
| `spec.computeResources.spotIamFleetRole` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.computeResources.launchTemplate` | `AwsBatchLaunchTemplate` |  |  |  |
| `spec.computeResources.launchTemplate.launchTemplateId` | `string \| valueFrom` | yes |  | AwsLaunchTemplate (`status.outputs.launch_template_id`) |
| `spec.computeResources.launchTemplate.version` | `string` |  |  |  |
| `spec.computeResources.ec2Configurations` | `[]AwsBatchEc2Configuration` |  |  |  |
| `spec.computeResources.ec2Configurations[].imageType` | `string` |  |  |  |
| `spec.computeResources.ec2Configurations[].imageIdOverride` | `string` |  |  |  |
| `spec.computeResources.ec2Configurations[].imageKubernetesVersion` | `string` |  |  |  |
| `spec.computeResources.placementGroup` | `string` |  |  |  |
| `spec.computeResources.resourceTags` | `map<string, string>` |  |  |  |
| `spec.eksConfiguration` | `AwsBatchEksConfiguration` |  |  |  |
| `spec.eksConfiguration.eksClusterArn` | `string \| valueFrom` | yes |  | AwsEksCluster (`status.outputs.cluster_arn`) |
| `spec.eksConfiguration.kubernetesNamespace` | `string` | yes |  |  |
| `spec.updatePolicy` | `AwsBatchUpdatePolicy` |  |  |  |
| `spec.updatePolicy.terminateJobsOnUpdate` | `bool` |  |  |  |
| `spec.updatePolicy.jobExecutionTimeoutMinutes` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the Batch compute environment is created.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.state

`string` · optional (explicit presence)

Whether the compute environment accepts jobs from associated queues.
When DISABLED, no new jobs are dispatched to it, but running jobs finish
and the environment can still scale in. Disabling is also how a compute
environment is drained before deletion or replacement behind a queue.

- default: `ENABLED`
- rule: {"string":{"in":["ENABLED","DISABLED"]}}

### spec.serviceRole

`string | valueFrom`

The IAM role AWS Batch assumes to manage compute on your behalf.
LEAVE UNSET for the recommended path: AWS Batch then uses (and
auto-creates) the AWSServiceRoleForBatch service-linked role -- which is
also a precondition for in-place infrastructure updates (see the message
comment). Set a custom role only when your org mandates one, and expect
most compute_resources changes to replace the environment in that mode.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.computeResources

`AwsBatchComputeResources` · required

The compute infrastructure backing this environment: resource type,
vCPU scaling bounds, VPC placement, and (for EC2/SPOT) instance
selection. Required because the modules always create a MANAGED
environment, and AWS requires compute resources for MANAGED.

- rule: {"required":true}
- rule: EC2 and SPOT compute environments need instance_role -- reference an AwsIamInstanceProfile (wrapping the ECS instance role) or pass a literal instance profile ARN
- rule: SPOT compute environments using the BEST_FIT allocation strategy (or no strategy, which defaults to BEST_FIT) need spot_iam_fleet_role -- either reference the Spot Fleet role, or pick a capacity-optimized allocation_strategy which does not use Spot Fleet
- rule: FARGATE and FARGATE_SPOT compute environments require at least one security group in security_group_ids (Fargate task ENIs must have security groups)
- rule: FARGATE and FARGATE_SPOT environments are serverless -- remove the EC2-only fields (min/desired vCPUs, instance types, allocation strategy, instance role, key pair, bid percentage, Spot fleet role, launch template, EC2 configuration, placement group, resource tags)
- rule: bid_percentage and spot_iam_fleet_role only apply to SPOT compute environments -- remove them or change type to SPOT

### spec.computeResources.type

`string` · required

The compute resource type.
  EC2:          On-Demand EC2 instances.
  SPOT:         EC2 Spot instances (interruptible, up to ~90% cheaper).
  FARGATE:      Serverless containers (AWS manages all instances).
  FARGATE_SPOT: Serverless containers at Spot pricing.
A job queue can only mix environments of one family (EC2/SPOT together,
or FARGATE/FARGATE_SPOT together) -- never both families.

- rule: {"required":true,"string":{"in":["EC2","SPOT","FARGATE","FARGATE_SPOT"]}}

### spec.computeResources.maxVcpus

`int32` · required

The maximum vCPUs the environment can scale out to. For Fargate this
caps total concurrent vCPU capacity across all running jobs. This is
the one sizing knob AWS allows updating on EVERY environment,
regardless of service role or allocation strategy.

- rule: {"required":true,"int32":{"gte":1}}

### spec.computeResources.minVcpus

`int32` · optional (explicit presence)

The vCPU floor maintained even when no jobs are runnable. EC2/SPOT
only. Keep the default 0 so the environment scales to zero when idle --
a non-zero floor keeps instances (and their cost) warm for
latency-sensitive queues.

- default: `0`

### spec.computeResources.desiredVcpus

`int32`

The initial vCPU target at environment creation. EC2/SPOT only. AWS
Batch continuously adjusts the actual desired capacity between
min_vcpus and max_vcpus based on queue demand, so treat this as a
starting point, not a setpoint.

### spec.computeResources.subnetIds

`[]string | valueFrom` · required

The VPC subnets where compute is launched. Spread across multiple
Availability Zones for capacity diversity -- especially for SPOT, where
more pools mean fewer interruptions.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.computeResources.securityGroupIds

`[]string | valueFrom`

The security groups attached to compute resources (and to Fargate task
ENIs). REQUIRED for FARGATE/FARGATE_SPOT. For EC2/SPOT they may be
omitted only when the launch template supplies its own.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.computeResources.instanceTypes

`[]string`

The EC2 instance types (or families) Batch may launch. EC2/SPOT only.
Use "optimal" to let Batch pick from the C, M, and R families to match
each job's resource requirements.
Examples: ["optimal"], ["m5", "c5"], ["c5.xlarge", "c5.2xlarge"].

### spec.computeResources.allocationStrategy

`string`

How Batch picks instance types (and Spot pools) when scaling out.
EC2/SPOT only; AWS defaults to BEST_FIT when omitted.
  BEST_FIT:                            cheapest fitting type only; may
                                       stall on capacity; no in-place
                                       infrastructure updates.
  BEST_FIT_PROGRESSIVE:                cheapest fitting types, falling
                                       forward when capacity runs out
                                       (recommended for EC2).
  BEST_FIT_PROGRESSIVE_ORDERED:        like BEST_FIT_PROGRESSIVE, but
                                       honors the instance_types list
                                       order as preference order.
  SPOT_CAPACITY_OPTIMIZED:             deepest Spot pools first --
                                       fewest interruptions (SPOT only).
  SPOT_PRICE_CAPACITY_OPTIMIZED:       balances Spot price and pool
                                       depth (recommended for SPOT).
  SPOT_CAPACITY_OPTIMIZED_PRIORITIZED: capacity-optimized, honoring the
                                       instance_types order (SPOT only).
Only BEST_FIT_PROGRESSIVE, SPOT_CAPACITY_OPTIMIZED, and
SPOT_PRICE_CAPACITY_OPTIMIZED support in-place infrastructure updates;
the others force replacement on most compute changes.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["BEST_FIT","BEST_FIT_PROGRESSIVE","BEST_FIT_PROGRESSIVE_ORDERED","SPOT_CAPACITY_OPTIMIZED","SPOT_PRICE_CAPACITY_OPTIMIZED","SPOT_CAPACITY_OPTIMIZED_PRIORITIZED"]}}

### spec.computeResources.instanceRole

`string | valueFrom`

The IAM instance profile applied to EC2/SPOT instances -- it wraps the
role that lets the ECS agent on each instance register with Batch's
underlying ECS cluster. Required for EC2 and SPOT. Reference an
AwsIamInstanceProfile's instance_profile_arn output or pass a literal
profile ARN.

- references: AwsIamInstanceProfile (`status.outputs.instance_profile_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamInstanceProfile, name: <that resource's name>, fieldPath: status.outputs.instance_profile_arn}} -- a bare string does not parse

### spec.computeResources.ec2KeyPair

`string`

The EC2 key pair name for SSH access to instances. EC2/SPOT only.
Prefer SSM Session Manager (via the instance profile) over SSH; omit
this unless direct SSH is genuinely needed.

### spec.computeResources.bidPercentage

`int32` · optional (explicit presence)

The maximum Spot price as a percentage of the On-Demand price
(e.g. 60 = pay at most 60% of On-Demand). SPOT only. Omit to default
to 100% -- with capacity-optimized strategies the actual price is
usually far below the cap anyway.

- rule: {"int32":{"lte":100,"gte":0}}

### spec.computeResources.spotIamFleetRole

`string | valueFrom`

The IAM role for the Amazon EC2 Spot Fleet that AWS Batch uses under
the BEST_FIT allocation strategy. AWS requires it ONLY for SPOT
environments using BEST_FIT (or no strategy, which defaults to
BEST_FIT); the modern capacity-optimized strategies do not use Spot
Fleet and need no role. CREATE-TIME ONLY: changing it replaces the
environment.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.computeResources.launchTemplate

`AwsBatchLaunchTemplate`

A custom EC2 launch template for instances -- custom AMIs, user data,
extra volumes, IMDSv2 posture. EC2/SPOT only. Adding or removing the
block replaces the environment; version changes update in place only
within the in-place-update envelope (see the spec comment).

### spec.computeResources.launchTemplate.launchTemplateId

`string | valueFrom` · required

The launch template. Reference an AwsLaunchTemplate's launch_template_id
output or pass a literal template ID ("lt-...").

- references: AwsLaunchTemplate (`status.outputs.launch_template_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLaunchTemplate, name: <that resource's name>, fieldPath: status.outputs.launch_template_id}} -- a bare string does not parse

### spec.computeResources.launchTemplate.version

`string`

The template version to launch: a version number, "$Latest", or
"$Default". Omit to use the template's default version. Batch caches
the resolved version at scale-out time; within the in-place-update
envelope, changing this triggers an infrastructure update that rolls
instances per update_policy.

### spec.computeResources.ec2Configurations

`[]AwsBatchEc2Configuration`

AMI selection for EC2/SPOT instances, keyed by image type. Maximum 2
entries (AWS allows one Linux and one Windows-family entry). CREATE-TIME
in practice: outside the in-place-update envelope any change here
replaces the environment.

- rule: {"repeated":{"maxItems":"2"}}

### spec.computeResources.ec2Configurations[].imageType

`string`

The image family Batch should launch. Common values: "ECS_AL2023"
(current Amazon Linux 2023 ECS AMI), "ECS_AL2" (Amazon Linux 2 --
required for GPU instance types together with ECS_AL2_NVIDIA),
"ECS_AL2_NVIDIA" (GPU), "EKS_AL2023"/"EKS_AL2" (Batch on EKS). AWS
defaults to ECS_AL2 when the whole block is omitted.

### spec.computeResources.ec2Configurations[].imageIdOverride

`string`

A specific AMI ID that overrides the image_type default -- the way to
pin a custom or hardened AMI while keeping Batch's image-family
semantics.

- rule: {"string":{"maxLen":"256"}}

### spec.computeResources.ec2Configurations[].imageKubernetesVersion

`string`

The EKS-optimized AMI's Kubernetes version, for Batch-on-EKS
environments (image types EKS_AL2023/EKS_AL2). Ignored for ECS image
types.

- rule: {"string":{"maxLen":"256"}}

### spec.computeResources.placementGroup

`string`

The EC2 placement group for tightly-coupled multi-node parallel jobs
that need low-latency networking between instances. EC2/SPOT only.
CREATE-TIME ONLY: changing it replaces the environment.

### spec.computeResources.resourceTags

`map<string, string>`

Tags applied to the launched compute resources themselves (EC2
instances and Spot Fleet requests) -- these propagate to the EC2
console, cost reports, and IAM tag conditions. EC2/SPOT only; Fargate
task ENIs cannot be tagged this way. Distinct from the identity tags
Planton applies to the compute environment resource.

### spec.eksConfiguration

`AwsBatchEksConfiguration`

Attach this compute environment to an EKS cluster so Batch schedules
jobs as Kubernetes pods instead of ECS tasks. CREATE-TIME ONLY: both
fields replace the environment when changed. The referenced cluster must
exist before the environment is created, and Batch-on-EKS job workloads
are registered outside this graph (the AwsBatchJobDefinition kind models
ECS-based container jobs).

### spec.eksConfiguration.eksClusterArn

`string | valueFrom` · required

The EKS cluster that receives Batch-managed nodes. Reference an
AwsEksCluster's cluster_arn output or pass a literal cluster ARN.

- references: AwsEksCluster (`status.outputs.cluster_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsEksCluster, name: <that resource's name>, fieldPath: status.outputs.cluster_arn}} -- a bare string does not parse

### spec.eksConfiguration.kubernetesNamespace

`string` · required

The Kubernetes namespace Batch launches job pods into. The namespace
must exist in the cluster and be RBAC-configured for Batch before the
environment is created.

- rule: {"required":true}

### spec.updatePolicy

`AwsBatchUpdatePolicy`

How infrastructure updates treat RUNNING jobs when Batch replaces
instances during an in-place update (EC2/SPOT environments). When unset,
Batch waits for jobs to finish on the old instances (up to 30 minutes)
before terminating them.

### spec.updatePolicy.terminateJobsOnUpdate

`bool`

Terminate running jobs when their instance is replaced, instead of
waiting for them to finish. Jobs are restarted per their retry
strategy. Leave false for long jobs that checkpoint poorly.

### spec.updatePolicy.jobExecutionTimeoutMinutes

`int32` · optional (explicit presence)

How long (in minutes, 1-360) Batch waits for running jobs to finish on
old instances before terminating them anyway. Only meaningful when
terminate_jobs_on_update is false. AWS default: 30.

- rule: {"int32":{"lte":360,"gte":1}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsBatchComputeEnvironment, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.compute_environment_arn` | `string` | The Amazon Resource Name (ARN) of the compute environment -- what job queues reference in their compute_environment_order. |
| `status.outputs.compute_environment_name` | `string` | The compute environment's name (derived from metadata.name). |
| `status.outputs.ecs_cluster_arn` | `string` | The ARN of the ECS cluster AWS Batch provisions behind a MANAGED compute environment -- useful for monitoring and debugging the tasks Batch actually runs. |
| `status.outputs.status` | `string` | The environment's current status (e.g. "VALID", "INVALID"). A queue can only associate environments whose status is VALID. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.serviceRole` | AwsIamRole | `status.outputs.role_arn` |
| `spec.computeResources.subnetIds` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.computeResources.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.computeResources.instanceRole` | AwsIamInstanceProfile | `status.outputs.instance_profile_arn` |
| `spec.computeResources.spotIamFleetRole` | AwsIamRole | `status.outputs.role_arn` |
| `spec.computeResources.launchTemplate.launchTemplateId` | AwsLaunchTemplate | `status.outputs.launch_template_id` |
| `spec.eksConfiguration.eksClusterArn` | AwsEksCluster | `status.outputs.cluster_arn` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsBatchJobQueue | `spec.computeEnvironmentOrder[].computeEnvironment` | `status.outputs.compute_environment_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
