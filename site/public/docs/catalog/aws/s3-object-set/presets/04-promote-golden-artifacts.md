---
title: "Promote Golden Artifacts"
description: "This preset copies released build artifacts from a golden artifacts bucket into an environment's bucket, server-side — the bytes move inside S3 and never through the deploy host, so artifact size is..."
type: "preset"
rank: "04"
presetSlug: "04-promote-golden-artifacts"
componentSlug: "s3-object-set"
componentTitle: "S3 Object Set"
provider: "aws"
icon: "package"
order: 4
---

# Promote Golden Artifacts

This preset copies released build artifacts from a golden artifacts bucket into an environment's bucket, server-side — the bytes move inside S3 and never through the deploy host, so artifact size is unconstrained. The first object shows the promotion guard: an ETag precondition that fails the deploy if the source artifact changed since the release was cut. The second shows a copy that stamps its own headers and metadata onto the destination (`replaceMetadata`) instead of preserving the source's.

## When to Use

- Promoting build artifacts between environment buckets (staging → production) as a tracked infrastructure change
- Seeding a new environment's bucket from a golden source bucket
- Re-homing large objects with new placement or metadata without re-uploading them

## Key Configuration Choices

- **Server-side copy** (`copyFrom`) -- The source object is named by bucket + key; the source bucket can also be an `AwsS3Bucket` component reference or an access-point ARN
- **Promotion guard** (`copyIfMatch`) -- The copy only proceeds if the source's ETag still matches the recorded release ETag; a changed source fails the deploy instead of silently promoting different bytes
- **Metadata control** (`replaceMetadata`) -- The default preserves everything the source object carried; setting it writes this entry's own `contentType`, `cacheControl`, and `metadata` to the copy

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region of the DESTINATION bucket (e.g., `us-east-1`) | Must match the destination bucket's region |
| `<destination-bucket-name>` | Bucket receiving the promoted artifacts | AWS S3 console or `AwsS3Bucket` status outputs |
| `<golden-artifacts-bucket-name>` | Bucket holding the released artifacts | AWS S3 console or `AwsS3Bucket` status outputs |
| `<source-artifact-key>` | Key of the artifact to promote (e.g., `builds/app/1.2.3/app.zip`) | Your release pipeline's output location |
| `<source-artifact-etag>` | The source artifact's ETag at release time | `aws s3api head-object`, or the set's `object_etags` output where the artifact was created |
| `<source-manifest-key>` | Key of the release manifest to promote | Your release pipeline's output location |
