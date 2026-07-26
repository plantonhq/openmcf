---
title: "Dev single instance preset"
description: "The smallest declarable Neo4j: community edition, a declared admin password (materialized as a Kubernetes Secret, never in rendered values), a 10Gi data volume on the cluster's default StorageClass,..."
type: "preset"
rank: "01"
presetSlug: "01-dev-single-instance"
componentSlug: "neo4j"
componentTitle: "Neo4j"
provider: "kubernetes"
icon: "package"
order: 1
---

# Dev single instance preset

The smallest declarable Neo4j: community edition, a declared admin
password (materialized as a Kubernetes Secret, never in rendered
values), a 10Gi data volume on the cluster's default StorageClass,
and resources at the chart's own floor — 500m CPU / 2Gi memory, below
which the chart refuses to install at all. For developers who need a
real Cypher endpoint, and for the first iteration of a knowledge
graph or agent-memory experiment.

Do not shrink this preset to save resources — there is no smaller
Neo4j; the floor is the chart's, not this preset's. And do not leave
the placeholder password in place beyond a private dev cluster: the
data volume defaults to the cluster's default StorageClass with the
graph's entire state on it, so a guessable credential is the whole
ballgame. In-cluster clients connect via the default Service in the
stack outputs (bolt 7687); nothing is exposed outside the cluster.

Change first: the `auth.password` (or switch to `existing_secret`
carrying the chart's NEO4J_AUTH contract), then `data_volume.size` if
your graph will outgrow 10Gi — growing a PVC later depends on the
StorageClass allowing expansion.

See [01-dev-single-instance.yaml](./01-dev-single-instance.yaml) for
the manifest.
