# KubernetesCilium Guide

The judgment this guide carries: Cilium is the cluster's DATAPLANE, and
the load-bearing choice is primary-CNI versus chaining — pick wrong and
either the cluster was built incompatibly (nodes stuck NotReady with no
Cilium to install) or you ripped out networking that didn't need
replacing.

## Primary CNI vs chaining — decide with the cluster's birth

- **Primary CNI** (default): Cilium owns pod networking — and the
  cluster must be CREATED without another CNI (the per-platform knobs
  are on [reference.md](reference.md)). Nodes stay NotReady until Cilium
  installs; that is by design, not failure — the same
  NotReady-by-design class as OpenBao's seal lifecycle. This choice
  belongs in the cluster proposal, not as an afterthought.
- **Chaining**: Cilium runs ON TOP of the incumbent CNI — it keeps
  IPAM/routing while Cilium adds eBPF policy, load-balancing, and Hubble
  observability. The no-rip-and-replace path for existing clusters.

## What turning on its arms unlocks (and requires)

- `gatewayApi: true` creates the "cilium" GatewayClass — one of the
  catalog's two Gateway API implementations (the
  [Gateway anchor](../../kubernetesgateway/v1/GUIDE.md) carries the
  chain) — and REQUIRES `kubeProxyReplacement` plus the Gateway API CRDs.
- Cilium enforces standard NetworkPolicy (plus its own L7 policies) —
  the [NetworkPolicy guide](../../kubernetesnetworkpolicy/v1/GUIDE.md)'s
  rules assume a CNI that enforces them; Cilium is one.

## Once per cluster, in kube-system

The agent DaemonSet and generated CNI config are cluster singletons
(release name fixed); upstream convention installs into `kube-system` —
cluster infrastructure, not an app namespace.

## On the diagram

Cilium renders as a shared-layer singleton. Its couplings (the cluster's
no-CNI birth setting, kube-proxy replacement, the GatewayClass it
creates) are configuration, not drawn edges — reviewers check the
cluster-creation posture whenever Cilium is primary CNI.

## Pairs well with

- KubernetesGateway / KubernetesGatewayClass — north-south exposure via
  the "cilium" class when `gatewayApi` is on.
- KubernetesNetworkPolicy — the policies this dataplane enforces.
  (Hubble flow observability is part of this spec's own surface, not a
  separate kind.)
