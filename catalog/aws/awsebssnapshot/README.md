# AwsEbsSnapshot

One EBS snapshot as a three-way source union — snapshot a live volume, copy an existing snapshot (same- or cross-region, optionally re-encrypting), or import a disk image through VM Import/Export — with archive tiering, fast snapshot restore, and cross-account share grants managed in-line.

## Highlights

- **Three ways a snapshot is born, one kind**: `volume_id` XOR `copy_from` XOR `import_from`, mirroring the provider's three resources; the exactly-one CEL keeps the arms honest.
- **Copying is the ONLY re-encryption path**: snapshots never re-encrypt in place — the copy arm carries `encrypted` + `kms_key_id` for exactly that reason, taught on the fields.
- **The cost levers are explicit**: `storage_tier: archive` is ~75% cheaper with a restore step; `fast_restore_availability_zones` bills per zone-hour while enabled — each a deliberate dial, never a default.
- **Import honesty pre-declared**: only the volume arm is importable at the provider (the copy, import, and share-grant resources ship no importer) — declared in the import catalog before any adopter finds out the hard way.

## Both Engines

Both modules render the same arm selection with the FSR and share satellites and export the same outputs: `snapshot_id` (import ID, volume arm), `snapshot_arn`, `owner_id`, `volume_size_gb`.

## Chart Wiring

`volume_id` references an AwsEbsVolume; `source_snapshot_id` references another AwsEbsSnapshot; `kms_key_id` references an AwsKmsKey. The `snapshot_id` output is what AwsEbsVolume restores and copies reference back.
