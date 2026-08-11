# Hourly Enhanced Policy

This preset creates a low-RPO policy on the V2 (enhanced) generation: a backup every 4 hours inside a 12-hour working-day window, a week of instant-restore snapshots, and age-based archive tiering. For VMs whose data loses real money by the hour.

## When to Use

- Databases and stateful applications with a recovery point objective under a day
- Workloads whose backup churn must stay inside a defined window
- Long monthly retention that should ride archive-tier prices

## Key Configuration Choices

- **`frequency: Hourly` needs `policyType: V2`** -- V1 policies back up at most once a day (validated at manifest time)
- **The window arithmetic** -- `time: 08:00` + `hourDuration: 12` bounds backups to 08:00-20:00; `hourInterval: 4` spaces three backups inside it. The duration must be a multiple of the interval
- **`instantRestoreRetentionDays: 7`** -- a week of snapshot-speed restores (V1 caps this at 5)
- **`TierAfter` 3 months** -- every point past the age archives; `TierRecommended` is the alternative that lets Azure pick cost-optimal points
- **Month-days retention grammar** -- "the 1st" via `days: [1]`; mutually exclusive with the week-of-month grammar

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The vault's AzureResourceGroup | Your resource group resource's name |
| `<your-recovery-services-vault>` | The AzureRecoveryServicesVault the policy lives in | Its name output wires automatically |

Hourly points multiply retention math -- `retentionDaily.count: 30` here means 30 DAYS of hourly points. Watch the storage line the first month.
