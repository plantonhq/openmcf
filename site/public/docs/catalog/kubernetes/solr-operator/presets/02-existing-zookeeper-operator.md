---
title: "Existing zookeeper operator preset"
description: "The install for clusters that already run a zookeeper-operator: `install: false` skips the bundled dependency (whose fixed-name, cluster-scoped RBAC would conflict with the existing one), and..."
type: "preset"
rank: "02"
presetSlug: "02-existing-zookeeper-operator"
componentSlug: "solr-operator"
componentTitle: "Solr Operator"
provider: "kubernetes"
icon: "package"
order: 2
---

# Existing zookeeper operator preset

The install for clusters that already run a zookeeper-operator:
`install: false` skips the bundled dependency (whose fixed-name,
cluster-scoped RBAC would conflict with the existing one), and
`use_existing: true` tells the Solr operator an ensemble provisioner
is nonetheless present. KubernetesSolr resources keep working exactly
as with the standard preset — provided ZooKeeper ensembles are
reconciled by the operator that was already there.

Do not use this on a cluster with no zookeeper-operator: provided
ensembles would sit unprovisioned forever (SolrCloud waits on
ZooKeeper), and the failure is a silent hang, not an error. Also skip
the `use_existing` flag if every Solr cluster will connect to an
EXTERNAL ensemble (`zookeeper.external` on KubernetesSolr) — then
plain `install: false` alone is the honest configuration.

Change first: verify which release owns the existing
zookeeper-operator and which version it runs — the Solr operator
assumes a compatible ZookeeperCluster CRD is being served.

See
[02-existing-zookeeper-operator.yaml](./02-existing-zookeeper-operator.yaml)
for the manifest.
