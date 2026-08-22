# AwsS3VectorBucket — Pulumi module (Go)

Manages one S3 vector bucket with its indexes and policy (`s3.VectorsVectorBucket` + `s3.VectorsIndex` + `s3.VectorsVectorBucketPolicy`) - three provider resources under one declarative owner.

Module facts worth knowing before editing:

- **Every index property is create-only** — an index is replaced, not edited; dimension must match the embedding model BEFORE the first vector lands.
- **Indexes are named `index-{name}`** — stable across list reorders; a rename replaces the index (and its vectors).
- **`DataType` is module-pinned to float32** — the provider's enum holds exactly that one value (recorded parity exclusion).
- **Bucket encryption is fixed for life**; the bridge models `encryption_configuration` as an ARRAY (`EncryptionConfigurations`) though the provider caps it at one — this module always sends exactly one element.
- **The policy is JSON-normalized by AWS**; **`ForceDestroy` is config-only** and never round-trips on import (declared in the import catalog).

Outputs mirror the Terraform module key-for-key: `vector_bucket_arn` (import ID), `index_arns` (keyed by index name).
