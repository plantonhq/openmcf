# AwsS3DirectoryBucket

One S3 directory bucket (S3 Express One Zone): single-digit-millisecond object storage that lives in exactly ONE availability zone (or Local Zone), co-located with the compute that hammers it.

## Highlights

- **The derived name that can never disagree**: AWS mandates `{base}--{zone_id}--x-s3` — both modules build it from `metadata.name` + `zone_id`, export it as `bucket_name`, and the spec never carries a hand-assembled copy.
- **The zone IS the product**: you pick the zone by ID (`use1-az4`, never the letter name), the redundancy class derives from the zone type, and everything except `force_destroy` replaces the bucket — an Express bucket is placed, not edited.
- **Honest single-AZ semantics**: the data has no cross-AZ redundancy by design — the speed trade taught on the fields, not discovered in an incident.

## Both Engines

Both modules derive the identical bucket name and export the same outputs: `bucket_name` (import ID), `bucket_arn`.

## Chart Wiring

Co-locate by wiring the same zone into the compute (AwsEc2Instance via its subnet's AZ, EKS node groups) and this bucket. Standalone otherwise — Express buckets authenticate via the S3 Express session API against the bucket name.
