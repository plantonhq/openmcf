---
title: "ECS Service HTTP"
description: "This preset creates an HTTP target group with the `ip` target type -- the shape ECS awsvpc tasks and Kubernetes pod-IP integrations register into. The group stays empty at deploy time; the ECS..."
type: "preset"
rank: "01"
presetSlug: "01-ecs-service-http"
componentSlug: "lb-target-group"
componentTitle: "LB Target Group"
provider: "aws"
icon: "package"
order: 1
---

# ECS Service HTTP

This preset creates an HTTP target group with the `ip` target type -- the
shape ECS awsvpc tasks and Kubernetes pod-IP integrations register into. The
group stays empty at deploy time; the ECS service (or controller) registers
and deregisters task IPs as deployments roll. Reference this group's
`target_group_arn` output from an `AwsLbListener` forward action, an
`AwsLbListenerRule`, or an `AwsEcsService`.

## When to Use

- ECS services in awsvpc networking mode behind an ALB
- EKS or self-managed Kubernetes services using pod-IP target registration
- Any containerized HTTP backend whose targets come and go with deployments

## Key Configuration Choices

- **`targetType: ip`** -- tasks register by IP, not instance ID; required for
  awsvpc tasks, which have no instance identity
- **Health check on `/healthz`** -- an application readiness endpoint with a
  strict `matcher: "200"`; a 15-second interval with a threshold of 3 marks
  targets healthy or unhealthy within ~45 seconds
- **`deregistrationDelaySeconds: 60`** -- short-lived HTTP requests do not
  need the 300-second AWS default; a shorter drain speeds up every rollout
- **No static `targets`** -- the ECS service owns registration; listing
  targets here would fight the scheduler

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<service-name>` | Name for the target group (max 32 chars after truncation) | Your service's name (e.g., `api`, `checkout`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<vpc-resource-name>` | Name of the AwsVpc resource the tasks run in | Your AwsVpc manifest's `metadata.name` |

## Common Additions

- Set `protocolVersion: GRPC` for gRPC services (enables gRPC status-code
  matchers on the health check)
- Set `slowStartSeconds` (30–900) if targets need cache warm-up before full
  load
- Set `stickiness` (`type: lb_cookie`) only if the service keeps per-client
  state -- it skews load distribution and conflicts with slow start

## Related Presets

- **02-nlb-tcp-passthrough** -- Layer-4 pass-through for non-HTTP protocols
- **03-lambda-function** -- invoke a Lambda function instead of addressing targets
