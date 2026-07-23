---
title: "Kind / Local Dev Cluster"
description: "This preset installs Cilium as the primary CNI on a local development cluster (kind, or any single-node cluster created without a default CNI). It uses `kubernetes` IPAM — the control plane already..."
type: "preset"
rank: "01"
presetSlug: "01-kind-dev-cluster"
componentSlug: "cilium"
componentTitle: "Cilium"
provider: "kubernetes"
icon: "package"
order: 1
---

# Kind / Local Dev Cluster

This preset installs Cilium as the primary CNI on a local development
cluster (kind, or any single-node cluster created without a default CNI).
It uses `kubernetes` IPAM — the control plane already allocates per-node
PodCIDRs, so Cilium simply consumes them — and runs a single operator
replica, because the chart's 2-replica HA default can never fully schedule
on a single node (the replicas carry pod anti-affinity) and the rollout
would never settle.

## When to Use

- kind clusters created with `disableDefaultCNI: true` (nodes sit NotReady
  until Cilium installs — that is by design)
- Single-node or laptop-scale clusters for developing and testing
  NetworkPolicy, Cilium policies, or Hubble locally
- Any cluster whose control plane already assigns node PodCIDRs (kind,
  kubeadm defaults)

## Key Configuration Choices

- **`namespace: kube-system`** (`createNamespace: false`) — the upstream
  convention; several chart defaults assume it, and kube-system always
  exists
- **`ipam.mode: kubernetes`** — pods get IPs from the node's
  Kubernetes-assigned PodCIDR; no cluster-pool CIDR planning needed
- **`operator.replicas: 1`** — mandatory on single-node clusters; the
  second replica of the chart default cannot schedule
- **Everything else at chart defaults** — tunnel routing (vxlan), Hubble
  enabled in the agent (no relay/UI), policy enforcement `default`

## Placeholders to Replace

No placeholders — this preset is directly deployable.

## Related Presets

- **02-eks-chaining** — run Cilium's policy/observability on top of the
  AWS VPC CNI on EKS
- **03-self-managed-primary-kpr** — the kube-proxy-free primary-CNI
  posture for self-managed multi-node clusters
- **04-production-observability** — production hardening: Hubble
  relay/UI/metrics, Prometheus, WireGuard encryption
