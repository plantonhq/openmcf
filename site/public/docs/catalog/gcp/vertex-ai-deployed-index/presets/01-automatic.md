---
title: "Automatic Serving Compute"
description: "Deploys an index onto an endpoint with Vertex-managed serving compute — the zero-configuration path. Vertex AI picks the machine types; replicas scale between the bounds you set."
type: "preset"
rank: "01"
presetSlug: "01-automatic"
componentSlug: "vertex-ai-deployed-index"
componentTitle: "Vertex AI Deployed Index"
provider: "gcp"
icon: "package"
order: 1
---

# Automatic Serving Compute

Deploys an index onto an endpoint with Vertex-managed serving compute —
the zero-configuration path. Vertex AI picks the machine types; replicas
scale between the bounds you set.

## What this preset creates

A DeployedIndex named `products_v1` placing the `product-embeddings`
index onto the `vector-serving` endpoint in `us-central1`, with
`automaticResources` scaling between 2 and 5 replicas. Deploying takes
tens of minutes; the endpoint serves queries only after `indexSyncTime`
catches up with the index's `updateTime`.

## When to use

- First deployment of an index, before capacity needs are known
- Workloads with variable query volume where GCP-managed scaling is
  preferable to manual capacity planning
- Any deployment where machine-type control is not a requirement

## Remix ideas

- Raise `maxReplicaCount` any time without a redeploy — the replica
  bounds are the only in-place-updatable fields.
- Keep `minReplicaCount` at 2 or higher for availability; Vertex AI
  offers no SLA at 1 replica.
- Choose `deployedIndexId` like an API name (letter start,
  letters/numbers/underscores) — it is the query handle and immutable.
- Switch to `dedicatedResources` when you need predictable performance
  or a specific machine type — but that switch redeploys.
