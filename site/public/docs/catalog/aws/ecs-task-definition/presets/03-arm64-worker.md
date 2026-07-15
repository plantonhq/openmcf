---
title: "ARM64 Background Worker"
description: "This preset defines a queue-consuming background worker on Graviton (ARM64) Fargate: no ports, no load balancer, the smallest Fargate task size, and a long stop timeout so in-flight work drains..."
type: "preset"
rank: "03"
presetSlug: "03-arm64-worker"
componentSlug: "ecs-task-definition"
componentTitle: "ECS Task Definition"
provider: "aws"
icon: "package"
order: 3
---

# ARM64 Background Worker

This preset defines a queue-consuming background worker on Graviton
(ARM64) Fargate: no ports, no load balancer, the smallest Fargate task
size, and a long stop timeout so in-flight work drains cleanly on
deployments and scale-in.

## When to Use

- SQS/queue consumers, schedulers, and batch-ish daemons with no inbound
  traffic
- Cost-sensitive fleets -- Fargate ARM pricing runs ~20% below x86 for
  the same vCPU/memory (images must be arm64-built)

## Key Configuration Choices

- **`cpuArchitecture: ARM64`** -- the Graviton discount with a one-line
  change; multi-arch images make this transparent
- **No `portMappings`** -- workers pull work; nothing routes to them
- **`stopTimeoutSeconds: 120`** -- the Fargate maximum grace period
  between SIGTERM and SIGKILL, so a message mid-flight finishes instead
  of being re-queued
- **Task role carries the queue permissions** -- the worker's AWS
  identity composes with `AwsIamRole` instead of embedding credentials

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<worker-name>` | The task-definition family name | Your worker's name |
| `<aws-region>` | AWS region code | Your deployment region |
| `<account-id>` | Your AWS account ID | AWS console |
| `<repository>:<tag>` | The ECR repository and image tag (arm64) | Your ECR console |
| `<execution-role-resource-name>` | Name of the AwsIamRole for the ECS agent | Your role manifest's `metadata.name` |
| `<task-role-resource-name>` | Name of the AwsIamRole for the worker | Your role manifest's `metadata.name` |
| `<sqs-queue-url>` | The queue the worker consumes | Your SQS console or AwsSqsQueue outputs |

## Common Additions

- A `healthCheck` probing the worker's liveness endpoint or a heartbeat
  file, so a wedged worker is replaced
- `ephemeralStorageGib` for workers that stage large payloads on disk
