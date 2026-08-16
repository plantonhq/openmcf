# Backup Coverage Audit

This preset creates the first-questions framework: are resources
actually covered by a backup plan, and are recovery points retained at
least 35 days.

## When to Use

- The first Audit Manager framework in an account — coverage and
  retention are what every compliance review asks first
- Catching resources that drifted out of backup-plan selections

## What You Get

- Continuous evaluation of plan coverage across the account
- A 35-day minimum-retention check on every recovery point
- Results as Config rules in the Backup console's Audit Manager view

## Customize

- Tune `requiredRetentionDays` to your regime's floor
- Add `BACKUP_RECOVERY_POINT_ENCRYPTED` and
  `BACKUP_RECOVERY_POINT_MANUAL_DELETION_DISABLED` for the
  security-posture pair
- Scope controls to tagged subsets with `scope` (one tag pair max)
- Remember: the region needs an ACTIVE Config recorder or deployment
  lands FAILED
