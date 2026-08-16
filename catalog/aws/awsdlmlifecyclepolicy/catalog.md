# AWS Data Lifecycle Manager Policy

Snapshot hygiene on autopilot: tag your volumes once, and DLM creates the backups on schedule, keeps exactly the retention you declared, replicates to other regions, archives the old ones, and deletes what expired — no cron jobs, no forgotten snapshots billing forever.

## What Gets Managed

- The policy, in one of two modes: AWS's simplified default posture (every volume or instance, daily-ish, N-day retention, with exclusions) or the full custom engine (tag-targeted, up to four schedules with create/retain/archive/copy/share/deprecate rules).
- Its enabled/disabled state and the IAM role DLM acts through.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with DLM and EC2 permissions.

### AWS Prerequisites

- An execution role trusting `dlm.amazonaws.com` with snapshot permissions — AWS's `AWSDataLifecycleManagerDefaultRole` (create once with `aws dlm create-default-role`) or your own AwsIamRole.

## After You Deploy

- Tag resources with the target tags — coverage is dynamic; new volumes carrying the tag are covered the next time a schedule fires.
- Watch the policy's state: AWS flips it to ERROR when the role loses permissions.

## Common Changes

- Add a tier: a second schedule (weekly, monthly) in the same policy — up to four.
- DR replication: a `cross_region_copy_rules` entry with the destination region's key.
- Pause without deleting: `disabled: true` — schedules stop firing, existing snapshots stay.
