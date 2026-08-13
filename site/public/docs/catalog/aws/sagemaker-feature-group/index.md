---
title: "SageMaker Feature Group"
description: "SageMaker Feature Group deployment documentation"
icon: "package"
order: 100
componentName: "awssagemakerfeaturegroup"
---

# AWS SageMaker Feature Group

ML features as managed infrastructure — a declared feature schema over
an online store for low-latency serving and/or an offline S3/Glue store
for training, with the record identifier and event time anchoring every
record and TTL-based expiry on the serving side.

## What Gets Created

- A feature group named after `metadata.name`, with a schema of 1–2500
  features (scalars, or List / Set / Vector collections — vectors carry
  a dimension).
- An online store (Standard or InMemory storage, optional KMS, optional
  record TTL) and/or an offline store (S3 location, auto-created Glue
  table in Glue or Iceberg format, optional KMS) — at least one.
- On-demand (default) or provisioned read/write throughput.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker Feature Store control-plane
  permissions (`sagemaker:CreateFeatureGroup` and its siblings).

### AWS Account

- An IAM role trusting `sagemaker.amazonaws.com` that can write the
  offline store's bucket (`role_arn` — referenceable from an
  AwsIamRole).
- For an offline store: the S3 bucket behind `s3_uri`.

## Deploy

### Console

Create the resource from the AWS catalog, declare the feature schema
and the record-identifier / event-time features, pick the stores, and
deploy.

### CLI

```bash
planton apply -f feature-group.yaml
```

## After Deploy

- `feature_group_name` / `feature_group_arn` identify the group;
  ingestion (PutRecord) and serving (GetRecord) key on the name.
- Only the online store's TTL and the throughput settings update in
  place — everything else replaces the group.
- Deleting the group does NOT delete the offline store's S3 objects
  (AWS design) — clean the bucket yourself if the data must go.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
