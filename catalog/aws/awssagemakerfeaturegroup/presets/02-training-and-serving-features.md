# Training and Serving Features

This preset is a dual-store feature group: every write serves online at
low latency AND lands in S3 under an auto-created Glue table for
training datasets and point-in-time queries.

## When to Use

- Features that both serve real-time inference and feed model training
- Avoiding train/serve skew by writing one record to both stores

## What You Get

- An online store with a 30-day record TTL for serving
- An offline store on `s3://my-feature-store/customers/` with a Glue
  Data Catalog table created automatically, partitioned by
  `event_time`

## Customize

- Point `s3Uri` at your bucket, and make sure `roleArn` can write it
- Switch `tableFormat` to `Iceberg` (create-time) for faster queries
  and compaction
- Remember: deleting the group leaves the S3 objects in place by AWS
  design — the bucket's lifecycle is yours to manage
