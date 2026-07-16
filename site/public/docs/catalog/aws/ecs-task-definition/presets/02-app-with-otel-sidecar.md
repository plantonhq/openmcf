---
title: "Application with OpenTelemetry Sidecar"
description: "This preset defines a two-container task: the application plus the AWS Distro for OpenTelemetry collector as a sidecar. The app waits for the collector to start, exports telemetry to it over..."
type: "preset"
rank: "02"
presetSlug: "02-app-with-otel-sidecar"
componentSlug: "ecs-task-definition"
componentTitle: "ECS Task Definition"
provider: "aws"
icon: "package"
order: 2
---

# Application with OpenTelemetry Sidecar

This preset defines a two-container task: the application plus the AWS
Distro for OpenTelemetry collector as a sidecar. The app waits for the
collector to start, exports telemetry to it over localhost (containers in
one awsvpc task share a network namespace), and the collector restarts in
place if it crashes -- without cycling the whole task.

## When to Use

- Services shipping traces/metrics through OTLP without embedding an
  exporter endpoint per environment
- Any task that needs a sidecar pattern (log routers, proxies, agents) --
  swap the collector image for your sidecar

## Key Configuration Choices

- **Per-container sizing** -- the task's 1 vCPU / 2 GiB is subdivided
  (app 768/1536, collector 256/512) so a leaking sidecar cannot starve
  the application
- **`essential: false` on the sidecar** -- a collector crash never kills
  the app; paired with `restartPolicy.enabled` the sidecar heals in place
- **`dependsOn: START`** -- the app starts after the collector is up, so
  the first telemetry batch has somewhere to go
- **`localhost` wiring** -- awsvpc tasks share one network namespace;
  sidecars are reachable on localhost, no discovery needed

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<service-name>` | The task-definition family name | Your service's name |
| `<aws-region>` | AWS region code | Your deployment region |
| `<account-id>` | Your AWS account ID | AWS console |
| `<repository>:<tag>` | The ECR repository and image tag | Your ECR console |
| `<execution-role-resource-name>` | Name of the AwsIamRole for the ECS agent | Your role manifest's `metadata.name` |
| `<task-role-resource-name>` | Name of the AwsIamRole for the application | Your role manifest's `metadata.name` |

## Common Additions

- `condition: HEALTHY` on the dependency (give the sidecar its own
  `healthCheck`) when the app must not start until the collector is
  verifiably ready
- A FireLens log router as a third container for non-CloudWatch log
  destinations
