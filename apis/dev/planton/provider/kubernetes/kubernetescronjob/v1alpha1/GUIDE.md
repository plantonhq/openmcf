# KubernetesCronJob Guide

The judgment this guide carries: scheduled work usually operates on some
other component's state — a database it backs up, a queue it drains — and
the composition mistakes happen in that relationship, not in the schedule
expression.

## Compose it beside what it operates on

A cron job that touches another component's data belongs in that
component's namespace, wired through the same KubernetesNamespace
reference — which also puts the relationship on the diagram. Credentials
come the platform way: reference the target's exported secret handles
(e.g. a database's credential Secret documented on its reference page)
rather than copying values into the job's environment.

## Overlap is opt-in — keep it that way deliberately

The platform defaults `concurrencyPolicy` to Forbid (stricter than the
Kubernetes default — the reasoning is on [reference.md](reference.md)).
Before opting out, name which property you need: `Replace` suits jobs
where only the latest run matters (a sync whose stale run is worthless);
`Allow` is only safe when runs share nothing — two backups writing one
target is the classic incident the default exists to prevent.

## Namespace ownership

`spec.namespace` is a required foreign key targeting KubernetesNamespace.
`createNamespace: true` makes this cron job the namespace's owner in IaC
state — almost never right for a job that exists to operate on neighbors
in the same namespace. The failure story and the `valueFrom` wiring:
[namespace-ownership pattern](../../../patterns/namespace-ownership.md).

## Pipeline-attached scheduled work

When the job runs your own code as a Service Hub deploy target, the
pipeline injects the image at `spec.jobTemplate.container.app.image` (the
deploy-target contract on [reference.md](reference.md)) — the scheduled
job stays on the same release train as the service it belongs to.

## On the diagram

A cron job wired by `valueFrom` to its namespace and its target's outputs
renders with the edges that explain WHY it exists. One created with
`createNamespace: true` and copied credential strings renders as an
orphan node beside the system it actually serves.

## Pairs well with

- KubernetesNamespace — the namespace owner (pattern above).
- The component whose state the job operates on — wire its exported
  outputs via `valueFrom` instead of copying values.
- KubernetesJob — the same batch controls for one-shot work (its
  reference page carries the comparison).
