# AWS S3 Directory Bucket

S3 at memory-adjacent speed: a bucket that lives in one availability zone next to your compute, serving objects in single-digit milliseconds — the storage tier for ML training data, feature stores, and hot intermediate results that regular S3 is too slow for.

## What Gets Managed

- The directory bucket: its zone (availability zone or Local Zone, by ID), redundancy class, and force-destroy posture. The full AWS name (`{base}--{zone_id}--x-s3`) is derived — you name the base, the modules handle AWS's format.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with S3 Express permissions.

### AWS Prerequisites

- An Express-supported zone in your region (not every AZ supports One Zone — check the region's zone list).

## After You Deploy

- Address objects at the `bucket_name` output through the S3 Express session API (SDKs handle the session under `CreateSession` permissions).
- Co-locate the readers: the latency win exists only for compute in the SAME zone.

## Common Changes

- There are almost none by design: everything except `force_destroy` replaces the bucket. Moving zones or renaming means a new bucket and a data copy.
