# GcpSpannerInstance — Research and Design Documentation

## 1. Cloud Spanner's Resource Hierarchy

Spanner uses a three-level hierarchy, and Planton models the first two as separate composable kinds:

1. **Instance** (this component) — the unit of compute and storage allocation. Pins the geographic topology (instance configuration) and the capacity envelope every database shares. Infrastructure-shaped: rarely recreated, shared by many databases.
2. **Database** ([GcpSpannerDatabase](../../gcpspannerdatabase/v1)) — schemas, data, encryption, and retention. Application-shaped: created and destroyed as applications evolve.
3. **Tables/indexes** — schema objects inside a database; ongoing DDL belongs to migration tooling, not IaC.

Backup schedules are a fourth API object ([GcpSpannerBackupSchedule](../../gcpspannerbackupschedule/v1)) — many per database with their own lifecycle, so they are their own kind rather than a bundled block.

## 2. Capacity Model

- **Nodes** — coarse units: ~10,000 QPS reads, ~2,000 QPS writes, 10 TB storage each.
- **Processing units** — 1 node = 1000 PUs; below 1000, multiples of 100 (the smallest billable Spanner is 100 PUs).
- **Autoscaling** — Spanner adjusts capacity within explicit min/max bounds based on high-priority CPU and storage targets. Google's guidance: CPU target 65% for regional, ~45% for multi-region (failover headroom); storage 80%. While enabled, fixed capacity fields become read-only reflections.
- **Asymmetric autoscaling options** — multi-region instances can bound one replica location's node range independently (a read-heavy region scales alone). Requires ENTERPRISE or ENTERPRISE_PLUS. The provider wraps the replica location in a single-field `replica_selection` block; the spec flattens it to `replica_location` — a pure API-shape artifact, not a semantic choice.

All capacity operations are online: no capacity change recreates the instance.

## 3. Terraform Provider Floor

Designed from `google_spanner_instance` on the released Terraform Google provider 6.x line (`~> 6.0`); the Spanner surface is fully GA (the GA and beta provider schemas are identical). Both engines enable `spanner.googleapis.com` before creating the instance, and both build the `instance_id` output from the provider-resolved project so the ambient-project fallback stays honest.

### Field coverage

| Provider surface | Modeled | Notes |
|---|---|---|
| `config`, `display_name`, `name`, `project` | ✅ | name defaults to `metadata.name`; project falls back to the ambient provider project |
| `num_nodes` / `processing_units` / `autoscaling_config` | ✅ | exactly-one enforced by CEL pre-deploy (the provider's ConflictsWith set) |
| `autoscaling_limits`, `autoscaling_targets` | ✅ | same-unit + max≥min + 0-100 bounds enforced pre-deploy |
| `asymmetric_autoscaling_options` | ✅ | replica location + node-range overrides (the released surface) |
| `instance_type`, `edition`, `default_backup_schedule_type` | ✅ | FREE_INSTANCE conflict set enforced by CEL |
| `labels` | ✅ | user labels merged beneath Planton attribution labels, both engines identically |
| `force_destroy` | ✅ | destroy-time backup deletion lever |

### Recorded skips (evidence-based)

| Feature | Reason |
|---|---|
| `deletion_policy` | Client-side Terraform lever (PREVENT/ABANDON) that conflicts with managed destroy semantics; catalog-wide exclusion. |
| Per-replica autoscaling target overrides / disable flags, processing-unit override limits | Present on the provider's main branch only; the released 6.x asymmetric surface is replica location + node bounds. Add when released. |
| `autoscaling_targets.total_cpu_utilization_percent` | Main-branch-only; the released 6.x targets are high-priority CPU + storage. |
| Custom instance configurations (`google_spanner_instance_config`) | Own resource for user-managed replica topologies — a genuine enterprise niche; candidate kind on concrete pull. |
| Instance partitions (`google_spanner_instance_partition`) | Own resource for geo-partitioning (ENTERPRISE_PLUS); candidate kind on concrete pull. |
| Instance IAM trio (`google_spanner_instance_iam_*`) | Resource-scoped IAM stays unmodeled catalog-wide; grants compose via IAM kinds. |

## 4. Editions

| Feature | STANDARD | ENTERPRISE | ENTERPRISE_PLUS |
|---|---|---|---|
| Regional SLA | 99.99% | 99.99% | 99.99% |
| Multi-region SLA | N/A | 99.999% | 99.999% |
| Asymmetric autoscaling | — | ✅ | ✅ |
| Incremental backups | — | ✅ | ✅ |
| Pricing | Lowest | Medium | Highest |

Upgrades apply in place; downgrades require disabling higher-edition features first.

## 5. Downstream Composition

The `instance_name` output is the composition key:

```
GcpProject (project_id)
  └── GcpSpannerInstance (instance_name)
        ├── GcpSpannerDatabase (instance)
        │     └── GcpSpannerBackupSchedule (database)
        └── GcpSpannerBackupSchedule (instance)
```

`default_backup_schedule_type: AUTOMATIC` gives every new database a GCP-managed default schedule; explicit GcpSpannerBackupSchedule resources are the full-control path.

## 6. Immutability

ForceNew (recreate on change): `instance_name`, `config`, `project_id`. Everything else — capacity, autoscaling, edition, backup schedule type, labels, display name — updates in place.
