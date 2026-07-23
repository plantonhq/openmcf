---
title: "Replica Set"
description: "This preset declares the production MongoDB posture: a three-member replica set with automated failover (a new primary is elected in seconds when the current one dies), explicit resources, a..."
type: "preset"
rank: "02"
presetSlug: "02-replica-set"
componentSlug: "mongodb"
componentTitle: "MongoDB"
provider: "kubernetes"
icon: "package"
order: 2
---

# Replica Set

This preset declares the production MongoDB posture: a three-member
replica set with automated failover (a new primary is elected in
seconds when the current one dies), explicit resources, a
PodDisruptionBudget, zone-spread anti-affinity, required TLS, and a
declarative application user. Applications connect at the exported
`kube_endpoint` with `?replicaSet=rs0` so the driver follows failovers.

## When to Use

- Any production MongoDB whose loss or downtime matters
- The 30-second choice: this is the standard production shape — add
  the backup block (below) before real data arrives

## Key Configuration Choices

- **`size: 3`** — automated failover needs a majority; three members
  survive one loss. Even numbers waste a vote (declare the `arbiter`
  instead of a fourth data member)
- **`pod_disruption_budget.max_unavailable: 1`** — node drains and
  cluster upgrades can never take down two members at once; the
  majority always survives voluntary disruptions
- **`anti_affinity_topology_key: topology.kubernetes.io/zone`** — the
  upstream default spreads members one-per-node; zone spreading means
  a whole-zone loss costs one member, not the quorum
- **Explicit resources** — WiredTiger sizes its cache from the memory
  limit; an unlimited database pod is at the mercy of every noisy
  neighbor
- **`tls.mode: requireTLS`** — clients must speak TLS (the spec
  default is `preferTLS`); point `tls.issuer` at a cert-manager
  (Cluster)Issuer for an organization-trusted chain
- **A declarative user with no inline password** — the operator
  generates the credential into the `prod-mongo-user-app` Secret and
  keeps the `readWrite` role reconciled; applications never use the
  admin account

**Backups are the one deliberate omission.** Declare a `backup` block
with a named storage (S3/S3-compatible, GCS, or Azure Blob), a nightly
task, and `pitr.enabled: true` before the database holds anything you
cannot lose — the README's production example shows the full shape.
The omission here keeps the preset deployable without a real bucket,
not because backups are optional.

## Placeholders to Replace

None — this preset deploys as-is once KubernetesPerconaMongoOperator
runs in the `percona-mongo` namespace (rename `metadata.name` per
database; every derived object follows it).

## Related Presets

- **01-single-instance** — the development shape: one member, no
  backups, `unsafe.replset_size`
