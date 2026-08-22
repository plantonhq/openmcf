# Backup Operator Key

This preset mints an account-wide Spaces credential for the rare tool that legitimately needs every bucket -- a backup system sweeping all storage, or an admin migration job.

## When to Use

- Backup/restore tooling that must reach every bucket, including ones created later
- One-off account-wide migrations or audits

## Key Configuration Choices

- **fullaccess with no bucket** -- the provider's account-wide grammar; validation rejects a fullaccess grant that names a bucket.
- **Treat it as an admin credential** -- store the write-once secret in a manager, rotate on a schedule (rotation = destroy and recreate), and never hand this key to a workload that only needs its own bucket (use the per-bucket preset for that).

## What You Get

One credential that unlocks all Spaces data in the account -- maximum blast radius, so reserve it for operators, not applications.
