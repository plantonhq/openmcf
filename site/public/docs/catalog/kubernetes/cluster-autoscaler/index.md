---
title: "Cluster Autoscaler"
description: "Cluster Autoscaler deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesclusterautoscaler"
---

# Kubernetes Cluster Autoscaler

Installs the Kubernetes Cluster Autoscaler from the official Helm chart,
with a typed spec over the chart's meaningful configuration surface. The
autoscaler grows and shrinks EXISTING node groups: when pods are
unschedulable it raises the desired size of a matching group (an EC2
Auto Scaling group, an Azure VMSS, a Cluster API MachineDeployment,
...), and it scales groups back down when nodes sit underutilized.
Exactly one provider arm per installation — AWS, Azure, GCE, Cluster
API, Civo, or the KWOK simulation sandbox. One installation per cluster.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set; `kube-system` (pre-existing) is the upstream convention
- **Helm Release** (`cluster-autoscaler`) — the autoscaler Deployment
  (chart default 1 replica; extras leader-elect as warm standbys), RBAC,
  the provider-specific credential Secret when declared credentials are
  used, and the chart-derived service account cloud-side keyless
  bindings are written against

## Prerequisites

- Node capacity organized as pre-defined groups with size bounds — ASGs
  tagged for auto-discovery on AWS (`k8s.io/cluster-autoscaler/enabled`
  + `k8s.io/cluster-autoscaler/<cluster_name>`), VMSS on Azure, MIG name
  prefixes on GCE, annotated MachineDeployments on Cluster API
- Cloud-side identity for the keyless postures: an IRSA role, GCP
  Workload Identity binding, or Azure workload/managed identity written
  against the chart's derived service account
- On GKE/AKS: a deliberate decision to run self-managed autoscaling —
  both platforms ship a managed autoscaler as a node-pool toggle
- With `prometheus.service_monitor`: the Prometheus operator CRDs — the
  release fails to install without them

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesClusterAutoscaler
metadata:
  name: cluster-autoscaler
spec:
  namespace:
    value: kube-system
  aws:
    region: us-west-2
    autoDiscovery:
      clusterName: my-eks-cluster
    irsaRoleArn: arn:aws:iam::111111111111:role/cluster-autoscaler
```

The autoscaler then manages every ASG carrying the discovery tags for
the cluster — new node groups enroll by tagging alone.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (always `cluster-autoscaler`) |
| `service_account_name` | Chart-derived service account (embeds the provider arm, e.g. `cluster-autoscaler-aws-cluster-autoscaler`) — the subject cloud-side keyless bindings are written against |

## Next Steps

Tune the `scaling` block as the cluster grows: `least-waste` is the
common production expander, `balance_similar_node_groups` keeps
one-group-per-zone layouts even, and the `scale_down` thresholds control
how aggressively idle capacity is reclaimed. Reach for `extra_args` for
the autoscaler's long tail of flags — it is the chart's own contract.
Keep the autoscaler's minor version aligned with the cluster's
Kubernetes minor per upstream guidance. For AWS clusters that would
rather launch right-sized machines on demand than manage groups,
KubernetesKarpenter is the modern alternative.
