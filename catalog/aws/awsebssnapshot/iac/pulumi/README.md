# AwsEbsSnapshot — Pulumi module (Go)

Manages one EBS snapshot as a three-way source union (`ebs.Snapshot` XOR `ebs.SnapshotCopy` XOR `ebs.SnapshotImport`) with fast snapshot restore (`ec2.FastSnapshotRestore`) and cross-account grants (`ec2.SnapshotCreateVolumePermission`) in-line.

Module facts worth knowing before editing:

- **Exactly one arm exists** — the spec's union CEL guarantees it; all three arms feed the same exported outputs.
- **Only the volume arm imports** — the copy, import, and grant resources ship NO importer at the provider (declared honestly in the import catalog).
- **Tiering dials update in place** — `StorageTier`, `PermanentRestore`, `TemporaryRestoreDays`; every source field replaces the snapshot. The three per-arm apply helpers exist because the bridge generates three distinct args structs for the same dials.
- **Fast snapshot restore bills per zone-hour while enabled** — one resource per zone (`fast-restore-{zone}`).
- **Share grants are per-account resources** (`share-{account}`); encrypted snapshots additionally need the KMS key shared out-of-band.
- **Imports run through VM Import/Export** — the `vmimport` service role (or `RoleName`) must pre-exist or the create fails at AWS after several minutes.

Outputs mirror the Terraform module key-for-key: `snapshot_id` (import ID, volume arm), `snapshot_arn`, `owner_id`, `volume_size_gb`.
