<p align="center">
  <img src="logo.svg" alt="AWS Backup Settings" width="80"/>
</p>

# AWS Backup Settings

Manage [AWS Backup account and region settings](https://docs.aws.amazon.com/aws-backup/latest/devguide/service-opt-in.html)
— the account's global settings (cross-account backup) and the
region's resource-type opt-in and management preferences.

## What Gets Managed

- **The global arm** (account-wide — one instance across all regions
  should own it): the global settings map, notably
  `isCrossAccountBackupEnabled`.
- **The region arm** (one instance per region): which resource types
  AWS Backup protects in the region
  (`resource_type_opt_in_preference`) and which it fully manages
  (`resource_type_management_preference`).

This is a SETTINGS SINGLETON: `metadata.name` never reaches AWS — the
global arm's identity is the account, the region arm's identity is the
region. **Destroy is a no-op on both arms**: whatever was last applied
stays in effect; to revert a setting, apply the desired value before
destroying.

These settings are deliberately NOT fields on
[AwsBackupVault](../awsbackupvault) or
[AwsBackupPlan](../awsbackupplan) — multiple vaults/plans would fight
over the one settings object.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
