<p align="center">
  <img src="logo.svg" alt="AWS Backup Restore Testing Plan" width="80"/>
</p>

# AWS Backup Restore Testing Plan

Manage an [AWS Backup restore testing plan](https://docs.aws.amazon.com/aws-backup/latest/devguide/restore-testing.html)
— scheduled, automated restore tests that prove recovery points
actually restore, with the per-resource-type selections folded in.

## What Gets Managed

- **The testing plan** (`spec.plan_name` is the AWS name — restore
  testing names forbid hyphens and periods, so `metadata.name` stays
  Planton-side): the test schedule (cron + timezone), the start
  window, and the **recovery point selection** (latest vs random
  within a lookback window, across included/excluded vaults, for
  snapshot and/or continuous points).
- **Selections** folded as name-keyed entries: one per protected
  resource type, with the IAM role AWS Backup assumes, coverage by
  explicit ARNs XOR tag conditions, restore metadata overrides, and
  the validation window.

AWS runs each test into a temporary copy and deletes it after the
validation window — tests bill as regular restores plus the temporary
resource's runtime.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
