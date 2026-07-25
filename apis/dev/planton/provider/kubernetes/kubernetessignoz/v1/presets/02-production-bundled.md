# Production on the bundled ClickHouse

A production-shaped single-cluster observability platform: replicated
ClickHouse (one shard, two replicas over a 3-node ZooKeeper quorum),
autoscaled ingestion, alert emails over authenticated SMTP, and an
external URL so alert links land on the real UI hostname.

Three things to know before applying:

- **The database volume is the retention budget.** Retention itself is
  set inside SigNoz (UI → Settings → General); this preset's 200Gi per
  replica is what that retention has to fit into. Watch usage and expand
  the volume class permitting.
- **The SMTP password is a reference** to a Secret you create first
  (`signoz-smtp-auth`, key `password`) — it is wired as an env-from-secret
  and never lands in rendered configuration.
- **Multi-replica layouts need a multi-node cluster**, and the layout is
  a day-0 choice — the chart marks reshaping a live installation
  experimental.

**When to use:** a team or org running its observability inside one
cluster, wanting the one-product experience with HA storage.

**When to move on:** when the database deserves its own lifecycle
(scaling, cold tiers, deep tuning, its own operations), move the storage
to a KubernetesClickHouse with `03-external-clickhouse`.
