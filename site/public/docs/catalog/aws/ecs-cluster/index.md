---
title: "ECS Cluster"
description: "ECS Cluster deployment documentation"
icon: "package"
order: 100
componentName: "awsecscluster"
---

# AWS ECS Cluster

Deploys an Amazon ECS cluster through one declarative manifest: the scheduling boundary for services and tasks, its capacity (Fargate built-ins, EC2 capacity providers wrapping auto-scaling groups, or a blend), and cluster-wide posture -- Container Insights, ECS Exec auditing, Fargate storage encryption, and the Service Connect default namespace.

## What Gets Created

When you deploy an AwsEcsCluster resource, Planton provisions:

- **ECS Cluster** — keyed by `metadata.name`, with Container Insights, exec auditing, managed-storage encryption, and Service Connect defaults
- **EC2 Capacity Providers** — from `ec2CapacityProviders`, one per entry, each wrapping a referenced auto-scaling group with ECS-managed scaling and draining
- **Capacity Provider Association** — one association putting the union of Fargate built-ins and EC2 provider names onto the cluster, together with the default strategy

The cluster itself is free -- only the tasks and instances it schedules cost money. All resources are tagged with Planton metadata (organization, environment, resource kind, resource ID).

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **An auto-scaling group** (only for EC2 capacity) — reference an `AwsAutoScalingGroup` whose launch template uses an ECS-optimized AMI joining this cluster via user data
- **A KMS key** (optional) — reference an `AwsKmsKey` for exec-session or Fargate storage encryption
- **A Cloud Map namespace** (optional) — for Service Connect defaults

## Quick Start

Create a file `cluster.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEcsCluster
metadata:
  name: my-cluster
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
```

Deploy:

```shell
planton apply -f cluster.yaml
```

This creates a serverless cluster where one task is guaranteed On-Demand and ~80% of scaled capacity rides Fargate Spot.

## Configuration Reference

### Observability

| Field | Type | Description |
|-------|------|-------------|
| `containerInsights` | `string` | `enabled` (task/service metrics), `enhanced` (container-level observability with automatic dashboards -- the production posture), or `disabled`. Unset keeps the account default. |

### Capacity

| Field | Type | Description |
|-------|------|-------------|
| `capacityProviders` | `string[]` | The AWS-managed serverless built-ins to associate: `FARGATE` and/or `FARGATE_SPOT`. EC2 capacity is defined separately below; both sets associate onto the cluster together. |
| `ec2CapacityProviders` | `object[]` | Per-name EC2 capacity providers. Each entry: `name` (what services put in a strategy; may not start with `aws`/`ecs`/`fargate`), `autoScalingGroupArn` (reference an `AwsAutoScalingGroup`'s `autoscaling_group_arn` output), `managedScaling`, `managedTerminationProtection`, `managedDraining`. |
| `ec2CapacityProviders[].managedScaling` | `object` | ECS drives the group's desired count: `status` (`ENABLED`/`DISABLED`), `targetCapacity` (1-100% utilization target -- 80 keeps headroom for instant placement), `minimumScalingStepSize`/`maximumScalingStepSize`, `instanceWarmupPeriodSeconds`. |
| `ec2CapacityProviders[].managedTerminationProtection` | `string` | `ENABLED` protects instances running tasks from scale-in -- requires the group's own new-instance scale-in protection. |
| `ec2CapacityProviders[].managedDraining` | `string` | `ENABLED` (AWS default) drains tasks gracefully off terminating instances. |
| `defaultCapacityProviderStrategy` | `object[]` | What ECS uses when a service names no strategy: `capacityProvider` (any associated name), `base` (guaranteed tasks; only one entry may set non-zero), `weight` (relative share beyond bases). Validated against the associated set. |

### Cluster-Wide Posture

| Field | Type | Description |
|-------|------|-------------|
| `executeCommandConfiguration` | `object` | ECS Exec auditing: `logging` (`DEFAULT` = task's own log config, `OVERRIDE` = the explicit destinations below, `NONE` = unaudited), `logConfiguration` (CloudWatch group and/or S3 bucket, with encryption requirements), `kmsKeyId` (ref -- encrypts session traffic). |
| `managedStorageConfiguration` | `object` | Customer-managed KMS keys: `fargateEphemeralStorageKmsKeyId` (the key policy must grant the Fargate service principal) and `kmsKeyId`. |
| `serviceConnectNamespaceArn` | `string` | The Cloud Map namespace ARN Service Connect uses by default for services in this cluster. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `cluster_name` | The cluster name (mirrors `metadata.name`) |
| `cluster_arn` | The join key — `AwsEcsService.cluster_arn` references it |
| `capacity_provider_names` | The full strategy vocabulary: built-ins plus folded EC2 provider names |
| `capacity_provider_arns` | The EC2 capacity providers' ARNs (empty for Fargate-only clusters) |

## Related Resources

- [AWS ECS Service](/docs/catalog/aws/ecs-service) — the long-running workloads scheduled into this cluster
- [AWS ECS Task Definition](/docs/catalog/aws/ecs-task-definition) — the container blueprints those services run
- [AWS Auto Scaling Group](/docs/catalog/aws/auto-scaling-group) — the instance fleet behind an EC2 capacity provider
- [AWS Launch Template](/docs/catalog/aws/launch-template) — the instance shape that fleet launches from
- [AWS KMS Key](/docs/catalog/aws/kms-key) — exec-session and Fargate storage encryption
