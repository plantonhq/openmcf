# KubernetesGhaRunnerScaleSet Guide

The judgment this guide carries: workflows reach this fleet by NAME —
`runs-on: <the scale set name>` — not by labels; the GitHub credential
always rides in a Secret; and `docker build` jobs fail on the default
runner unless a container mode is chosen up front.

## Routing is by name, and the name is the contract

The fleet registers in GitHub under its name (default:
`metadata.name`, max 45 chars) and workflows select it with
`runs-on: <that name>` — runner labels are NOT how scale-set routing
works (the reference page states it). Renaming the fleet breaks every
workflow that targets it; treat the name like an API.

## Credentials and Docker builds — the two up-front choices

- **Auth is Secret-native**: reference an existing Secret
  (`auth.existingSecretName`, the recommended posture — pair with
  KubernetesExternalSecret when the PAT/GitHub App lives in a secrets
  manager) or declare inline and let the module materialize the Secret;
  inline values are sensitive-marked and never reach chart values.
- **Container mode**: the default runner runs plain jobs only — `docker
  build` needs `containerMode.mode: dind` (the reference page has the
  arms). A CI fleet proposed for image builds without a container mode
  fails on its first real job.

## The controller must exist first

KubernetesGhaRunnerScaleSetController is the prerequisite — one per
cluster, watching everything by default (its
[guide](../kubernetesgharunnerscalesetcontroller/GUIDE.md) has the
fencing judgment;
[operator-prerequisite pattern](../../_patterns/operator-prerequisite.md)
for the mechanism). One fleet per repo/org/enterprise registration.

## Namespace ownership

Fleets commonly share a CI namespace with their secrets — wire
`spec.namespace` through a dedicated KubernetesNamespace
([namespace-ownership pattern](../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesGhaRunnerScaleSetController — required, once per cluster.
- KubernetesExternalSecret — the GitHub credential's managed path.
