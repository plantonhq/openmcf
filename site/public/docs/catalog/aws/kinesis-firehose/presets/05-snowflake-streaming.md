---
title: "Snowflake Streaming Preset"
description: "Direct streaming into a Snowflake table via Snowpipe Streaming — no intermediate S3 staging, no external Snowpipe configuration, near-real-time inserts. Credentials are sourced from AWS Secrets..."
type: "preset"
rank: "05"
presetSlug: "05-snowflake-streaming"
componentSlug: "kinesis-firehose"
componentTitle: "Kinesis Firehose"
provider: "aws"
icon: "package"
order: 5
---

# Snowflake Streaming Preset

Direct streaming into a Snowflake table via Snowpipe Streaming — no intermediate S3 staging, no external Snowpipe configuration, near-real-time inserts. Credentials are sourced from AWS Secrets Manager.

## When to Use

- **Near-real-time Snowflake ingestion** — Events queryable in Snowflake seconds after they are produced
- **Zero-ops Kafka/Kinesis-to-Snowflake** — Pair with a Kinesis or MSK source to stream a topic into a table with no consumer code
- **Replacing batch COPY pipelines** — Retire staged-file loaders and their scheduling in favor of managed streaming

## Key Configuration

- **Snowpipe Streaming defaults** — 0s interval / 1 MiB buffering for near-real-time delivery
- **Secrets Manager credentials** — the key pair lives in Secrets Manager (`{"user": ..., "private_key": ..., "key_passphrase": ...}`); the manifest and IaC state never carry it
- **Dedicated Snowflake role** — `FIREHOSE_INGEST` with insert-only privileges (least privilege)
- **JSON_MAPPING loading** — top-level JSON keys map to same-named columns; switch to `VARIANT_CONTENT_MAPPING` to land whole records in one VARIANT column
- **S3 backup of failed data** — undeliverable records land in the backup bucket for replay

## Prerequisites

| Resource | Description |
|----------|-------------|
| **Snowflake account** | With the target database/schema/table created and a user configured for key-pair authentication. |
| **Secrets Manager secret** | Holding the key-pair credential in the shape above. |
| **S3 bucket** | Backup bucket for failed records. |
| **IAM roles** | (1) Delivery role with `secretsmanager:GetSecretValue` on the secret; (2) backup role with S3 write on the bucket. |

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `myaccount` | Your Snowflake account identifier |
| `ANALYTICS` / `PUBLIC` / `EVENTS` | Target database / schema / table |
| `FIREHOSE_INGEST` | Snowflake ingestion role |
| `snowflake-firehose-keypair-AbCdEf` | Secrets Manager secret name suffix |
| `my-firehose-backup-bucket` | S3 backup bucket |
| `123456789012` | Your AWS account ID |

## PrivateLink

For private connectivity, set `privateLinkVpceId` to your Snowflake PrivateLink VPCE ID (`com.amazonaws.vpce.<region>.vpce-svc-<id>`) and use the privatelink account URL.
