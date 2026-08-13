# Overview

The **AzureBackupPolicyFileShare** component creates an Azure Backup policy for Azure Files shares -- WHEN shares under it are backed up (the schedule) and HOW LONG each backup is kept (retention layered daily/weekly/monthly/yearly, like grandfather-father-son tape rotation). The policy is a free configuration object, an ARM child of its Recovery Services vault.

## Purpose

- **Backup schedules as declarative infrastructure**: frequency, time, and the whole retention ladder -- reviewed and versioned, never clicked together.
- **Front-loaded service contracts**: the schedule-shape rule (Daily names a time; Hourly configures a window) and the vault-standard snapshot-retention bound validate at manifest time instead of failing minutes into an apply.
- **A deliberate tier choice**: `backupTier` decides whether backups live as share snapshots in the storage account (fast, but sharing the account's fate) or additionally as vaulted copies (surviving account deletion and compromise).
- **Typed references**: the vault wires by name reference -- chart-ready.

## Key Features

- Full azurerm v5 surface: Daily/Hourly schedules with the hourly window dials (interval 4/6/8/12, start time, duration), all four retention layers with both selection grammars (week-of-month and month-days), the snapshot/vault-standard tier switch with local snapshot retention, timezone.
- The provider's schedule-shape and tier contracts mirrored as CEL: a manifest that validates renders a schedule ARM accepts.
- File-share retention bounds encoded honestly: daily/weekly 1-200, monthly 1-120, yearly 1-10 -- shorter than VM retention by the service's own design.

## Use Cases

- **The everyday share policy**: one nightly snapshot, 30 days of dailies, 12 monthly, a few yearly.
- **Vaulted protection**: `vault-standard` copies backups into the vault so a deleted or compromised storage account cannot take its backups with it.
- **Low-RPO file shares**: hourly windows backing up every 4 hours during business hours.

## Future Enhancements

- Azure Files vaulted backup gains capabilities release by release (cross-region restore among them) -- the tier vocabulary extends additively as azurerm models them.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
