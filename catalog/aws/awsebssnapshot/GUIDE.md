# AwsEbsSnapshot — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Snapshots are incremental — deletion is not loss

Each snapshot stores only blocks changed since the previous one, and deleting an intermediate snapshot is safe: AWS re-parents the blocks. What you pay for is the unique block set, not the snapshot count — prune by policy (DLM), not by fear.

## Fast snapshot restore is a running meter

FSR bills per snapshot per zone-HOUR while enabled, whether or not anyone restores. Enable it for the restore-latency-critical window (a migration day, a DR drill), then remove the zones — the spec's list makes both directions one edit.

## The import arm fails late without the vmimport role

VM Import/Export validates its service role only when the task runs — a missing/misconfigured `vmimport` role fails the create after several minutes, not at plan. Create the role once per account with AWS's documented trust and S3 policy before the first import.

## Sharing encrypted snapshots is two grants, not one

`share_with_account_ids` grants createVolumePermission, but an encrypted snapshot is useless to the peer until the KMS key ALSO grants them decrypt — two resources, two owners. The unencrypted default aws/ebs key can never be shared: re-encrypt under a customer key first (the copy arm).

## Archive math has a floor

Archived snapshots bill a 90-day minimum: archiving something you will delete next week costs MORE than keeping it standard. Archive is for compliance tails, not short-lived backups.
