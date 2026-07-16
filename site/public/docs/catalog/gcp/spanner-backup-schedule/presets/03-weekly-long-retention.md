---
title: "Weekly Long-Retention Archive"
description: "Creates a full backup every Sunday and keeps each one for 366 days — the maximum Spanner allows — encrypted with an explicit customer-managed key. The compliance-archive pattern that runs alongside a..."
type: "preset"
rank: "03"
presetSlug: "03-weekly-long-retention"
componentSlug: "spanner-backup-schedule"
componentTitle: "Spanner Backup Schedule"
provider: "gcp"
icon: "package"
order: 3
---

# Weekly Long-Retention Archive

Creates a full backup every Sunday and keeps each one for 366 days — the maximum Spanner allows — encrypted with an explicit customer-managed key. The compliance-archive pattern that runs alongside a shorter-cadence operational schedule.

## When to Use

- Regulatory or audit requirements measured in months, not days
- As the second schedule on a database whose operational restore points come from a daily or incremental schedule
- When backup copies must be encrypted with a specific, auditable KMS key

## Key Configuration

- **Weekly cadence** — `0 4 * * 0`, Sundays at 04:00 UTC
- **366-day retention** — `31622400s`, the Spanner maximum
- **Explicit CMEK** — `CUSTOMER_MANAGED_ENCRYPTION` with a `GcpKmsKey` reference; the Spanner service agent needs `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key. For multi-region instances use `kmsKeyNames` (one key per region) instead.

## Customization Notes

- Schedules are many-per-database: deploy this beside 01 or 02, not instead of them
- Retention changes apply only to backups created after the change
- Omit `encryptionConfig` entirely to inherit the database's own encryption posture

## Related Presets

- **01-daily-full-backups** — the operational baseline
- **02-incremental-enterprise** — cost-efficient high-frequency restore points
