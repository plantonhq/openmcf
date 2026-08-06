# AWS Auto Scaling Group: The Fleet Orchestrator

## What an Auto Scaling Group Is

An EC2 Auto Scaling group maintains a set of instances at a desired size:
it launches from a template, spreads capacity across availability zones,
replaces members that fail health checks, registers instances into load
balancer target groups, and adjusts the size in response to policies,
schedules, and forecasts. It is the resource that turns "an instance" into
"a fleet."

`AwsAutoScalingGroup` models the group as a pure ORCHESTRATOR. What to
launch lives in the referenced `AwsLaunchTemplate` (AMI, instance type,
storage, IAM identity, metadata posture); where traffic comes from lives in
the referenced `AwsLbTargetGroup` nodes. The group owns how many, where,
and when: capacity bounds, subnet spread, purchase-option mix, health
model, scaling behavior, and instance lifecycle. That separation is the
composition story -- one hardened template backs many groups, and one
group serves many target groups.

## Why the Satellites Fold In

AWS models scaling policies, scheduled actions, lifecycle hooks, and
notifications as standalone resources, and most IaC tools mirror that.
This component folds all four into the group's spec, because each is a
sub-resource of exactly ONE group, is referenced by nothing else in the
graph, and has no meaning in isolation -- a scaling policy cannot exist
without its group, and nothing composes against "a scheduled action."
Standalone kinds here would be a sea of glue nodes.

The fold is spec-level only: both IaC modules still manage each entry as
its own provider resource (`aws_autoscaling_policy`,
`aws_autoscaling_schedule`, `aws_autoscaling_lifecycle_hook`,
`aws_autoscaling_notification`), keyed by name -- so adding, editing, or
removing one is an in-place update that never touches the group itself.
Lifecycle hooks deliberately use the standalone hook resource rather than
the group's `initial_lifecycle_hook` block, which is create-only: inline
hooks would force group replacement on every hook edit.

One satellite stays out on purpose: the CloudWatch alarms that TRIGGER
step and simple policies are real, referenceable, first-class resources --
they are `AwsCloudwatchAlarm` nodes whose alarm actions target the
policy's ARN. The graph keeps observability and scaling composed, not
bundled.

## Launch Template XOR Mixed Instances

The group launches in exactly one of two modes, mirroring AWS's own
ExactlyOneOf:

- **`launchTemplate`** -- a single template reference with an optional
  version pin. Leaving the version unset follows the template's default
  version (`$Default`), which is what lets a template update roll the
  fleet; pinning a numeric version freezes the fleet until the pin moves.
- **`mixedInstancesPolicy`** -- the cost architecture real fleets run: an
  On-Demand base (`onDemandBaseCapacity`) for guaranteed capacity, a
  Spot majority above it (`onDemandPercentageAboveBaseCapacity: 0` means
  all-Spot above the base -- the field is optional precisely so explicit
  zero is expressible), drawn from several instance pools. Overrides add
  pools by explicit type, by a different template (an arm64 template for
  Graviton types), or by attribute-based `instanceRequirements`;
  `price-capacity-optimized` is the AWS-recommended Spot strategy, and
  `capacityRebalance` proactively replaces at-risk Spot instances.

Weighted capacity makes heterogeneous sizes count fairly (an
`m5.2xlarge` at weight 4 next to an `m5.large` at weight 1), and
`desiredCapacityType` switches the counting unit to vCPUs or memory when
sizes vary -- weights are strings at AWS but honest integers here.

## Health, Rollouts, and the Refresh Model

Two decisions dominate fleet reliability:

- **`healthCheckType: ELB`** whenever the group serves target groups.
  With the default `EC2` checks, a wedged-but-running process passes
  status checks forever; `ELB` trusts the target group's application
  health check and replaces what it fails.
- **Instance refresh** turns template changes into bounded rollouts:
  `minHealthyPercentage`/`maxHealthyPercentage` set the dip and surge
  (100/110 means launch-before-terminate and never dip below full
  strength), `checkpointPercentages` stage canary waves,
  `alarms` (references to `AwsCloudwatchAlarm` nodes) fail the refresh on
  regression, and `autoRollback` returns the fleet to the previous
  configuration. `instanceMaintenancePolicy` applies the same bounds to
  EVERY replacement operation, not just refreshes.

Scheduled actions express absent capacity values as "leave unchanged"
(AWS's -1 convention); the spec models them as optional integers so
explicit zero -- scale to nothing overnight -- is distinguishable from
unset.

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **`launch_configuration`** -- the legacy predecessor of launch
  templates; AWS itself steers new fleets off it. This component is
  launch-template-only.
- **Classic ELB attachment (`load_balancers`)** -- Classic Load Balancers
  are legacy; the target-group references cover ALB/NLB.
- **`traffic_source` (VPC Lattice arm)** -- the unified attachment block
  matters only for VPC Lattice, which has no Planton kind yet; the
  type-specific `target_group_arns` covers the real case today.
- **Capacity reservation targeting** -- deferred until a capacity
  reservation kind exists to reference.
- **`context`** -- an AWS reserved field with no public semantics.
- **Predictive scaling's customized metric-data-query specifications** --
  the predefined metric pairs (`ASGCPUUtilization`, `ASGNetworkIn/Out`,
  `ALBRequestCount`) cover the real predictive use cases; the raw
  query form triples the surface for a niche need. Target tracking DOES
  carry the full customized/metric-math form, where backlog-per-instance
  style metrics are common.
- **`ignore_failed_scaling_activities`, `force_delete_warm_pool`,
  `min_elb_capacity`/`wait_for_elb_capacity`** -- engine-behavior escape
  hatches for pathological states; `wait_for_capacity_timeout` and
  `force_delete` cover the operational needs.

## Immutability and Naming

Only the group name (from `metadata.name`; AWS limit 255 characters) is
create-only -- everything else updates in place, which is exactly why the
group is safe to enrich over time. Deletion drains instances unless
`forceDelete` is set.

## Dual-Engine Implementation

`AwsAutoScalingGroup` ships both a Terraform/OpenTofu module and a Pulumi
(Go) module at behavioral parity. Both express the identity tags through
the group's native key/value/propagate-at-launch triple (so launched
instances inherit them), materialize the four folded satellites as
per-name provider resources, convert the same honest integers to AWS's
string-typed fields (weights, warmup, buffers), honor presence semantics
on the optional zeros, and export the same outputs
(`autoscaling_group_name`, `autoscaling_group_arn`). Whichever engine a
team standardizes on, the fleet behaves identically.
