# AwsEbsVolume

One standalone EBS volume as a create-XOR-copy union — a fresh volume in a chosen availability zone (empty or restored from a snapshot) or a clone of an existing volume — with instance attachments managed in-line.

## Highlights

- **The persistent-disk building block**: databases on EC2, stateful services, scratch space that outlives instances — provisioned independently of any instance and attached where needed.
- **Two arms, one downstream surface**: create fresh (zone + size/snapshot + encryption) or `copy_from` an existing volume; either way charts reference the same `volume_id`/`volume_arn` outputs.
- **AWS's real constraint graph as CELs**: iops needs an iops-capable type, io1/io2 REQUIRE iops, throughput is gp3-only, multi-attach is io-family-only, more than one attachment requires multi-attach — every wall AWS enforces at apply is caught at validation.
- **Attachments live with the volume**: each entry is one instance + device pair with its detach posture (`force_detach`, `skip_destroy`, `stop_instance_before_detaching`) — removing an entry detaches cleanly.

## Both Engines

Both modules render the same arm selection and export the same outputs: `volume_id` (import ID), `volume_arn`, `availability_zone` (notably useful on copies, which inherit the source's), `size_gb`, `create_time`.

## Chart Wiring

`instance_id` references an AwsEc2Instance output; `snapshot_id` references an AwsEbsSnapshot; `kms_key_id` references an AwsKmsKey. The `volume_id` output is what AwsEbsSnapshot's volume arm and other volumes' `copy_from` reference back.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
