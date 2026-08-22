# DigitalOcean Spaces Bucket -- Operational Guide

Judgment calls that matter when you run DigitalOcean Spaces buckets.

## Spaces is a second credential plane

The DigitalOcean API token cannot read or write Spaces. Both provisioners and the verifier speak the S3-compatible Spaces API with `SPACES_ACCESS_KEY_ID` / `SPACES_SECRET_ACCESS_KEY` (the provider's own env names). A deploy that only has `DIGITALOCEAN_TOKEN` will create nothing and the verifier will refuse to run. Mint a Spaces key in the control panel (API -> Spaces Keys) and put both values in the environment, or on the DigitalOcean provider connection.

## Names collide globally

A Spaces bucket name must be unique within its region across ALL DigitalOcean customers, not just your account. Prefix with an org or project slug. The e2e scenarios append a run-id token for the same reason.

## Region is optional until a satellite is set

Unset `region` lets the provider apply its own default (`nyc3`). That is the smallest real bucket. The moment you set `corsRules`, `policy`, or `logging`, the spec requires an explicit region: those are separate provider resources whose region argument is required, and they cannot inherit the bucket's default. Changing `region` later replaces the bucket.

Spaces is not in every DigitalOcean region. The spec accepts only the Spaces-capable slugs: `ams3`, `atl1`, `blr1`, `fra1`, `lon1`, `nyc3`, `sfo2`, `sfo3`, `sgp1`, `syd1`, `tor1`.

## Versioning cannot be undone

`versioningEnabled: true` keeps every overwrite and delete as a recoverable version. Flipping it back to false only suspends versioning — existing versions stay, and new writes no longer create them. Pair versioning with a `lifecycleRules` entry that expires noncurrent versions, or the storage bill grows silently.

## Lifecycle rules need an action

A rule with only `id`/`prefix`/`enabled` does nothing. Set at least one of: `expiration` (exactly one of `date`, `days`, or `expiredObjectDeleteMarker`), `noncurrentVersionExpiration.days`, or `abortIncompleteMultipartUploadDays`. The spec rejects an expiration that sets more than one trigger — the provider would send only one (date wins) and the others would silently disappear.

## CORS rides the standalone resource

The bucket's inline `cors_rule` argument is deprecated at the pinned provider and never reads back (no drift detection). This kind writes CORS through `digitalocean_spaces_bucket_cors_configuration`, which does call `GetBucketCors`. Do not also manage CORS out-of-band on the same bucket — last writer wins.

## `forceDestroy` empties the bucket on destroy

When true, destroy deletes every object AND every object version before removing the bucket. That is the only way a non-empty bucket will tear down. Leave it false for production data; set it true for ephemeral and test buckets so teardown cannot stall.

## Access logging needs a second bucket

`logging.targetBucket` is a literal name or a `DigitalOceanBucket` reference (resolved from `status.outputs.bucket_id` — a Spaces bucket's id IS its name). The two buckets must be in the same region. Logging a bucket to itself works but compounds: every log read writes more logs.

## Importing an existing bucket

Import uses the `<region>,<name>` composite (the `region` and `bucket_id` outputs). Expect `acl` and `forceDestroy` to stay at their configured values after import: the API never reports the canned ACL, and the importer hardcodes `forceDestroy` to false. Policy JSON is normalized (whitespace, key order) on read — a formatting-only diff is suppressed by the provider.

## What is deliberately NOT here

Uploading object contents (data-plane work, not infrastructure). Spaces access keys (their own kind — they have an independent create/rotate lifecycle). CDN endpoints (a separate resource in front of a public bucket). The deprecated inline `cors_rule` on the bucket itself.
