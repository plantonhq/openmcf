---
title: "Karpenter Node Pool"
description: "Karpenter Node Pool deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskarpenternodepool"
---

# Kubernetes Karpenter Node Pool

Declares a Karpenter NodePool — the fleet shape Karpenter provisions
machines from. A NodePool answers "what nodes may exist": the
instance-type/zone/capacity-type constraints, the taints and labels new
nodes carry, how long nodes live, how aggressively under-used nodes are
consolidated away, and the resource ceiling the pool may reach. Clusters
routinely run several — a default on-demand pool, a spot pool, a GPU
pool — and Karpenter picks among them by weight. The spec holds 100%
fidelity with the upstream `karpenter.sh/v1` NodePool CRD, with the
CRD's own validation rules mirrored so mistakes surface at validate
time.

## What Gets Created

- **NodePool** (cluster-scoped, named after `metadata.name`) — the
  `karpenter.sh/v1` custom resource carrying the NodeClaim template
  (requirements, taints, labels, lifetime), the disruption policy
  (consolidation and budgets), pool-wide limits, and weight

## Prerequisites

- A Karpenter installation (KubernetesKarpenter) on the cluster — the
  NodePool CRD does not exist before it
- A NodeClass carrying the machine template — on AWS, a
  KubernetesKarpenterEc2NodeClass the pool references through
  `node_class_ref` (one class is typically shared by several pools)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKarpenterNodePool
metadata:
  name: general-spot
spec:
  template:
    nodeClassRef:
      name:
        valueFrom:
          kind: KubernetesKarpenterEc2NodeClass
          name: default-al2023
          fieldPath: status.outputs.node_class_name
    requirements:
      - key: karpenter.sh/capacity-type
        operator: In
        values: [spot]
      - key: kubernetes.io/arch
        operator: In
        values: [amd64]
  limits:
    cpu: "1000"
```

Karpenter then launches nodes from this pool whenever pending pods fit
its constraints, consolidates under-used nodes per the disruption
policy, and recycles nodes at expiry (720h by default).

## Stack Outputs

| Output | Description |
|---|---|
| `node_pool_name` | Name of the cluster-scoped NodePool — the value of the `karpenter.sh/nodepool` label on every node the pool launches |

## Next Steps

Layer pools by purpose: taint dedicated pools (GPU, batch) so only
tolerating pods land there, rank overlapping pools with `weight`, and
bound spot churn with `min_values` instance diversity and disruption
budgets (a `"0"` budget on a cron schedule freezes disruption during
business hours). Cap each pool with `limits` so a runaway workload
cannot provision unbounded capacity. The machine template itself — AMIs,
subnets, security groups, IAM — evolves independently on the
KubernetesKarpenterEc2NodeClass the pool references.
