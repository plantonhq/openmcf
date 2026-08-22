# AwsEbsVolume — Terraform/OpenTofu module

Manages one EBS volume as a create-XOR-copy union (`aws_ebs_volume` XOR `aws_ebs_volume_copy`) with attachments (`aws_volume_attachment`) in-line.

Module facts worth knowing before editing:

- **Exactly one arm exists** — the spec's union CEL guarantees `copy_from` XOR the create-arm fields, and both arms expose the identical downstream surface through the `volume_id`/`volume_arn` locals.
- **A copy inherits placement and encryption** — `aws_ebs_volume_copy` has no zone or encryption arguments (the copy lands in the source's zone with the source's posture); only size/type/iops/throughput may be overridden.
- **`size`/`type`/`iops`/`throughput` update in place** — everything else replaces the volume.
- **Attachments are keyed by `device:instance`** — stable across list reorders; each attachment is ForceNew per (volume, instance, device). Multi-attach (io1/io2 only) conventionally reuses the same device name per instance.
- **`final_snapshot` is config-only at AWS** — never read back, so imports do not round-trip it (declared in the import catalog).
- **`create_time` is empty on the copy arm** — the provider does not expose it there; the output documents that honestly.

Outputs mirror the Pulumi module key-for-key: `volume_id` (import ID), `volume_arn`, `availability_zone`, `size_gb`, `create_time`.
