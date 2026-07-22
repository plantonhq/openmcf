# Overview

The AwsAutoScalingGroup API resource provisions an EC2 Auto Scaling group:
the fleet manager that keeps instances launched from a launch template at
the desired size, replaces unhealthy members, spreads capacity across
subnets, and scales on policies and schedules.

## Why We Created This API Resource

The auto-scaling group is where self-managed EC2 compute becomes a fleet.
Modeling it as a pure orchestrator -- what to launch lives in the
referenced `AwsLaunchTemplate`, where traffic comes from lives in the
referenced `AwsLbTargetGroup` nodes -- lets you:

- **Separate blueprint from fleet**: rotate an AMI by updating the launch
  template; the group's instance refresh rolls the fleet with health-bound
  waves, checkpoints, and automatic rollback.
- **Blend purchase options honestly**: a mixed-instances policy expresses
  an On-Demand base plus a Spot majority drawn from many instance pools
  (explicit types or attribute-based requirements) -- the cost
  architecture real fleets run.
- **Wire the fleet into the load-balancing graph**: target-group
  references register instances on launch and drain them on termination,
  with `ELB` health checks replacing wedged-but-running processes.

## Key Features

### Capacity and Placement

- **Bounds and units**: min/max/desired, with `desired_capacity_type`
  counting in instances, vCPUs, or memory for heterogeneous fleets.
- **Multi-AZ spread** across referenced subnets with balanced-only or
  best-effort zone distribution and automatic rebalancing.
- **Warm pools**: pre-initialized (stopped, running, or hibernated)
  instances that cut scale-out latency from minutes to seconds.

### Scaling

- **Target tracking** (predefined group metrics, custom CloudWatch
  metrics, or metric-math expressions -- e.g. queue backlog per
  instance), **step scaling** against alarm breaches, legacy **simple
  scaling**, and **predictive scaling** on forecast load.
- **Scheduled actions**: cron or one-shot capacity changes with time
  zones -- business-hours scale-up, overnight scale-down.
- **Capacity rebalance**: proactively replace Spot instances at elevated
  interruption risk.

### Lifecycle and Rollouts

- **Instance refresh**: rolling replacement with min/max-healthy bounds,
  surge (launch-before-terminate), checkpointed canary waves, CloudWatch
  alarm watch, and auto-rollback.
- **Lifecycle hooks**: pause launching or terminating instances for
  warm-up or drain logic, delivered to an SNS topic reference.
- **Fleet hygiene**: max instance lifetime, termination policies,
  scale-in protection, and process suspension for maintenance windows.

## Benefits

- **Composability**: the group references its launch template, subnets,
  target groups, alarms, and SNS topics through `valueFrom`, so the
  architecture graph shows exactly how compute, routing, and
  observability connect.
- **Zero-downtime change**: template-driven instance refresh makes fleet
  rollouts reviewable, bounded, and reversible.
- **Consistency**: identical behavior across Terraform and Pulumi.

## Stack outputs

- `autoscaling_group_name`: group name (CloudWatch dimensions, ECS capacity providers, CLI)
- `autoscaling_group_arn`: group ARN (IAM policies, EventBridge rules)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
