# Production preset

A production single server: credentials referenced from a
pre-existing Secret (never declared in the manifest), an explicit
heap/page-cache memory split sized to the container, a 100Gi data
volume on fast storage, transaction-timeout and telemetry settings in
neo4j.conf, node-selector scheduling, and the ClusterIP service
posture with a commented LoadBalancer recipe for teams that need
direct exposure.

Know what this is not: community edition means ONE server — no
replicas, no failover; availability is the StatefulSet rescheduling
the pod and the PVC surviving it. If the graph is critical enough to
need cluster availability, that is the enterprise edition (a license,
`accept_license_agreement`, and multiple KubernetesNeo4j resources
sharing a `cluster_name`) — not a bigger version of this preset.
Backups are also not declared here; run them operationally against
the bolt endpoint.

Change first: create the `neo4j-auth` Secret (key NEO4J_AUTH, value
`neo4j/<password>`) before applying — the chart reads it at template
time and the install fails without it. Then set the `storage_class`
placeholder (a literal class name, or a valueFrom reference to a
KubernetesStorageClass) and re-derive the memory split if you resize
the container: initial heap = max heap, page cache from what remains
after heap and roughly 2Gi of OS overhead.

See [02-production.yaml](./02-production.yaml) for the manifest.
