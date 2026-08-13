# Overview

The **AzureBackupPolicyVm** component creates an Azure Backup policy for IaaS virtual machines -- WHEN VMs under it are backed up (the schedule) and HOW LONG each backup is kept (retention layered daily/weekly/monthly/yearly, like grandfather-father-son tape rotation). The policy is a free configuration object, an ARM child of its Recovery Services vault.

## Purpose

- **Backup schedules as declarative infrastructure**: frequency, time, and the whole retention ladder -- reviewed and versioned, never clicked together.
- **Front-loaded service contracts**: the frequency/retention coupling, the hourly-needs-V2 rule, and Azure's daily-retention floor all validate at manifest time instead of failing minutes into an apply.
- **Cost-aware retention**: archive tiering moves aged recovery points to cheaper storage; instant-restore days and layered retention are the levers backup bills are made of.
- **Typed references**: the vault wires by name reference -- chart-ready.

## Key Features

- Full azurerm v5 surface: Hourly (V2)/Daily/Weekly schedules with the hourly window dials, all four retention layers with both selection grammars (week-of-month and month-days), archive tiering (recommended or age-based), instant-restore retention and its named snapshot resource group, crash-only consistency (V2), timezone.
- The provider's whole CustomizeDiff lattice mirrored as CEL: a manifest that validates renders a schedule ARM accepts.
- Azure's undocumented-until-apply rules encoded: daily retention is 1 or 7+ (the service rejects 2-6).

## Use Cases

- **The everyday daily policy**: one nightly backup, 30 days of dailies, 12 monthly, 7 yearly.
- **Enhanced hourly protection**: V2 policies backing up every 4 hours for low-RPO workloads.
- **Archive-tiered long retention**: multi-year compliance retention at archive-tier prices.

## Future Enhancements

- The file-share backup policy kind extends the same retention grammar to Azure Files as its contract lands in the catalog.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
