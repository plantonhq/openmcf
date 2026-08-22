# AwsEbsVolume — Pulumi module (Go)

Manages one EBS volume as a create-XOR-copy union (`ebs.Volume` XOR `ebs.VolumeCopy`) with attachments (`ec2.VolumeAttachment`) in-line.

Module facts worth knowing before editing:

- **Exactly one arm exists** — the spec's union CEL guarantees `copy_from` XOR the create-arm fields; both arms feed the same exported outputs.
- **A copy inherits placement and encryption** — `ebs.VolumeCopy` has no zone or encryption arguments; only size/type/iops/throughput may be overridden.
- **`Size`/`Type`/`Iops`/`Throughput` update in place** — everything else replaces the volume.
- **Attachments are named `attachment-{device}-{instance}`** — one resource per pair, ForceNew at the provider; multi-attach (io1/io2 only) conventionally reuses the same device name per instance.
- **`FinalSnapshot` is config-only at AWS** — never read back, so imports do not round-trip it (declared in the import catalog).
- **`create_time` is empty on the copy arm** — the provider does not expose it there.

Outputs mirror the Terraform module key-for-key: `volume_id` (import ID), `volume_arn`, `availability_zone`, `size_gb`, `create_time`.
