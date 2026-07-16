---
title: "App Runner Observability Configuration"
description: "App Runner Observability Configuration deployment documentation"
icon: "package"
order: 100
componentName: "awsapprunnerobservabilityconfiguration"
---

# AWS App Runner Observability Configuration

The reusable tracing policy for AWS App Runner services -- reference it from a service and every instance forwards request traces to AWS X-Ray through a managed OpenTelemetry collector.

## Why a First-Class Resource

AWS models tracing as a shared, versioned configuration referenced by ARN, and so does Planton: one configuration turns on tracing for an entire fleet, and because each change registers a new revision with a new ARN, referencing services roll to the new posture through the resource graph.

## Key Capabilities

- **X-Ray distributed tracing** -- request spans flow to X-Ray for service maps and latency analysis; instrument the application with the ADOT SDK.
- **Reference-is-the-switch composition** -- a service enables tracing simply by referencing the configuration; removing the reference disables it. No boolean to drift out of sync.
- **Revision semantics** -- trace settings are create-time immutable; each change registers the next revision under the same name.

## Composes With

- `AwsAppRunnerService` -- adopts the configuration via `observabilityConfigurationArn`.
