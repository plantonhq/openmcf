---
title: "Hardened Monitored Deployment"
description: "This preset is the production posture: secure egress through the workspace's managed network, Application Insights on, a patient startup probe for slow-loading models, honest request limits, and..."
type: "preset"
rank: "02"
presetSlug: "02-hardened-monitored-deployment"
componentSlug: "machine-learning-online-deployment"
componentTitle: "Machine Learning Online Deployment"
provider: "azure"
icon: "package"
order: 2
---

# Hardened Monitored Deployment

This preset is the production posture: secure egress through the workspace's managed network, Application Insights on, a patient startup probe for slow-loading models, honest request limits, and sampled model data collection for drift monitoring.

## When to Use

- Production model serving in locked-down workspaces
- Models whose load time exceeds the default probe patience
- Deployments feeding a drift-monitoring pipeline

## Key Configuration Choices

- **`egressPublicNetworkAccessEnabled: false`** -- image and model pulls traverse the workspace's managed network only; prove the registry's and storage's private endpoints exist BEFORE deploying, or provisioning fails with what looks like a pull error.
- **`startupProbe`** -- sixty failures at ten-second periods gives a model ten minutes to load before the service declares it dead; the startup probe gates liveness and readiness.
- **`requestSettings`** -- the timeout raised above the clients' own; concurrency of 2 only because the container is known to parallelize -- the default 1 is correct for single-threaded scorers.
- **`dataCollector`** -- 20% sampling keeps storage growth honest on busy endpoints; payloads inherit the workspace storage's access controls, so review who can read them.

## After Deployment

Mirror traffic to `green` first (`mirrorTraffic: {green: 20}`), watch logs and latency, then shift the live map.
