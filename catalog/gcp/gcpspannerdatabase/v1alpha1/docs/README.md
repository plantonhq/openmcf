# GcpSpannerDatabase — Research and Design Documentation

## 1. Architecture: Instance vs Database vs Backup Schedule

```
Project
└── Spanner Instance (compute + topology)         GcpSpannerInstance
    ├── Database A (schema + data + encryption)   GcpSpannerDatabase
    │   ├── Tables / Indexes / Views              (migration tooling)
    │   └── Backup Schedule(s)                    GcpSpannerBackupSchedule
    └── Database B ...
```

A single instance hosts many databases, all sharing its compute capacity. Databases carry independent schemas, encryption posture, retention, and lifecycle settings. Backup schedules are separate API objects — many per database, own lifecycle — so they are their own kind rather than a bundled block.

The instance/database split reflects production reality: platform teams manage instances (capacity, topology), application teams manage databases (schema, retention, encryption), and both compose in infra charts (one instance, N databases).

## 2. Terraform Provider Floor

Designed from `google_spanner_database` on the released Terraform Google provider 6.x line (`~> 6.0`); the Spanner surface is fully GA (GA and beta schemas identical). Both engines enable `spanner.googleapis.com` before creating the database, and both build `database_id` from the created resource's resolved attributes so the ambient-project fallback stays honest.

### Field coverage

| Provider surface | Modeled | Notes |
|---|---|---|
| `instance`, `name`, `project` | ✅ | instance by FK to the GcpSpannerInstance `instance_name` output; name defaults to `metadata.name`; project ambient fallback |
| `database_dialect` | ✅ | GOOGLE_STANDARD_SQL / POSTGRESQL, permanent |
| `ddl` | ✅ | append-only lifecycle documented; the provider forces recreation when an existing entry changes |
| `version_retention_period` | ✅ | 1h-7d point-in-time recovery window; applied via ALTER DATABASE |
| `enable_drop_protection` | ✅ | GCP API-side lock (also blocks parent-instance deletion) |
| `encryption_config.kms_key_name` / `kms_key_names` | ✅ | regional XOR multi-region CMEK enforced pre-deploy; both FK → GcpKmsKey `key_id` (the fully qualified crypto-key path the Spanner API requires) |
| `default_time_zone` | ✅ | IANA zone; applied via ALTER DATABASE at create |
| `deletion_protection` | ✅ | IaC-side guard, spec default TRUE, set explicitly on BOTH engines for identical destroy semantics |

### The two deletion guards, deliberately both modeled

- `deletion_protection` is client-side: the engines refuse to destroy the resource. It exists on both providers with default true — the spec surfaces it (default TRUE) and both modules always set it explicitly so neither engine's provider default can drift the destroy behavior.
- `enable_drop_protection` is server-side: GCP refuses deletion through every interface and blocks deleting the parent instance. Compliance-grade; default false because it makes even intentional teardown a two-step operation.

### Recorded skips (evidence-based)

| Feature | Reason |
|---|---|
| `deletion_policy` | Client-side Terraform lever (PREVENT/ABANDON) that conflicts with managed destroy semantics; catalog-wide exclusion. |
| Database IAM trio (`google_spanner_database_iam_*`) | Resource-scoped IAM stays unmodeled catalog-wide; grants compose via IAM kinds. |
| Ongoing schema management | Tables/indexes beyond the initial DDL belong to migration tooling (Liquibase, Flyway), not IaC. |

## 3. CMEK Notes

- The key must live in the same location as the instance configuration; multi-region configurations need one key per region (`kms_key_names`).
- The Spanner service agent (`service-{project_number}@gcp-sa-spanner.iam.gserviceaccount.com`) needs `roles/cloudkms.cryptoKeyEncrypterDecrypter` on every key before creation.
- Encryption posture is immutable — switching between Google-managed and CMEK (or rotating to a different key resource) recreates the database.

## 4. Immutability

ForceNew (recreate on change): `instance`, `database_name`, `database_dialect`, `encryption_config`, `project_id`, and any edit to an existing `ddl` entry. Mutable in place: `version_retention_period`, `enable_drop_protection`, appended `ddl`, `deletion_protection`.

## 5. Downstream Composition

`database_name` is the composition key for GcpSpannerBackupSchedule; `database_id` is the handle for IAM bindings and API callers. In the spanner application pattern:

```
GcpSpannerInstance (instance_name)
  └── GcpSpannerDatabase (database_name)
        └── GcpSpannerBackupSchedule
```
