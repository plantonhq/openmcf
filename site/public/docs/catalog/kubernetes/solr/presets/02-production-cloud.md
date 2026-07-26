---
title: "Production cloud preset"
description: "The production topology: three Solr nodes on persistent fast storage, a three-member ZooKeeper quorum with its own persistent volumes, basic-auth security bootstrapped by the operator, shard-aware..."
type: "preset"
rank: "02"
presetSlug: "02-production-cloud"
componentSlug: "solr"
componentTitle: "Solr"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production cloud preset

The production topology: three Solr nodes on persistent fast storage,
a three-member ZooKeeper quorum with its own persistent volumes,
basic-auth security bootstrapped by the operator, shard-aware managed
rolling updates bounded to one pod and one shard replica at a time,
and the cluster-wide PodDisruptionBudget. Retain reclaim policy means
the index volumes outlive the resource — deletion is not data loss.

Do not use this for development (the dev preset costs a tenth of the
resources) and do not treat it as internet-ready: basic auth secures
the API, but the Service is in-cluster only — compose exposure
deliberately (a KubernetesIngress over the common service, or the
`external` block when clients need per-node addressability for
CloudSolrClient). Also do not scale ZooKeeper like Solr: the ensemble
stays at three (or five) members regardless of how many Solr nodes
you add.

Change first: the two `storage_class` placeholders (a literal class
name, or a valueFrom reference to a KubernetesStorageClass), then
`storage.persistent.size` and the memory/`java_mem` pair to match
your index volume — keeping heap at about half the container memory.

See [02-production-cloud.yaml](./02-production-cloud.yaml) for the
manifest.
