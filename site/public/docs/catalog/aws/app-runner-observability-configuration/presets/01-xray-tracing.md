---
title: "X-Ray Tracing"
description: "The standard observability configuration: distributed request tracing to AWS X-Ray, shared by every service that references it."
type: "preset"
rank: "01"
presetSlug: "01-xray-tracing"
componentSlug: "app-runner-observability-configuration"
componentTitle: "App Runner Observability Configuration"
provider: "aws"
icon: "package"
order: 1
---

# X-Ray Tracing

The standard observability configuration: distributed request tracing to AWS X-Ray, shared by every service that references it.

## When to Use

- Any production App Runner fleet that needs request-level tracing
- Debugging latency across services that call each other or downstream AWS APIs

## What It Configures

- **`vendor: AWSXRAY`** — each service instance runs an OpenTelemetry collector sidecar forwarding spans to X-Ray

## What to Customize

- Replace `<aws-region>` with your region
- Reference from services via `observabilityConfigurationArn` — the reference itself is the enable switch
- Instrument the application with the AWS Distro for OpenTelemetry (ADOT) SDK; without instrumentation the collector has nothing to forward
- The service's `instanceRoleArn` role needs X-Ray write permissions (`AWSXRayDaemonWriteAccess` covers it)
