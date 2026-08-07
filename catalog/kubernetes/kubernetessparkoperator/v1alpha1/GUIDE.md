# KubernetesSparkOperator Guide

The judgment this guide carries: this is an engine with no catalog CR to
pair — Spark workloads (SparkApplication, SparkCluster) are authored
directly and carried by KubernetesManifest today, the same
engine-plus-declarations shape as KEDA and the policy engines.

## Engine here, workloads authored beside it

Installing this operator runs no Spark. Declare the actual work as
SparkApplication (one batch/streaming job, run to completion) or
SparkCluster (a long-lived standalone cluster) custom resources — via
KubernetesManifest, in a namespace the operator watches. A proposal for
"Spark on the cluster" names both halves: this engine AND the
application/cluster CRs, or nothing runs. Diagram consequence: the Spark
workloads themselves ride inside opaque manifest nodes — the actual jobs
are invisible in the rendered architecture, so name them in the
proposal's prose.

## Watch scope and CRD lifecycle

Cluster-wide watch by default (fence with `watchNamespaces`); the CRDs
are keep-on-uninstall — destroying the operator strands, not deletes,
declared Spark objects (the reference page carries both). The
invisible-edge mechanism:
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md).

## Namespace ownership — the infra exception

Dedicated single-tenant namespace, `createNamespace: true` — the
[namespace-ownership pattern](../../../_patterns/namespace-ownership.md)'s
sole-tenant case.

## Pairs well with

- KubernetesManifest — carries the SparkApplication/SparkCluster CRs
  today.
- KubernetesSeaweedFs — the S3-compatible store Spark jobs commonly
  read/write.
