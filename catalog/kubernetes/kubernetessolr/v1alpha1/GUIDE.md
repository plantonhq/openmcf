# KubernetesSolr Guide

The judgment this guide carries: a SolrCloud is really two dependencies —
its operator and a ZooKeeper ensemble — and its default storage is
ephemeral. Any proposal that omits persistence or assumes ZooKeeper
appears by magic ships data loss.

## ZooKeeper is not optional — decide how it arrives

Every SolrCloud needs a ZooKeeper ensemble. The default (empty
`zookeeper`) provisions a 3-node ensemble via the zookeeper-operator that
the Solr operator's chart bundles; `zookeeper.external` points at one you
already run (the field doc on [reference.md](reference.md)). Either way it
is part of the architecture — a SolrCloud proposed with no account of its
ZooKeeper is incomplete.

## Ephemeral storage is the default — override it

Storage defaults to emptyDir: data is LOST on pod eviction. Declare
`storage.persistent` for anything beyond a throwaway experiment. This is
the single most common Solr mistake because nothing fails at deploy —
the loss happens later, on the first reschedule.

## Operator and credentials

KubernetesSolrOperator is the registry prerequisite; it watches ALL
namespaces by default (fence with `watchNamespaces`), and the mechanics
of why the dependency shows no diagram edge are in the
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md).
Basic-auth bootstraps operator-generated credentials in a Secret —
consume by reference, never inline.

## Namespace ownership

The cluster shares its namespace with consumers — wire `spec.namespace`
to a dedicated KubernetesNamespace, not `createNamespace: true`
([namespace-ownership pattern](../../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesSolrOperator — required (see its
  [guide](../../kubernetessolroperator/v1alpha1/GUIDE.md)).
- KubernetesNamespace — the shared namespace's owner.
