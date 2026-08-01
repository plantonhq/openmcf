---
title: "KubeRay Operator"
description: "KubeRay Operator deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskuberayoperator"
---

# KubeRay Operator

The controller that turns `RayCluster` declarations into running Ray
clusters, from the official `kuberay-operator` Helm chart (1.6.2 =
operator v1.6.2). One operator per cluster: it watches every namespace
by default, runs leader election out of the box, and reconciles the
clusters declared with `KubernetesRayCluster` (plus RayJob/RayService
CRs authored directly).

## Highlights

- **Removing the operator never deletes Ray declarations** — the three
  ray.io CRDs ride the chart's `crds/` directory: installed once,
  never upgraded by chart bumps, kept on uninstall by upstream design.
- **No webhook, no cert-manager** — the operator validates in its
  reconcile loop; a bad RayCluster surfaces on the CR's status
  conditions, not as an admission rejection.
- **Feature gates that don't bite** — the chart's `featureGates` is a
  Helm list, and lists replace rather than merge; the modules merge
  your flips over the chart's five defaults so unlisted gates never
  silently change state.
- **Two image seams, honestly separated** — `image_registry` mirrors
  the operator's own `quay.io/kuberay/operator` image; Ray cluster
  images ride each `KubernetesRayCluster`'s own field, and mirroring
  one does nothing for the other.
- **Fail-loud, not fail-later** — names over the 47-character budget
  are rejected at apply time, and the atomic install waits for the
  operator to become Available instead of surfacing as RayClusters
  that never reconcile.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
