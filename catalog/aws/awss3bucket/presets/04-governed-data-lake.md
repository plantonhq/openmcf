# Governed Data Lake

This preset creates a private, versioned data-lake bucket with the full governance toolkit switched on: scheduled inventory reports, storage-class analytics, request metrics, and S3 Metadata tables — so the bucket audits and optimizes itself instead of relying on ad-hoc listing scripts.

## When to Use

- Data-lake buckets (raw/curated zones) large enough that `ListObjects` is impractical for audits
- Buckets with compliance duties: prove encryption/replication posture per object from the weekly inventory
- Cost-optimization work: the analytics export tells you when to add lifecycle transitions, the day-zero INTELLIGENT_TIERING rule handles unpredictable access on the raw zone

## Key Configuration Choices

- **Day-zero tiering** (`transitions.days: 0`) — raw-zone objects enter INTELLIGENT_TIERING at upload, letting AWS optimize storage cost per access pattern from day one
- **Weekly inventory to itself** (`inventoryConfigurations`) — Parquet manifests of every object version, with encryption/replication/storage-class columns, delivered under `governance/inventory/`; the self-ARN works because bucket ARNs are name-derived
- **Storage-class analysis with export** (`analyticsConfigurations`) — observes the curated zone's access pattern and exports daily CSV findings; read them before adding fixed-day transitions
- **Request metrics** (`metricsConfigurations`) — CloudWatch request-level visibility for the curated zone (standard CloudWatch per-metric cost)
- **S3 Metadata tables** (`metadataConfiguration`) — queryable Iceberg journal + live inventory tables; find objects with Athena instead of listing. Journal records expire after 90 days to bound table cost
- **SSE-C blocked** (`blockedEncryptionTypes`) — states AWS's own March-2026 posture explicitly so it survives account-level changes

## Placeholders to Replace

- `<aws-region>` — the region for the bucket
- `my-data-lake` — rename to your bucket name, INCLUDING inside the two self-referencing `bucketArn` literals (`arn:aws:s3:::my-data-lake`)
- The `raw/` and `curated/` prefixes — match your lake's zone layout

## Notes

- Inventory and analytics DELIVERY to the bucket requires a bucket policy allowing `s3.amazonaws.com` to `s3:PutObject` (scoped to the delivery prefixes); AWS documents the exact statement. The configurations themselves apply without it — delivery just stays pending until the policy lands
- The metadata tables live in an AWS-managed table bucket and bill per S3 Tables pricing (near zero for small buckets, real for hot ones)

## Related Presets

- **01-private-encrypted** — the default posture for application data
- **03-log-archive-lifecycle** — fixed-day tiering when access patterns are known
