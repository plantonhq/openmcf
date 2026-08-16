# DR Cold Storage

This preset creates the disaster-recovery tier: weekly backups of
critical-tagged resources, archived to cold storage after 30 days,
kept a year, with a copy in a second vault (typically another region).

## When to Use

- Year-scale retention where warm storage would dominate the bill
- Cross-region (or cross-account) copies for real DR posture

## What You Get

- A weekly Sunday 03:00 UTC rule with cold-storage transition at 30
  days and expiry at 365 (honoring AWS's 90-day cold-storage minimum)
- A copy action into a second vault with its own one-year lifecycle
- A `tier: critical` tag selection assuming the referenced role

## Customize

- Point `destinationVaultArn` at a vault in another REGION for
  region-loss DR (same-region copies only protect against vault-level
  mistakes)
- Add `optInToArchiveForSupportedResources: true` to use the EBS
  archive tier instead of deletion
- Add a second rule at daily cadence with short retention for the
  operational-restore tier — one plan can carry both
