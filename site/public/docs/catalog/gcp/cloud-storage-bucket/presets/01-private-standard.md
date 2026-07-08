---
title: "Private Standard Bucket"
description: "The default posture for application data: IAM-only access control, public access impossible, versioned objects with bounded history, and an additive grant for the workload's service account."
type: "preset"
rank: "01"
presetSlug: "01-private-standard"
componentSlug: "cloud-storage-bucket"
componentTitle: "Cloud Storage Bucket"
provider: "gcp"
icon: "package"
order: 1
---

# Private Standard Bucket

The default posture for application data: IAM-only access control, public
access impossible, versioned objects with bounded history, and an additive
grant for the workload's service account.

## What this preset creates

A STANDARD-class bucket with uniform bucket-level access (no legacy
object ACLs), `publicAccessPrevention: enforced` (no IAM grant can ever
expose it), and object versioning whose growth is capped by a lifecycle
rule keeping the 3 newest noncurrent versions.

## Prerequisites

- A `GcpServiceAccount` named `app-runtime` (the workload identity that
  reads/writes objects). Replace the reference with your own service
  account, or drop the grant and manage access elsewhere.

## Composing storage

Every consumer references the bucket's `bucket_id` output: a
`GcpBackendBucket` origin, a Cloud Function's source bucket, Dataproc
staging, or a Pub/Sub Cloud Storage sink.

## Remix ideas

- Add `kmsKeyName` (a `GcpKmsKey` reference) for CMEK-at-rest.
- Add a `softDeletePolicy` with `retentionDurationSeconds: 0` for
  high-churn scratch data where the default 7-day recovery tail is pure
  cost — or lengthen it for precious data.
- Keep `forceDestroy` unset (false) so a teardown can never erase data
  silently; set it true only for ephemeral environments.
