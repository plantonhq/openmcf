# Log Archive with Lifecycle Tiering

This preset creates a private bucket purpose-built as a log/archive destination: objects tier down through cheaper storage classes as they age and are deleted after a year, with no manual housekeeping.

## When to Use

- ALB/NLB access logs, S3 server access logs, CloudFront logs, application logs
- Any write-once data whose value decays with age
- Archive destinations where cost control matters more than retrieval speed

## Key Configuration Choices

- **Fully private** — no `publicAccessBlock` block means all four guards stay on
- **Aging tiers** (`transitions`) — logs move to STANDARD_IA at 30 days (infrequent access, cheaper storage) and GLACIER at 90 days (archival, minutes-to-hours retrieval)
- **Hard expiration** (`expiration.days: 365`) — logs delete themselves after a year; adjust to your retention policy
- **Multipart cleanup** — incomplete uploads abort after 7 days, reclaiming invisible storage from failed writes
- **Prefix-scoped** (`filter.prefix: logs/`) — only the `logs/` tree is managed; drop the filter to cover the whole bucket

## Placeholders to Replace

- `<aws-region>` — the region for the bucket
- `my-log-archive` — rename to your destination bucket name
- The `logs/` prefix — match whatever prefix your producers write under

## Common Additions

- `versioningStatus: Enabled` plus a `noncurrentVersionExpiration` rule if writers may overwrite keys
- `logging` on OTHER buckets pointing here (`targetBucket` referencing this bucket) to centralize access logs
- A `policy` statement granting `logging.s3.amazonaws.com` or `logdelivery.elasticloadbalancing.amazonaws.com` write access, depending on the log producer
- `intelligentTieringConfigurations` instead of fixed-day transitions when access patterns are unpredictable

## Related Presets

- **01-private-encrypted** — the default posture for application data
- **02-public-static-website** — a deliberately public website bucket
