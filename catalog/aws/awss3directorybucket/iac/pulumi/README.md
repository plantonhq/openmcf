# AwsS3DirectoryBucket — Pulumi module (Go)

Manages one S3 directory bucket (`s3.DirectoryBucket`) - S3 Express One Zone.

Module facts worth knowing before editing:

- **The bucket name is DERIVED, never user-assembled** — AWS mandates `{base}--{zone_id}--x-s3`; `locals.BucketName` builds it from `metadata.name` + `spec.zone_id` so the name and the location block can never disagree. The derived name is the `bucket_name` output.
- **Everything except `ForceDestroy` replaces the bucket** — a directory bucket is replaced, not edited.
- **`DataRedundancy` is sent only when set** — the provider derives the only-valid value from the location type; the spec's CEL rejects a mismatched explicit pair before any preview.
- **`ForceDestroy` is config-only at AWS** — never read back, so imports do not round-trip it (declared in the import catalog).

Outputs mirror the Terraform module key-for-key: `bucket_name` (import ID), `bucket_arn`.
