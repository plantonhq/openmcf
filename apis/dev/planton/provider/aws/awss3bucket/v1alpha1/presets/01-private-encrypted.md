# Private Encrypted Bucket

This preset creates a fully private, versioned S3 bucket encrypted with a customer-managed KMS key. It is the right starting point for application data, backups, and anything holding real information.

## When to Use

- Application data, documents, and user uploads
- Backup and archive targets
- Any bucket subject to compliance requirements that mandate customer-controlled encryption keys and audit trails

## Key Configuration Choices

- **Fully private by default** — no `publicAccessBlock` block means all four public-access guards are enabled; no policy or ACL can expose the bucket accidentally
- **Versioning enabled** (`versioningStatus: Enabled`) — every overwrite and delete keeps the previous version, protecting against mistakes and ransomware-style overwrites
- **SSE-KMS with a bucket key** (`encryption`) — CloudTrail-audited key usage and customer control over the key policy; `bucketKeyEnabled: true` cuts KMS request costs by up to 99%
- **Noncurrent version pruning** (`lifecycleRules`) — deletes superseded versions after 30 days so versioning does not silently multiply storage costs, and aborts failed multipart uploads after 7 days

## Placeholders to Replace

- `<aws-region>` — the region for the bucket (e.g., `us-west-2`)
- `<kms-key-arn>` — the ARN of your KMS key; or replace the `value` arm with a `valueFrom` reference to an `AwsKmsKey` resource:

```yaml
kmsKeyId:
  valueFrom:
    kind: AwsKmsKey
    name: my-data-key
    fieldPath: status.outputs.key_arn
```

## Common Additions

- Drop the `encryption` block entirely to fall back to free SSE-S3 (AES256) — AWS encrypts every object by default
- Add storage-class `transitions` to the lifecycle rule to tier aging data into STANDARD_IA or GLACIER classes
- Add a `replication` block (with an IAM role and a versioned destination bucket) for cross-region disaster recovery
- Add `notification` with `eventbridge: true` to route object events into EventBridge rules

## Related Presets

- **02-public-static-website** — a deliberately public bucket serving a static site
- **03-log-archive-lifecycle** — a log destination with aggressive storage tiering and expiration
