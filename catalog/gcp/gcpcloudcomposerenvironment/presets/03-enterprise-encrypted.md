# Enterprise Encrypted

A large Cloud Composer environment with the full security and
operations surface: CMEK encryption, a bring-your-own DAG bucket,
private networking, web UI allowlisting, disaster-recovery snapshots,
and data retention bounds.

## When to Use

- Organizations that mandate customer-managed encryption keys (CMEK)
- Regulated workloads (PCI-DSS, HIPAA) needing private access and
  auditable retention policies
- Large-scale production pipelines with generous resource needs
- Teams that manage the DAG bucket as its own governed resource

## Key Configuration

- **ENVIRONMENT_SIZE_LARGE + HIGH_RESILIENCE** — maximum managed
  capacity with multi-zone redundancy (bump to
  `ENVIRONMENT_SIZE_EXTRA_LARGE` for the heaviest fleets)
- **kmsKeyName** — every Composer-managed resource (GKE nodes, Cloud
  SQL, storage) encrypted with your key; immutable after creation
- **storageBucket** — DAGs, plugins, and data live in a bucket you
  govern (lifecycle rules, IAM, CMEK) instead of the auto-created one;
  immutable after creation
- **Private endpoint with VPC peering** — no public Airflow UI
- **recoveryConfig** — daily 04:00 UTC environment snapshots for
  disaster recovery
- **dataRetentionConfig** — task logs go to both Cloud Logging and the
  environment bucket; on Composer 3, enable Airflow metadata retention
  (30-730 days) instead
- **labels** — merged beneath Planton's platform attribution labels

## What to Customize

- `kmsKeyName` — your KMS key; the Composer service agent needs
  `roles/cloudkms.cryptoKeyEncrypterDecrypter` on it, and the key must
  be in the environment's region
- `storageBucket` — your governed bucket (or drop the field to let
  Composer create one)
- `recoveryConfig.snapshotLocation` — a GCS path you control
- Workload sizes and worker min/max — sized here for heavy fleets;
  trim to your actual load
- `webServerNetworkAccessControl.allowedIpRanges` — your corporate and
  VPN ranges

## Important Notes

- Environment creation takes 25-45 minutes.
- CMEK, the storage bucket, and all networking are immutable — settle
  them before the first deploy.
- LARGE sizing with these workload allocations carries significant
  cost; monitor and adjust.

## Related Presets

- **01-dev-small** — minimal development environment
- **02-production-private** — private networking without CMEK or
  retention controls
