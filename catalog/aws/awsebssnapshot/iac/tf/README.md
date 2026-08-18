# AwsEbsSnapshot — Terraform/OpenTofu module

Manages one EBS snapshot as a three-way source union (`aws_ebs_snapshot` XOR `aws_ebs_snapshot_copy` XOR `aws_ebs_snapshot_import`) with fast snapshot restore (`aws_ebs_fast_snapshot_restore`) and cross-account grants (`aws_snapshot_create_volume_permission`) in-line.

Module facts worth knowing before editing:

- **Exactly one arm exists** — the spec's union CEL guarantees it, and all three arms expose the identical downstream surface through the `snapshot_*` locals.
- **Only the volume arm imports** — `aws_ebs_snapshot_copy`, `aws_ebs_snapshot_import`, and `aws_snapshot_create_volume_permission` ship NO importer at the provider (declared honestly in the import catalog).
- **Tiering dials update in place** — `storage_tier`, `permanent_restore`, `temporary_restore_days` are the only in-place surface; every source field replaces the snapshot.
- **Fast snapshot restore bills per zone-hour while enabled** — one resource per zone, keyed by zone name; treat the list as a deliberate cost decision.
- **Share grants are per-account resources** — keyed by account id; encrypted snapshots additionally need the KMS key shared with the same accounts (out of band — grants alone are not enough).
- **Imports run through VM Import/Export** — the `vmimport` service role (or `role_name`) must pre-exist with the documented trust + S3 policy, or the create fails at AWS after several minutes.

Outputs mirror the Pulumi module key-for-key: `snapshot_id` (import ID, volume arm), `snapshot_arn`, `owner_id`, `volume_size_gb`.
