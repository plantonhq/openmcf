---
title: "Routed Endpoint with Identity"
description: "This preset creates a batch endpoint with a system-assigned identity and an explicit default-deployment pointer -- the shape for teams whose role-grant conventions target the endpoint object and..."
type: "preset"
rank: "02"
presetSlug: "02-routed-endpoint-with-identity"
componentSlug: "machine-learning-batch-endpoint"
componentTitle: "Machine Learning Batch Endpoint"
provider: "azure"
icon: "package"
order: 2
---

# Routed Endpoint with Identity

This preset creates a batch endpoint with a system-assigned identity and an explicit default-deployment pointer -- the shape for teams whose role-grant conventions target the endpoint object and whose rollout process moves the pointer between recipe versions.

## When to Use

- Organizations that grant Azure roles to endpoint identities as a convention
- Endpoints whose submissions should route to a well-known default recipe
- Rollout processes that move the default pointer between deployment versions (`production`, `production-v2`, ...)

## Key Configuration Choices

- **`identity: SYSTEM_ASSIGNED`** -- created with the endpoint. Note the batch data path does NOT use it (jobs run under the submitter's token plus the compute pool's identity); set it for grant conventions, not for job execution.
- **`defaultDeploymentName: production`** -- names the deployment that answers unrouted submissions. It usually does not exist yet at endpoint creation -- ARM accepts the pointer, and submissions fail until the deployment attaches; create the recipe promptly or set the pointer afterwards.

## After Deployment

Attach an **Azure Machine Learning Batch Deployment** named `production` wiring `endpointId` by reference; unrouted job submissions then run that recipe.
