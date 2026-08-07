# KubernetesSolrOperator Guide

The judgment this guide carries: one operator serves the cluster and it
watches everything by default — and it bundles the zookeeper-operator that
KubernetesSolr clusters lean on, so it is more than a single-CRD
controller.

## Watches all namespaces by default

Empty `watchNamespaces` means the operator reconciles SolrCloud resources
in every namespace; set an explicit list only to fence it. One install
per cluster in the shared-cluster chart is the normal shape. The
invisible-edge mechanism and the three watch postures:
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md).

## It brings ZooKeeper's operator with it

The chart bundles the zookeeper-operator, which is what lets a
KubernetesSolr with the default `zookeeper` block get a provisioned
3-node ensemble (the Solr guide covers that choice). Worth knowing when
reasoning about what this one install actually places on the cluster.

## Namespace ownership — the infra exception

A dedicated single-tenant namespace with `createNamespace: true` is the
normal shape — the
[namespace-ownership pattern](../../../_patterns/namespace-ownership.md)'s
sole-tenant case.

## Pairs well with

- KubernetesSolr — the SolrCloud clusters this operator reconciles (see
  its [guide](../../kubernetessolr/v1alpha1/GUIDE.md)).
