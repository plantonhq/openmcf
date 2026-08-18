# AwsLambdaLayer

A Lambda layer version — a shared code archive (libraries, custom runtimes, config files) that functions attach by ARN — published from S3, with its cross-account and organization share grants managed in-line.

## Highlights

- **Immutable by contract, taught up front**: every layer-version argument is ForceNew — a config change publishes a NEW version, and functions keep the exact version ARN they pinned, so a replacement never breaks consumers mid-run.
- **S3 is the one code source**: modules run remote from the author's machine, so the local-file arm the provider offers is deliberately out; stage the archive with your pipeline (or AwsS3ObjectSet) and reference the bucket.
- **Grants are per-version statements**: each `permissions` entry is one policy statement keyed by statement_id — account-scoped, or organization-scoped through the wildcard principal.

## Both Engines

Both modules publish the version and attach the grants identically and export the same outputs: `layer_arn` (the identity that persists across versions), `layer_version_arn` (what functions attach — the import ID), `version`, `code_sha256`, plus the `permission_revision_ids` map keyed by statement_id.

## Chart Wiring

`code.bucket` → AwsS3Bucket `bucket_id`; `layer_version_arn` → AwsLambda's `layers` list. Publish the layer and the functions that attach it in one chart — the version ARN wires by reference.
