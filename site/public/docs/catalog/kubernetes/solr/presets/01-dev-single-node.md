---
title: "Dev single node preset"
description: "The smallest declarable SolrCloud that actually serves: one Solr node, a single-member provided ZooKeeper, ephemeral storage, and no authentication. For developers and CI who need the real SolrCloud..."
type: "preset"
rank: "01"
presetSlug: "01-dev-single-node"
componentSlug: "solr"
componentTitle: "Solr"
provider: "kubernetes"
icon: "package"
order: 1
---

# Dev single node preset

The smallest declarable SolrCloud that actually serves: one Solr
node, a single-member provided ZooKeeper, ephemeral storage, and no
authentication. For developers and CI who need the real SolrCloud
API — collections, schemas, distributed queries — without production
ceremony.

The trade-offs are total: emptyDir storage means every pod eviction
erases the indices, a one-member ZooKeeper is not a quorum so any
restart of it pauses the cluster, and the endpoints are open to
anything in the cluster network. Nothing here transfers to production
except the manifest's shape — do not grow this in place; graduate to
the production preset's three-and-three topology, persistent storage
and basic auth before anything depends on the data.

The first thing to change is `version` (pin the line your clients
test against). If your CI indexes more than toy data, raise the
memory/`java_mem` pair together, keeping heap at about half the
container memory.

See [01-dev-single-node.yaml](./01-dev-single-node.yaml) for the
manifest.
