# Private Versioned Bucket

This preset creates a private Spaces bucket with versioning on and a lifecycle rule that controls versioning's storage cost: noncurrent object versions are deleted after 90 days and abandoned multipart uploads after 7.

## When to Use

- Application data, backups, or artifacts that must never be publicly readable
- Data where accidental overwrites or deletions must be recoverable (versioning keeps every version)
- Any bucket that will hold real data long-term -- the lifecycle rule keeps version sprawl from becoming a silent storage bill

## Key Configuration Choices

- **Private ACL** (`accessControl: PRIVATE`) -- objects are reachable only with Spaces credentials. Finer-grained access (per-prefix, per-IP) is the `policy` field's job.
- **Versioning** (`versioningEnabled: true`) -- every overwrite and delete keeps the previous version. Note versioning can never be removed once enabled, only suspended (flipping this back to false suspends it).
- **Noncurrent-version expiration** (`noncurrentVersionExpiration.days: 90`) -- old versions are kept 90 days, then deleted permanently. Without this, a frequently rewritten object accumulates versions forever.
- **Multipart-upload cleanup** (`abortIncompleteMultipartUploadDays: 7`) -- abandoned large-file uploads are aborted and their invisible parts removed.

## Placeholders to Replace

- `metadata.name` / `bucketName` -- your bucket's name. It must be globally unique within the region (across ALL DigitalOcean customers) and DNS-compatible.
- `region` -- a Spaces-capable region slug (ams3, atl1, blr1, fra1, lon1, nyc3, sfo2, sfo3, sgp1, syd1, tor1).
