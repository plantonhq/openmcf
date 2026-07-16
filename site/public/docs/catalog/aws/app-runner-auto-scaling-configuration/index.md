---
title: "App Runner Auto Scaling Configuration"
description: "App Runner Auto Scaling Configuration deployment documentation"
icon: "package"
order: 100
componentName: "awsapprunnerautoscalingconfiguration"
---

# AWS App Runner Auto Scaling Configuration

The reusable scaling policy for AWS App Runner services -- how many requests one instance absorbs before another launches, how far scale-out may go, and how many instances stay warm.

## Why a First-Class Resource

App Runner models scaling as a shared, versioned configuration referenced by ARN, and so does Planton: one configuration tunes an entire fleet from one place, and because every change registers a new revision with a new ARN, referencing services roll to the new posture through the resource graph -- no per-service edits.

## Key Capabilities

- **Concurrency-based scaling** -- the `maxConcurrency` dial (1-200) trades per-instance headroom against instance count.
- **Warm floor** -- `minSize` instances stay provisioned (memory-only billing while idle), eliminating cold starts for latency-sensitive APIs.
- **Cost ceiling** -- `maxSize` bounds scale-out during spikes.
- **Revision semantics** -- values are create-time immutable; each change registers the next revision under the same name and exports a new revision-carrying ARN.

## Composes With

- `AwsAppRunnerService` -- adopts this configuration via `autoScalingConfigurationArn`.
