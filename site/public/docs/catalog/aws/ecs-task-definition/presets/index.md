---
title: "Presets"
description: "Ready-to-deploy configuration presets for ECS Task Definition"
type: "preset-list"
componentSlug: "ecs-task-definition"
componentTitle: "ECS Task Definition"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-web-app"
    rank: "01"
    title: "Web Application"
    excerpt: "This preset defines a Fargate web application task: one container with a named HTTP port (ready for load-balancer and Service Connect wiring), secrets injected by reference at task start, an..."
  - slug: "02-app-with-otel-sidecar"
    rank: "02"
    title: "Application with OpenTelemetry Sidecar"
    excerpt: "This preset defines a two-container task: the application plus the AWS Distro for OpenTelemetry collector as a sidecar. The app waits for the collector to start, exports telemetry to it over..."
  - slug: "03-arm64-worker"
    rank: "03"
    title: "ARM64 Background Worker"
    excerpt: "This preset defines a queue-consuming background worker on Graviton (ARM64) Fargate: no ports, no load balancer, the smallest Fargate task size, and a long stop timeout so in-flight work drains..."
  - slug: "04-launch-time-ebs-volume"
    rank: "04"
    title: "Launch-Time EBS Volume Task"
    excerpt: "This preset declares a task whose volume is configured at launch time: the task definition names the volume slot, and the ECS service that runs it attaches a fresh, service-owned EBS volume per task..."
---

# ECS Task Definition Presets

Ready-to-deploy configuration presets for ECS Task Definition. Each preset is a complete manifest you can copy, customize, and deploy.
