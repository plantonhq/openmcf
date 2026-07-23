---
title: "Web Service Deployment"
description: "This preset deploys a single-replica web application fronted by a ClusterIP Service, with a readiness probe so traffic only reaches pods that are ready. It is the most common Kubernetes Deployment..."
type: "preset"
rank: "01"
presetSlug: "01-web-service"
componentSlug: "deployment"
componentTitle: "Deployment"
provider: "kubernetes"
icon: "package"
order: 1
---

# Web Service Deployment

This preset deploys a single-replica web application fronted by a ClusterIP Service, with a readiness probe so traffic only reaches pods that are ready. It is the most common Kubernetes Deployment pattern.

External exposure is composed, not embedded: attach a Gateway API route kind such as `KubernetesHttpRoute` that references this workload's exported `service` output to publish it at a hostname. That keeps every piece of exposure infrastructure visible in your resource graph.

## When to Use

- Standard web applications or REST APIs
- Development or low-traffic production services where a single replica is sufficient

## Key Configuration Choices

- **Single replica** (the default) — no autoscaling; see `02-web-service-with-hpa` for scaling
- **Port 8080 container / port 80 service** — the service port faces clients, the container port faces the app
- **Readiness probe on /healthz** — the piece that makes rolling updates zero-downtime; point it at your app's health endpoint
- **Default resources** — 50m–1000m CPU, 100Mi–1Gi memory

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace for the deployment | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-image>` | Container image repository (e.g., `ghcr.io/org/app`) | Your container registry |
| `<your-image-tag>` | Image tag or version (e.g., `v1.2.3`) | Your CI/CD pipeline output |

## Related Presets

- **02-web-service-with-hpa** — Adds horizontal pod autoscaling, zero-downtime rollout strategy, and a pod disruption budget
- **03-worker** — Background worker without a Service
- **04-hardened-production** — Restricted-profile security hardening with topology spread
