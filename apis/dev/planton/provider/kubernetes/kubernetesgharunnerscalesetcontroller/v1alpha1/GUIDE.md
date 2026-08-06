# KubernetesGhaRunnerScaleSetController Guide

The judgment this guide carries: one controller per cluster serves every
runner fleet — install it once, then declare fleets per repository or
organization. Fencing it is the exception with a real cost.

## Once per cluster, watching everything

The controller watches all namespaces by default, so every
KubernetesGhaRunnerScaleSet on the cluster is served by one install —
the sane default (its own guide carries the fleet-side judgment).
`flags.watchSingleNamespace` fences it for hard multi-tenancy only, and
then every fenced fleet must reference its controller explicitly — a
heavier shape to justify, not a tidy-by-default. The invisible-edge
mechanism: [operator-prerequisite pattern](../../../patterns/operator-prerequisite.md).

## Destroying the controller keeps the CRDs (and the fleets)

The actions.github.com CRDs are keep-on-uninstall — destroying this
controller strands, not deletes, the declared fleets' objects (verified
live per the reference page). Same family as the Percona operators;
opposite of RabbitMQ's. Sequence teardown accordingly.

## Namespace ownership — the infra exception

Dedicated single-tenant namespace, `createNamespace: true` — the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case.

## Pairs well with

- KubernetesGhaRunnerScaleSet — the fleets this controller reconciles
  (see its [guide](../../kubernetesgharunnerscaleset/v1alpha1/GUIDE.md)).
