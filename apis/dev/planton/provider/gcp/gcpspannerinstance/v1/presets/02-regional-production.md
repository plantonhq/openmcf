# Regional Production

Provisions a production-ready Cloud Spanner instance with fixed capacity (one node), ENTERPRISE edition, and automatic backup schedules for new databases. Suitable for production workloads with predictable capacity in a single GCP region.

## When to Use

- Production applications with predictable, steady traffic patterns
- Workloads that do not need multi-region replication (single-region is significantly cheaper)
- Teams that prefer explicit capacity management over autoscaling

## Key Configuration

- **1 node** — roughly 10,000 QPS reads, 2,000 QPS writes, 10 TB storage; capacity changes apply online (raise `numNodes`, or switch to `processingUnits` for finer granularity below one node)
- **ENTERPRISE edition** — granular sizing, asymmetric autoscaling, and incremental backups become available
- **AUTOMATIC backup schedule** — GCP attaches a default backup schedule to each new database; use `GcpSpannerBackupSchedule` resources for explicit control
- **Regional config** — all replicas in one region; 99.99% availability SLA

## Customization Notes

- Replace `config` with your region's configuration (e.g. `regional-europe-west1`); this choice is immutable
- `metadata.name` doubles as the instance name when `instanceName` is omitted (6-30 characters)
- `project_id` falls back to the provider's default project; set `projectId` to target another project
- `labels` merge with Planton's platform labels for cost attribution

## Related Presets

- **01-free-instance** — zero-cost instance for development
- **03-autoscaling-production** — autoscaling for variable workloads
