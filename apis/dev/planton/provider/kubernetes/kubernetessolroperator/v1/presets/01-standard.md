# Standard preset

The default install: one cluster-wide Solr operator plus the bundled
zookeeper-operator, on the stable 0.9.1 chart. This is the shape for
a cluster you administer — install once, and every KubernetesSolr
resource in any namespace gets reconciled, including its provided
ZooKeeper ensemble, with nothing else to set up.

Do not use this when a zookeeper-operator already runs in the
cluster: its cluster-scoped RBAC uses fixed names, so a second
install conflicts — use the existing-zookeeper-operator preset
instead. And on shared clusters where you only own certain
namespaces, prefer the namespace-fenced preset so this install does
not claim SolrCloud resources belonging to other teams.

The first thing to change is nothing — the defaults are the
recommended posture. If your platform requires it, add
`node_selector`/`tolerations` for pod placement or `resources` for
namespace quotas (the chart ships none; the operator is lightweight).

See [01-standard.yaml](./01-standard.yaml) for the manifest.
