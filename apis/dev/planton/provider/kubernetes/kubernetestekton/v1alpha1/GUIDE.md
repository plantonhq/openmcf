# KubernetesTekton Guide

The judgment this guide carries: Tekton is configured cluster-wide through
exactly one declaration, its namespace can never change in place, and its
teardown depends on the operator still being alive — three constraints
that make the ORDER of operations matter more than any field.

## Exactly one per cluster

Unlike the database CR kinds (one per cluster instance), Tekton's operator
admits exactly ONE TektonConfig, named `config` — the module renders that
fixed name regardless of `metadata.name` (the reference page states this).
Declare a single KubernetesTekton per cluster; it configures which
components run (Pipelines, Triggers, Dashboard, Chains), not one of many
instances. This is the singleton posture of the
[operator-prerequisite pattern](../../../patterns/operator-prerequisite.md),
not the per-namespace-instance one.

## The namespace is immutable; destroy order is load-bearing

`targetNamespace` cannot be changed on an existing installation — the
operator's webhook rejects it; moving Tekton means destroy-and-recreate.
And destroying this resource makes the operator tear down every installed
component through finalizers only it can process — so NEVER destroy
KubernetesTektonOperator before this resource, or the finalizers strand
and deletion hangs. Any teardown or namespace-move proposal must sequence
these deliberately.

## Operator prerequisite

KubernetesTektonOperator is installed with automatic component
installation disabled precisely so this declaration is the single owner of
the cluster's TektonConfig — installing the operator alone runs no Tekton.
Its guide carries the manager-side lifecycle.

## Pairs well with

- KubernetesTektonOperator — required; the manager whose lifecycle this
  config drives (see its
  [guide](../../kubernetestektonoperator/v1alpha1/GUIDE.md)).
- Gateway API kinds — external exposure for the dashboard when the `all`
  profile installs it.
