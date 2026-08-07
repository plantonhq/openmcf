---
title: "Team Namespace Governed"
description: "This preset is the full governance pair and the safe pattern for compute caps: aggregate CPU/memory caps on the namespace (the ResourceQuota) paired with per-container defaults (the companion..."
type: "preset"
rank: "01"
presetSlug: "01-team-namespace-governed"
componentSlug: "resourcequota"
componentTitle: "ResourceQuota"
provider: "kubernetes"
icon: "package"
order: 1
---

# Team Namespace Governed

This preset is the full governance pair and the safe pattern for compute caps: aggregate CPU/memory caps on the namespace (the ResourceQuota) paired with per-container defaults (the companion LimitRange). The pairing matters: once a quota caps `requests.cpu` or `limits.memory`, the API REJECTS pods that omit those requests/limits — the defaults are what keep naive workloads deployable, inheriting sane values instead of being rejected at admission.

## When to Use

- Onboarding a team or environment into a shared cluster with a compute budget
- Any namespace where you want aggregate CPU/memory caps without breaking `kubectl run` and other manifests that don't declare requests/limits
- As the standard per-namespace governance manifest, stamped out with per-team cap values

## Key Configuration Choices

- **All four compute caps (`requests.*` and `limits.*`)** — requests govern what the scheduler reserves; limits govern the runtime ceiling. Capping both bounds reservation AND burst; limits caps are double the requests caps here, allowing 2x aggregate burst
- **`limit_defaults` with a `container` item** — this is what creates the companion LimitRange (sharing the quota's name). Without it, this quota would make the API reject every pod that omits requests/limits
- **`default_request` well below `default_limit`** — naive containers land as Burstable QoS: a small scheduling reservation (100m/128Mi) with room to burst (500m/512Mi). The `default_request` is what each naive container is billed against the quota, so keep it modest
- **Sizing** — the caps here (10 CPU / 20Gi requests) are a starting point; size from observed usage (`kubectl describe resourcequota`) and adjust. Pre-existing usage counts immediately, and an over-quota namespace rejects the next creation

> For simple T-shirt-size governance on namespaces Planton creates, the `KubernetesNamespace` kind's resource profiles manage a quota and limit range internally — this preset is for full-fidelity control or namespaces Planton does not own.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace to govern — the quota and LimitRange land there and cap only that namespace | Your namespace management |

## Related Presets

- **02-object-count-caps** — the safer first step on a live namespace: caps counts, never rejects naive pods
- **03-besteffort-guard** — contain only the pods with no requests/limits at all, leaving declared workloads untouched
