# KubernetesKarpenterNodePool Guide

The judgment this guide carries: a NodePool is where a cluster's fleet
STRATEGY lives — most clusters run several, and Karpenter chooses among
them by weight. This guide covers that composition; the individual knobs
are upstream Karpenter, faithfully mirrored on
[reference.md](reference.md).

## Run several, by purpose

The common shape is not one pool but a few, each with a purpose and a
weight Karpenter picks by: a default on-demand pool, a cheaper spot pool
for interruptible work, a GPU pool for accelerated workloads. Model the
fleet as the set of shapes the cluster's workloads actually need —
requirements (instance types, zones, capacity type), the taints that
keep the wrong pods off specialized pools, and how aggressively
consolidation reclaims under-used nodes.

## Every pool references a NodeClass

`nodeClassRef` points at the AWS machine template
(KubernetesKarpenterEc2NodeClass — AMIs, subnets, security groups, IAM),
required. The pool answers "what shape of node may exist"; the NodeClass
answers "the cloud specifics of booting it." A pool without a resolvable
NodeClass provisions nothing.

## Prerequisite

Both this kind and its NodeClass need the Karpenter controller and its
CRDs on the cluster — see the
[Karpenter guide](../../kuberneteskarpenter/v1alpha1/GUIDE.md) and the
[operator-prerequisite pattern](../../../patterns/operator-prerequisite.md).
NodePools are cluster-scoped (no namespace).

## On the diagram

Each NodePool renders referencing the controller and its NodeClass — the
fleet strategy (how many pool shapes, spot vs on-demand vs GPU) is
readable from the graph, which is exactly the part worth seeing.

## Pairs well with

- KubernetesKarpenter — the controller (required).
- KubernetesKarpenterEc2NodeClass — the machine template every pool
  references.
