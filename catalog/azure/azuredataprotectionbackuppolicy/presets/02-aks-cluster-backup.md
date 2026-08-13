# AKS Cluster Backup

This preset creates the Kubernetes cluster policy: backups every four hours with two weeks of default retention, and the first backup of each day kept for eight weeks. The modern-backup capability most teams adopt Data Protection for.

## When to Use

- Protecting AKS cluster state (resources and persistent volumes) on an operational cadence
- Clusters whose recovery point objective is measured in hours, not days
- The retention shape to start from before tuning cadence and keeps to the workload

## Key Configuration Choices

- **`backupRepeatingTimeIntervals: R/.../PT4H`** -- every four hours; tighten or relax to the cluster's recovery-point objective
- **`defaultRetentionRule.lifeCycles` on `OperationalStore`** -- AKS retains near the source (the service's surface today; the vocabulary widens when the AKS vault tier lands)
- **`retentionRules[daily]`** -- the first backup of each day kept 8 weeks, layered over the 14-day default
- **The AKS backup extension and its permissions ride the backup INSTANCE**, not this policy -- the policy is pure configuration

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-backup-vault>` | The AzureDataProtectionBackupVault the policy lives in | Your backup vault resource's name |

The policy is free -- cost follows the protected clusters and their backup storage.
