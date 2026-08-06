# KubernetesOtelOperator Guide

The judgment this guide carries: one operator per cluster, watching
everything — the standard shape. The decisions all live on the collector
side, so this guide is deliberately short.

## Once per cluster, watches cluster-wide

The normal posture is a single install in the shared-cluster chart
reconciling OpenTelemetryCollector CRs in every namespace. The
invisible-edge mechanism and watch-posture table:
[operator-prerequisite pattern](../../../patterns/operator-prerequisite.md).
Its own dedicated namespace with `createNamespace: true` is the
sole-tenant case of the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md).

## The judgment lives on the collector

Mode, pipeline config, credential handling, and cluster-read RBAC are all
KubernetesOtelCollector decisions — see the
[collector guide](../../kubernetesotelcollector/v1alpha1/GUIDE.md). Installing
this operator alone runs no collectors.

## Pairs well with

- KubernetesOtelCollector — the collectors this operator reconciles (see
  its [guide](../../kubernetesotelcollector/v1alpha1/GUIDE.md)).
