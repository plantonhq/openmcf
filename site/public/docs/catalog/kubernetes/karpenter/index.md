---
title: "Karpenter"
description: "Karpenter deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskarpenter"
---

# Kubernetes Karpenter

Installs the Karpenter node-provisioning controller from the official
OCI-served Helm charts, with a typed spec over the chart's meaningful
configuration surface. Karpenter watches for unschedulable pods and
launches RIGHT-SIZED machines for them in seconds — no pre-created node
groups, no per-ASG scaling policies: it picks instance types from the
live catalog to fit the pending pods, consolidates under-used nodes
away, and handles spot interruptions. This kind installs the AWS
provider — the canonical implementation. One installation per cluster.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set; `kube-system` (pre-existing) is upstream's recommended home
- **CRD Helm Release** (`karpenter-crd`) — the NodePool, NodeClaim, and
  EC2NodeClass definitions as their own release, upstream's supported
  path for keeping CRDs upgradable — annotated to survive uninstall by
  default, so removing the release does not cascade-delete every fleet
  declaration in the cluster
- **Controller Helm Release** (`karpenter`) — the controller Deployment
  (chart default 2 replicas: leader plus warm standby), RBAC, and the
  `karpenter` service account; the chart pins controller pods away from
  Karpenter-provisioned nodes so the controller never disrupts its own
  machine

## Prerequisites

- An EKS (or EKS-compatible) cluster with capacity Karpenter does NOT
  manage for the controller itself — a managed node group or Fargate
- AWS-side identity for the controller: an IRSA role (trust policy
  against the `karpenter` service account) or an EKS Pod Identity
  association, with upstream's controller IAM policy
- For interruption handling: the SQS queue and EventBridge rules from
  upstream's guidance, named in `aws.interruption_queue`
- With `prometheus.service_monitor`: the Prometheus operator CRDs — the
  release fails to install without them

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKarpenter
metadata:
  name: karpenter
spec:
  namespace:
    value: kube-system
  cluster:
    name: my-eks-cluster
    eksControlPlane: true
  aws:
    irsaRoleArn: arn:aws:iam::111111111111:role/karpenter-controller
```

The engine alone provisions nothing: declare the fleet with
KubernetesKarpenterEc2NodeClass (the machine template) and
KubernetesKarpenterNodePool (the constraints), and Karpenter starts
launching nodes for pending pods.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Controller Helm release name (always `karpenter`) |
| `crd_release_name` | CRD Helm release name (`karpenter-crd`; empty when CRDs are managed elsewhere) |
| `service_account_name` | Always `karpenter` — the subject IRSA trust policies and EKS Pod Identity associations are written against |

## Next Steps

Create an EC2NodeClass and at least one NodePool — a general-purpose
on-demand pool is the usual first fleet; spot and GPU pools layer on as
siblings weighted by preference. Wire `aws.interruption_queue` before
running spot at scale so nodes drain ahead of interruptions. When
provisioning becomes load-bearing, add the ServiceMonitor for
provisioning-latency and disruption-decision telemetry. For clusters
whose capacity is organized as pre-defined node groups instead,
KubernetesClusterAutoscaler is the matching tool.
