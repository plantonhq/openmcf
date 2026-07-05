---
title: "EC2-Backed Cluster"
description: "This preset adds EC2 capacity to an ECS cluster by wrapping a referenced auto-scaling group as a capacity provider. ECS's managed scaling drives the group's instance count from task demand -- you..."
type: "preset"
rank: "03"
presetSlug: "03-ec2-capacity"
componentSlug: "ecs-cluster"
componentTitle: "ECS Cluster"
provider: "aws"
icon: "package"
order: 3
---

# EC2-Backed Cluster

This preset adds EC2 capacity to an ECS cluster by wrapping a referenced auto-scaling group as a capacity provider. ECS's managed scaling drives the group's instance count from task demand -- you size the group's bounds, ECS turns instances on and off. The right shape for workloads that need GPUs, specific instance families, instance-store disks, or the lower unit price of Reserved/Spot EC2.

## When to Use

- Workloads requiring GPU or other specialized instance types Fargate does not offer
- High-density clusters where EC2 unit economics beat Fargate at sustained load
- Daemon-style workloads (log shippers, monitoring agents) that need per-instance placement

## Key Configuration Choices

- **Auto-scaling group by reference** (`ec2CapacityProviders[].autoScalingGroupArn`) -- Reference an `AwsAutoScalingGroup`'s `autoscaling_group_arn` output; the group's launch template decides the instance shape (use an ECS-optimized AMI whose agent joins this cluster via user data)
- **Managed scaling at 80%** (`targetCapacity: 80`) -- Keeps 20% headroom so new tasks place immediately instead of waiting for an instance launch; 100 packs instances fully
- **Managed draining** (`managedDraining: ENABLED`) -- Scale-in and instance refresh drain tasks gracefully instead of killing them
- **Fargate kept alongside** -- Services choose per-strategy: EC2 for the steady base, Fargate for burst or one-off tasks

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region where the cluster will be created (e.g., `us-west-2`) | AWS region list |
| `<auto-scaling-group-arn>` | ARN of the auto-scaling group providing the instances | `AwsAutoScalingGroup` status outputs |

## Related Presets

- **01-fargate-standard** -- Serverless-only baseline with no instances to manage
- **02-fargate-cost-optimized** -- Fargate + Fargate Spot blend without EC2
