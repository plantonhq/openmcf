# Production HA preset

Three-node Galera cluster (quorum-safe synchronous replication) with
three HAProxy replicas, nightly XtraBackup to S3, dedicated PITR
storage, zone-spread anti-affinity, and cert-manager TLS via a
ClusterIssuer.

Use when the workload needs MySQL with zero-data-loss replication and
scheduled off-cluster backups — the standard production posture for
stateful services on Kubernetes.

See [02-production-ha.yaml](./02-production-ha.yaml) for the manifest.
