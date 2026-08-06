# AWS ECS Service: Pure Scheduling on a Composed Graph

## What This Component Is

An ECS service keeps N copies of a task definition running: it starts
tasks, replaces unhealthy ones, drains and rolls them on deployments, and
registers their IPs into load-balancer target groups. `AwsEcsService`
models exactly that scheduling role and nothing else -- every resource
the service touches is a first-class node it references:

- the **task definition** (`AwsEcsTaskDefinition`) -- WHAT runs
- the **cluster** (`AwsEcsCluster`) -- WHERE it runs
- **subnets/security groups** -- the task network
- **target groups** (`AwsLbTargetGroup`) -- where traffic arrives,
  routed by `AwsLbListener`/`AwsLbListenerRule` nodes

The service registers task IPs into a referenced target group; it never
creates the group, the listener, or the routing rule. That split is what
makes multiple services share one ALB, blue/green pairs swap listener
rules between groups, and request-count autoscaling compose from the two
`arn_suffix` outputs.

## Deploys Travel Through the Graph

The `taskDefinition` reference points at the task definition's
revision-carrying ARN output. Registering a new revision (a new image
tag) changes that output, which changes this service's resolved input,
which rolls the service on its next deployment. There is no deploy
pipeline to configure -- the resource graph IS the pipeline, and the
deployment guards below decide whether the roll survives.

## Deployment Guards, In Layers

1. **Circuit breaker** -- tasks that repeatedly fail to reach steady
   state stop the rollout; `rollback: true` reverts to the last healthy
   deployment. Zero configuration; every production service should
   carry it.
2. **Alarm gating** -- referenced `AwsCloudwatchAlarm` nodes (by name --
   the CloudWatch API keys alarms on names) fail an in-flight deployment
   on error-rate/latency regressions the circuit breaker cannot see.
3. **Native blue/green** -- `deploymentConfiguration.strategy:
   BLUE_GREEN` with a target-group pair per load-balancer entry
   (`advancedConfiguration`), optional canary/linear traffic shifting,
   bake time, and Lambda lifecycle hooks at seven stages. ECS itself
   swaps the referenced production listener rule -- no CodeDeploy.

## Capacity Is a Spectrum

`launchType` and `capacityProviderStrategy` are mutually exclusive by
CEL: name one launch type, or blend providers with base/weight (the
FARGATE base-1 + FARGATE_SPOT weight-4 pattern keeps one guaranteed
on-demand task and runs ~80% of scaled capacity on Spot). With neither
set, both modules pin FARGATE explicitly rather than inheriting the
cluster default -- deployed results never depend on cluster-side mutable
state.

## Folded Autoscaling

Application Auto Scaling target + target-tracking policies are folded
into this spec because the scaler's identity IS the service
(`service/<cluster>/<service>`); nothing else can reference it. Three
policies: CPU, memory, and ALB requests-per-target -- the last scoped by
`<alb-arn-suffix>/<tg-arn-suffix>`, composed from the referenced nodes'
outputs. Both modules seed `desired_count` once and then ignore it
(lifecycle ignore on Terraform, `IgnoreChanges` on Pulumi), so the
scaler and operators own the live count.

## Deliberately Not Modeled

Bounded by the 90/10 rule; each skip is additive later if real
architectures pull for it:

- **`vpc_lattice_configurations`** -- VPC Lattice has no Planton kinds
  yet; modeling one arm of it would be a dead-end reference.
- **Classic-ELB wiring (`elb_name`, `iam_role`)** -- Classic Load
  Balancers are legacy; the LB family models ALB/NLB.
- **Apply-time knobs (`wait_for_steady_state`, `force_new_deployment`,
  `triggers`, `sigint_rollback`)** -- provisioner-session behavior, not
  durable resource state; a declarative spec would re-apply them
  meaninglessly.
- **Cloud Map namespace/service kinds** -- Service Connect's `namespace`
  and the legacy `serviceRegistries.registryArn` stay literal strings
  until a Cloud Map family exists (demand-gated).
- **Per-service custom tags** -- custom user tags are a platform-wide
  concern, not per-component scope.

## Failure Modes Worth Knowing

- AWS rejects `CreateService` when a referenced target group is not yet
  associated with a load-balancer listener -- compose the listener (or
  rule) node in the same graph so FK ordering guarantees it.
- The grace period is rejected on services without load balancers; the
  spec's CEL catches it before AWS does.
- DAEMON scheduling requires EC2 (no container instances exist on
  Fargate) and ignores desired_count/autoscaling -- all CEL-enforced.
- Service Connect's `port_name` must name a port in the task
  definition's `portMappings`; the join is by name, and a typo surfaces
  only at deploy time.
