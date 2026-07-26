# Production preset

Temporal sized for real workloads on the composed-PostgreSQL story:
replicated frontend/matching/worker, three history replicas (shard
ownership redistributes across them — history is where workflow state
and cache pressure live), a replicated stateless UI, ServiceMonitor
metrics, and 7-day retention on the `default` Temporal namespace.

The database is a same-namespace KubernetesPostgres named
`temporal-db` — give IT the production posture in its own manifest
(2+ instances for failover, real storage) and bootstrap both databases
there: `temporal` at initdb, `temporal_visibility` via post_init_sql.
The credential never appears here; the reference resolves to the
operator-maintained Secret, which survives database failovers.

Two numbers deserve respect: `num_history_shards` is IMMUTABLE after
first install (the preset states the 512 default explicitly so it is a
decision, not an accident), and `max_conns` is PER SERVICE INSTANCE —
nine server pods x 20 connections each needs the database's
max_connections sized for ~180 from Temporal alone.

Change first: archive closed histories (`archival` with an S3/GCS
bucket) before retention starts deleting anything you might need;
expose the frontend/UI through first-class exposure kinds over the
exported service handles — nothing in this kind does ingress.

See [02-production.yaml](./02-production.yaml) for the manifest.
