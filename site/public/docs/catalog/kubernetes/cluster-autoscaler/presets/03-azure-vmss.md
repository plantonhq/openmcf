---
title: "Azure VMSS"
description: "This preset installs the Cluster Autoscaler against Azure VM scale sets with tag-based auto-discovery and federated workload identity — the keyless credential posture. It is the arm for AKS clusters..."
type: "preset"
rank: "03"
presetSlug: "03-azure-vmss"
componentSlug: "cluster-autoscaler"
componentTitle: "Cluster Autoscaler"
provider: "kubernetes"
icon: "package"
order: 3
---

# Azure VMSS

This preset installs the Cluster Autoscaler against Azure VM scale sets
with tag-based auto-discovery and federated workload identity — the
keyless credential posture. It is the arm for AKS clusters that opt out
of the managed autoscaler and for self-managed clusters running on VMSS.

Note: AKS ships a MANAGED autoscaler configured as a toggle on the node
pool itself — deploying this component on AKS is the exception, not the
rule. Use it when you need autoscaler versions, flags, or behavior the
managed toggle does not expose.

## When to Use

- Self-managed Kubernetes clusters on Azure VM scale sets
- AKS clusters that deliberately opt out of the managed autoscaler for
  flag-level control

## Key Configuration Choices

- **Tag-based VMSS auto-discovery (`clusterName`)** — scale sets tagged
  per upstream's Azure auto-discovery setup are managed automatically;
  the alternative is a static `nodeGroups` list with explicit size
  bounds (exactly one of the two)
- **`identity.useWorkloadIdentity: true`** — federated workload identity
  (recommended): the autoscaler's service account exchanges its token
  for the federated Entra application, configured Azure-side; no client
  secret is stored in the manifest. Managed identity and
  service-principal credentials are the alternative arms (exactly one)
- **`resourceGroup` is the NODE resource group** — for AKS, the `MC_*`
  group holding the VMSS instances, not the group the cluster resource
  lives in
- **`expander: least-waste`** — the common production choice (upstream
  default is `random`)

## Placeholders to Replace

| Placeholder                 | Description                                      | Where to Find                                  |
| --------------------------- | ------------------------------------------------ | ---------------------------------------------- |
| `<azure-subscription-id>`   | Subscription holding the scale sets              | Azure portal                                   |
| `<node-resource-group>`     | Node resource group (MC_* for AKS)               | Azure portal — the cluster's node resource group |
| `<cluster-name>`            | Cluster name used by the VMSS discovery tags     | Your VMSS tags per upstream's Azure setup      |

## Related Presets

- **01-eks-autodiscovery** — EKS Auto Scaling groups with IRSA
- **02-cluster-api** — Cluster API MachineDeployments
