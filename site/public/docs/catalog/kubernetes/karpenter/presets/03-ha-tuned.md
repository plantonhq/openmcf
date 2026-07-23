---
title: "HA Tuned"
description: "This preset hardens the EKS-standard installation for clusters where Karpenter is load-bearing: explicit two-replica sizing with the documented large-cluster resource starting point, batching tuned..."
type: "preset"
rank: "03"
presetSlug: "03-ha-tuned"
componentSlug: "karpenter"
componentTitle: "Karpenter"
provider: "kubernetes"
icon: "package"
order: 3
---

# HA Tuned

This preset hardens the EKS-standard installation for clusters where
Karpenter is load-bearing: explicit two-replica sizing with the
documented large-cluster resource starting point, batching tuned to
produce fewer and larger nodes, spot-to-spot consolidation for spot-heavy
fleets, and Prometheus scrape discovery for the controller's own metrics.

## When to Use

- Large or spot-heavy clusters where provisioning decisions and
  consolidation behavior are worth tuning and observing
- Production clusters that want the controller's failover and telemetry
  posture stated explicitly rather than inherited from chart defaults

## Key Configuration Choices

- **`controller.replicas: 2`** — the chart default made explicit: an
  active leader plus a warm standby (leader-elected failover capacity,
  not horizontal scale), zone-spread by the chart's topology constraints
- **`controller.resources: 1 CPU / 1Gi`** — the commonly documented
  starting point for large clusters (upstream otherwise leaves
  requests/limits to the operator)
- **`batching.maxDuration: 30s`** (chart default 10s) — pending pods wait
  longer per provisioning decision so more are considered at once, which
  usually produces fewer, larger nodes; `idleDuration: 2s` still closes
  quiet batches early
- **`featureGates.spotToSpotConsolidation: true`** — ALPHA gate allowing
  consolidation to replace spot nodes with cheaper spot nodes; off by
  default upstream because it can increase churn
- **`prometheus.serviceMonitor: true`** — requires the Prometheus
  operator CRDs on the cluster; the release fails to install without them

## Placeholders to Replace

| Placeholder                                           | Description                                 | Where to Find                          |
| ----------------------------------------------------- | ------------------------------------------- | -------------------------------------- |
| `<eks-cluster-name>`                                  | EKS cluster name                            | EKS console or `AwsEksCluster` outputs |
| `arn:aws:iam::123456789012:role/karpenter-controller` | IRSA role ARN — replace account id and name | IAM console                            |
| `karpenter-interruptions`                             | SQS interruption queue name                 | SQS console                            |

## Related Presets

- **01-eks-standard** — the baseline installation this preset hardens
- **02-eks-isolated-vpc** — isolated-VPC and custom-networking clusters
