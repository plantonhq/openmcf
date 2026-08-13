# Azure Backup Policy (File Share) -- Operational Guide

Judgment that saves real time when designing file-share backup policies. The field reference lives in the API Explorer; this is the operational layer above it.

## Snapshot or vault-standard is THE decision

`backupTier: snapshot` (the default) keeps backups as share snapshots INSIDE the storage account -- restores are fast, but the backups share the account's fate: delete or compromise the account and the "backups" go with it. `vault-standard` additionally copies backups into the vault, which is what "backup" usually means to an auditor. Pick vault-standard for anything whose loss would matter; keep snapshot for dev shares where speed beats durability.

## The snapshot budget is 200, and retention math spends it

Azure Files holds at most 200 snapshots per share. Your retention ladder SPENDS that budget: 30 dailies + 12 weeklies + 12 monthlies + 10 yearlies is 64 snapshots -- fine; 200 dailies alone maxes it out and the policy fails at create. The bounds are already shorter than VM policies (daily/weekly cap at 200, monthly 120, yearly 10) -- design inside them.

## Hourly is a window, not a metronome

An hourly schedule runs inside a WINDOW: `hourly.startTime` opens it, `hourly.windowDuration` bounds it (4-24 hours), `hourly.interval` spaces backups within it (4, 6, 8 or 12 hours). Interval 4 with a 12-hour window gives three backups inside business hours and none at night. There is no `time` field on hourly schedules -- the window replaces it (the manifest validation enforces the shape).

## vault-standard's snapshot retention is a strict bound

`snapshotRetentionInDays` (vault-standard only) keeps local snapshots alongside the vaulted copies for fast operational restores. It must be STRICTLY LESS than `retentionDaily.count` -- the provider rejects equality. Five local days against thirty vaulted dailies is the everyday shape.

## Changing a policy changes every share under it

A policy update applies to all protected shares at the next scheduled backup -- there is no per-share override. Shortening retention DELETES existing recovery points beyond the new horizon on the service side; treat retention reductions as a deliberate, announced operation.

## The policy is step one of three

A policy alone protects nothing. The full chain: register the share's storage account with the vault (AzureBackupContainerStorageAccount), then bind each share to this policy (AzureBackupProtectedFileShare). The policy can serve many shares across many registered accounts in the same vault.
