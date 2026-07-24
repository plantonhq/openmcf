# Production cluster preset

A production Qdrant: 3 nodes (the quorum posture — Raft survives one
member loss), generated read-write and read-only API keys living in
the chart-owned Secret (never in a manifest), 8Gi of memory per node
for the hot vector working set, a 50Gi data volume per pod on fast
storage plus an equally-sized separate snapshots volume, required
anti-affinity so a node loss takes one member instead of the quorum,
and a ServiceMonitor for Prometheus.

Know where deployment stops and data begins: this manifest creates a
3-node substrate, but collections replicate per their OWN
`replication_factor`, declared through the Qdrant API when the
collection is created. A collection created with the default factor
of 1 has one copy of each shard even on this cluster — set the factor
(and shard count) per collection to actually use the three nodes.
Backups are also data operations: the snapshots volume gives them a
safe place to land, but taking them is an API call, not spec.

Change first: set the `storage_class` placeholder (a literal class
name, or a valueFrom reference to a KubernetesStorageClass) before
applying. Then size memory to your corpus — the 8Gi here fits a
mid-size embedding set; quantization (a collection-level setting)
stretches it. Keep `snapshots.size` in step with `storage.size`
whenever you grow the data volume — the equal sizing is the point.
Drop `service_monitor_enabled` if the cluster has no Prometheus
Operator CRDs (the install fails on the unknown kind otherwise).

See [02-production-cluster.yaml](./02-production-cluster.yaml) for
the manifest.
