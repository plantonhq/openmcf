# KubernetesNetworkPolicy Guide

The judgment this guide carries: NetworkPolicy semantics invert the
moment a pod is selected by ANY policy — everything not explicitly
allowed is denied in that direction — so the composition question is
never one policy, it is the namespace's whole allow-list. And none of it
does anything without a CNI that enforces it.

## Default-deny is a namespace posture, not a one-off

The platform's namespace kind can open with default-deny plus explicit
allows (the [KubernetesNamespace guide](../kubernetesnamespace/GUIDE.md)
covers that surface) — reach for THIS standalone kind when policies must
target pods a Planton workload does not manage, span selectors across
workloads, or express rules beyond the namespace kind's arms. Remember
the additive-allow model on [reference.md](v1alpha1/reference.md): selecting a pod
with one narrow policy silently denies everything else to it; the
proposal must include the full set of allows the namespace's traffic
actually needs (DNS egress being the classic forgotten one).

## Policies need an enforcer

NetworkPolicy objects are inert on a CNI that does not implement them.
[KubernetesCilium](../kubernetescilium/GUIDE.md) is the catalog's
policy-enforcing dataplane; on managed clusters, verify the CNI in play
enforces policy before proposing isolation that will silently not exist.

## On the diagram

Policies select by label, not by reference — no edges are drawn. The
reviewable artifact is the namespace's policy SET; state it as a unit in
proposals ("default-deny + these allows") rather than scattering
individual rules.

## Pairs well with

- KubernetesNamespace — the namespace-level default-deny surface.
- KubernetesCilium — the enforcing dataplane.
- The workloads whose `selector_labels` outputs feed `podSelector`
  blocks (the Deployment page exports them for exactly this).
