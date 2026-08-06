# AWS ECS Cluster — Architecture and Design

## Overview

`AwsEcsCluster` models one Amazon ECS cluster at the full provider
surface: the capacity layer (AWS-managed Fargate built-ins plus folded
EC2 capacity providers wrapping referenced auto-scaling groups, unioned
onto the cluster by a single association), the default capacity provider
strategy, and cluster-wide posture -- Container Insights, ECS Exec
auditing, Fargate managed-storage encryption, and the Service Connect
default namespace.

A Fargate-capable cluster is a true leaf: nothing must exist before it,
and the cluster itself is free. EC2 capacity composes by reference --
the auto-scaling group (`AwsAutoScalingGroup`) owns the fleet's shape
through its launch template, and ECS's managed scaling drives the
group's desired count from task demand.

## Design Decisions

- **EC2 capacity providers FOLDED, materialized per-name.** The
  provider models a custom capacity provider as an account-level
  wrapper around exactly one auto-scaling group, associated to a
  cluster by a separate PUT-style resource. Services reference
  providers by NAME string; nothing references them by ARN; in
  practice each belongs to one cluster. That is the folded-satellite
  class: spec entries keyed by name, each its own provider resource in
  both engines (in-place edits), auto-joined to the association -- a
  provider is never listed twice.
- **One association resource, always.** `aws_ecs_cluster_capacity_providers`
  replaces the cluster's whole provider set and default strategy on
  every write. Both engines therefore create exactly one association
  over the union of built-ins and folded names; the built-ins list
  (`capacity_providers`) stays restricted to `FARGATE`/`FARGATE_SPOT`
  because those are the only providers that exist without being
  defined.
- **Strategy names validated against the associated set.** AWS rejects
  a default strategy naming an unassociated provider with an opaque
  runtime error; the spec's CEL rule fails in seconds and names the
  offending entry. AWS's one-non-zero-base rule is likewise CEL.
- **Exec configuration models what the API models.** Block presence is
  the audit configuration; `logging` carries the provider's own
  `DEFAULT`/`OVERRIDE`/`NONE` strings; `OVERRIDE` requires destinations
  and destinations require `OVERRIDE` (both directions CEL-enforced).
  There is no synthetic "exec disabled" sentinel -- exec availability
  is a per-service setting, not a cluster one.
- **Provider strings, not enums.** `container_insights` carries the
  API's `enabled`/`enhanced`/`disabled` values (the `enhanced`
  container-level tier is only expressible this way); unset keeps the
  account default rather than forcing one.
- **Managed termination protection stated honestly.** `ENABLED`
  requires the auto-scaling group's own new-instance scale-in
  protection -- a cross-resource coupling AWS validates at create; the
  field comment names it because CEL cannot see the referenced group.

## Deliberately Skipped Provider Surface

| Provider surface | Verdict | Reason |
|---|---|---|
| `managed_instances_provider` (ECS Managed Instances) | DEFER | A cluster-scoped provider where AWS also manages the instances (infrastructure role, launch-template-in-provider, attribute-based instance selection) -- a genuinely separate capacity product; revisit on concrete pull. |
| Standalone `AwsEcsCapacityProvider` kind | NOT BUILT | Fails the split test: referenced only by name strings from strategies, one-cluster in practice, meaningless without its group -- the folded per-name entries are the honest shape. |
| `setting` beyond `containerInsights` | N/A | `containerInsights` is the only setting name the ECS API defines today; the spec models the value directly instead of an open key/value list. |

## Billing Note

The cluster, its capacity providers, and the association are all free.
Container Insights bills per metric/log volume once tasks run
(`enhanced` more than `enabled`). The E2E scenarios launch zero tasks
and zero instances (the EC2-capacity chain runs a min=0/desired=0
group), so full lanes accrue effectively nothing.

## References

- [Amazon ECS clusters](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/clusters.html)
- [Auto Scaling group capacity providers](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/asg-capacity-providers.html)
- [Managed scaling](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/managed-scaling.html)
- [Container Insights with enhanced observability](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/ContainerInsights.html)
- [ECS Exec logging and auditing](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-exec.html#ecs-exec-logging)
