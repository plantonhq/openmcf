---
title: "Daily Disk Backup"
description: "This preset creates the everyday managed-disk policy: daily incremental snapshots at 02:00 UTC, a week of default retention, and the first backup of each week kept for 90 days. The right starting..."
type: "preset"
rank: "01"
presetSlug: "01-daily-disk-backup"
componentSlug: "data-protection-backup-policy"
componentTitle: "Data Protection Backup Policy"
provider: "azure"
icon: "package"
order: 1
---

# Daily Disk Backup

This preset creates the everyday managed-disk policy: daily incremental snapshots at 02:00 UTC, a week of default retention, and the first backup of each week kept for 90 days. The right starting point for protecting VM data disks and standalone managed disks.

## When to Use

- The first disk policy a vault needs -- the plain daily-snapshot contract
- VM data disks and standalone managed disks with a week of operational recovery
- As the safe, boring install manifest for backup-instance testing and charts

## Key Configuration Choices

- **`backupRepeatingTimeIntervals: R/.../P1D`** -- daily, phase-anchored at 02:00 UTC by the interval's date-time part
- **`defaultRetentionDuration: P7D`** -- a week of point-in-time recovery for every snapshot
- **`retentionRules[weekly]`** -- tags the first backup of each week for a 90-day keep (grandfather-father-son, one layer)
- **Remember policies are IMMUTABLE** -- changing anything here replaces the policy; instances re-bind to the replacement

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-backup-vault>` | The AzureDataProtectionBackupVault the policy lives in | Your backup vault resource's name |

The policy is free -- cost follows the protected disks and their snapshot storage.
