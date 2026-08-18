# AwsS3DirectoryBucket — Terraform/OpenTofu module

Manages one S3 directory bucket (`aws_s3_directory_bucket`) - S3 Express One Zone.

Module facts worth knowing before editing:

- **The bucket name is DERIVED, never user-assembled** — AWS mandates `{base}--{zone_id}--x-s3`; this module builds it from `metadata.name` + `spec.zone_id` so the name and the location block can never disagree. The derived name is the `bucket_name` output.
- **Everything except `force_destroy` replaces the bucket** — a directory bucket is replaced, not edited.
- **`data_redundancy` is sent only when set** — the provider derives the only-valid value (SingleAvailabilityZone / SingleLocalZone) from the location type; the spec's CEL rejects a mismatched explicit pair before any plan.
- **`force_destroy` is config-only at AWS** — never read back, so imports do not round-trip it (declared in the import catalog).

Outputs mirror the Pulumi module key-for-key: `bucket_name` (import ID), `bucket_arn`.
