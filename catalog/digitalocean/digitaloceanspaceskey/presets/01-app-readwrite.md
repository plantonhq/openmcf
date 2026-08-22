# App Read-Write Key

This preset mints the credential one application uses for its own bucket: readwrite on exactly that bucket, nothing else. The workload can upload and serve its files and cannot touch any other data in the account.

## When to Use

- The standard per-workload credential: one app, one bucket, one key
- CI pipelines that publish artifacts to a single bucket

## Key Configuration Choices

- **The bucket by reference** -- the grant follows the bucket resource; a chart deploys the bucket and its key together.
- **readwrite, not fullaccess** -- the key unlocks one bucket; a leak exposes that bucket, not the account.
- **Capture the secret at creation** -- `secret_key` is returned exactly once and can never be fetched again.

## What You Get

A free, scoped S3-style credential for the bucket's regional Spaces endpoint -- least privilege by construction.
