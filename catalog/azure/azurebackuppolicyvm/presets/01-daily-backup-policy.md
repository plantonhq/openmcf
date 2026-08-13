# Daily Backup Policy

This preset creates the everyday production policy: one nightly backup with the classic grandfather-father-son retention ladder -- 30 dailies, 12 weekly Sundays, 12 first-Sunday monthlies, 7 first-Sunday-of-January yearlies. The right default for most VMs.

## When to Use

- The first policy a vault needs -- one schedule most VMs share
- Workloads whose recovery point objective tolerates up to 24 hours
- Standard compliance shapes ("a month of dailies, a year of monthlies, seven years of yearlies")

## Key Configuration Choices

- **`policyType: V2`** -- the enhanced generation, stated explicitly: it unlocks hourly schedules and longer instant restore later, and the field is ForceNew (a replacement re-binds every protected VM)
- **`time: "23:00"`** -- one dial for the backup start AND every retention layer's timestamp; must land on the hour or half past
- **Week-of-month retention grammar** -- "First Sunday" for monthlies/yearlies; the alternative month-days grammar (`days`/`includeLastDays`) is mutually exclusive with it
- **30/12/12/7 ladder** -- each kept point is incremental backup storage; this shape balances restore granularity against cost

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The vault's AzureResourceGroup | Your resource group resource's name |
| `<your-recovery-services-vault>` | The AzureRecoveryServicesVault the policy lives in | Its name output wires automatically |

For multi-year retention at lower cost, add a `tieringPolicy` with `mode: TierRecommended` -- Azure moves archivable points to archive-tier prices.
