# CMEK-Encrypted Database

Provisions a compliance-grade Spanner database: customer-managed encryption (CMEK) by reference to a `GcpKmsKey`, GCP API-side drop protection, point-in-time recovery, and an explicit UTC time zone.

## When to Use

- Regulated workloads that require customer-owned encryption keys
- Databases whose deletion must be blocked at the GCP API level, not just in IaC
- Compliance regimes that audit key rotation and access separately from data access

## Key Configuration

- **CMEK by reference** — `kmsKeyName` resolves the KMS key's fully qualified `key_id` output; the key must live in the same location as the instance configuration. Encryption posture is immutable.
- **Multi-region instances** use `kmsKeyNames` instead — one key per region of the instance configuration (exactly one shape may be set)
- **enableDropProtection** — while true, NO interface (console, gcloud, API, IaC) can delete the database, and the parent instance cannot be deleted either
- **Deletion protection ON by default** — the IaC-side guard on top of the API-side lock

## Customization Notes

- Grant the Spanner service agent (`service-{project_number}@gcp-sa-spanner.iam.gserviceaccount.com`) the `roles/cloudkms.cryptoKeyEncrypterDecrypter` role on the key before creating the database
- To tear down: first set `enableDropProtection: false` (applies in place), then `deletionProtection: false`, then destroy

## Related Presets

- **01-basic-database** — Google Standard SQL database
- **02-postgresql-database** — PostgreSQL-dialect database
