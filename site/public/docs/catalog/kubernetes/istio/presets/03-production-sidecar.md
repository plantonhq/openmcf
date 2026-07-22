---
title: "Production Sidecar Mesh"
description: "This preset installs a hardened sidecar-mode control plane: highly available istiod, an organization-unique trust domain, egress lockdown, mesh-wide access logging, and the node-level CNI agent in..."
type: "preset"
rank: "03"
presetSlug: "03-production-sidecar"
componentSlug: "istio"
componentTitle: "Istio"
provider: "kubernetes"
icon: "package"
order: 3
---

# Production Sidecar Mesh

This preset installs a hardened sidecar-mode control plane: highly available istiod, an organization-unique trust domain, egress lockdown, mesh-wide access logging, and the node-level CNI agent in place of the injected privileged init-container.

## When to Use

- You are taking a sidecar mesh to production and want the availability and security posture set from day one
- Your platform forbids NET_ADMIN init containers (e.g. OpenShift) — the CNI agent is required there

## Key Configuration Choices

- **HA istiod** (`istiod.autoscale` min 2 / max 5 at 75% CPU, `podDisruptionBudget`, `priorityClassName: system-cluster-critical`) — the control plane survives node drains and is never evicted before the workloads that depend on it
- **Trust domain** (`meshConfig.trustDomain`) — set a stable, organization-unique value BEFORE production; changing it later re-identifies every workload in the mesh
- **Egress lockdown** (`outboundTrafficPolicyMode: REGISTRY_ONLY`) — sidecars only reach destinations in the service registry; every external service must be declared with a KubernetesServiceEntry
- **Access logging** (`accessLogFile: /dev/stdout`) — mesh-wide proxy access logs; per-workload control belongs to KubernetesTelemetry
- **Node-level CNI agent** (`cni.enabled: true`) — tighter pod security everywhere, required on platforms that forbid privileged init containers

## Placeholders to Replace

| Placeholder | Description |
|---|---|
| `prod.mesh.example.internal` | Your organization's trust domain — a stable, unique identity root for the mesh |

## Related Presets

- **01-standard** — chart defaults, the starting point
- **02-ambient** — the sidecar-less data plane
