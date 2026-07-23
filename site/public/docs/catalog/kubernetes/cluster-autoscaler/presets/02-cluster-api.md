---
title: "Cluster API"
description: "This preset installs the Cluster Autoscaler with the Cluster API provider in the self-managed mode: both the workload cluster and the CAPI management objects live in the same cluster, machine..."
type: "preset"
rank: "02"
presetSlug: "02-cluster-api"
componentSlug: "cluster-autoscaler"
componentTitle: "Cluster Autoscaler"
provider: "kubernetes"
icon: "package"
order: 2
---

# Cluster API

This preset installs the Cluster Autoscaler with the Cluster API provider
in the self-managed mode: both the workload cluster and the CAPI
management objects live in the same cluster, machine discovery is fenced
to one namespace, and RBAC is namespace-scoped — the least-privilege
posture for a single-cluster CAPI setup.

## When to Use

- Self-managed or multi-cloud clusters whose node capacity is declared as
  Cluster API MachineDeployments/MachineSets (annotated for autoscaling
  per upstream's Cluster API guidance)
- Single-cluster CAPI topologies where the management objects live in the
  workload cluster itself

## Key Configuration Choices

- **`mode: incluster-incluster`** (chart default, made explicit) — the
  autoscaler finds both the workload cluster and the CAPI objects in
  this cluster; the four kubeconfig-based modes exist for split
  management/workload topologies and additionally require a
  `kubeconfigSecret` mounted per upstream guidance
- **Namespace fence (`namespace`)** — only machine objects in the fenced
  namespace are considered for autoscaling; empty would watch all
  namespaces
- **`namespaceScopedRbac: true`** — restricts RBAC to namespace scope
  (`rbac.clusterScoped=false`); the least-privilege posture when the
  autoscaler only manages machines in one namespace
- **No cloud credentials** — scaling happens by editing CAPI objects
  through the Kubernetes API, so there is no IRSA/identity block

## Placeholders to Replace

| Placeholder                | Description                                    | Where to Find                        |
| -------------------------- | ---------------------------------------------- | ------------------------------------ |
| `<capi-cluster-namespace>` | Namespace holding the MachineDeployments to autoscale | Your Cluster API management manifests |

## Related Presets

- **01-eks-autodiscovery** — EKS Auto Scaling groups with IRSA
- **03-azure-vmss** — Azure VM scale sets with workload identity
