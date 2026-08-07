---
title: "Container Isolated"
description: "This preset scales on ONE container's CPU instead of the pod-level average. A plain `resource` metric blends every container in the pod — the app, the service-mesh proxy, the log shipper — into one..."
type: "preset"
rank: "02"
presetSlug: "02-container-isolated"
componentSlug: "horizontalpodautoscaler"
componentTitle: "HorizontalPodAutoscaler"
provider: "kubernetes"
icon: "package"
order: 2
---

# Container Isolated

This preset scales on ONE container's CPU instead of the pod-level average. A plain `resource` metric blends every container in the pod — the app, the service-mesh proxy, the log shipper — into one number, and that blend skews both ways: a hot sidecar can trigger scale-out while the app idles, and a heavy app can hide behind a large idle sidecar's requests. `container_resource` reads only the named container, so the scaling signal is the application's own load.

## When to Use

- Pods running service-mesh sidecars (Envoy/Istio proxies), log shippers, or other companions whose resource profile differs from the app's
- Any workload where "pod CPU" and "app CPU" have diverged enough to mis-scale
- Workloads whose sidecars have deliberately generous requests — the pod-level utilization denominator includes them, diluting the app's signal

## Key Configuration Choices

- **`container_resource` with `container`** — the average is computed over exactly this container across all the target's pods. The name must match the container name declared in the pod template; a wrong name reads as an unavailable metric and the autoscaler holds the current count
- **`utilization: 60`** — still measured against the NAMED container's requests, so that container must declare CPU requests
- **`min_replicas: 2` / `max_replicas: 10`** — the same floor/ceiling conversation as the basic preset: availability minimum and cost ceiling
- **Requires metrics-server only** — per-container readings come from the same resource-metrics pipeline; no custom adapter needed

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The scale target's own namespace — an HPA cannot scale across namespaces | Your namespace management |
| `<your-workload-name>` | The target Deployment's name | The workload's manifest |
| `<your-app-container-name>` | The app container's name as declared in the pod template | The workload's container spec |

## Related Presets

- **01-cpu-autoscale** — the pod-level average, fine for pods without significant sidecars
- **04-behavior-tuned** — add conservative scale-down on top of either CPU signal
