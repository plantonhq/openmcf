<p align="center">
  <img src="logo.svg" alt="AWS SageMaker Feature Group" width="80"/>
</p>

# AWS SageMaker Feature Group

Create and manage [Amazon SageMaker Feature Store feature groups](https://docs.aws.amazon.com/sagemaker/latest/dg/feature-store.html)
— the schema'd stores ML features are written to and served from, with
an online store for low-latency serving and/or an offline S3/Glue store
for training datasets and point-in-time queries.

## What Gets Created

- **A feature group** whose AWS name derives from `metadata.name`, with
  a declared schema of 1–2500 `feature_definitions` (Integral /
  Fractional / String; List / Set / Vector collections — a Vector
  pairs exactly with its `vector_dimension`).
- The `record_identifier_feature_name` and `event_time_feature_name`
  anchoring every record — both must be members of the schema.
- Optional **online store**: Standard or InMemory storage, KMS
  encryption, and a record TTL (`ExpiresAt = EventTime + ttl`).
- Optional **offline store**: an S3 location with an auto-created Glue
  Data Catalog table (Glue or Iceberg format), KMS encryption, or your
  own named catalog entry.
- Optional **throughput** settings: on-demand (default) or provisioned
  capacity.

## Almost Everything Is Create-Time Only

The schema, stores, and role are AWS's create-time contract — the ONLY
in-place updates are the online store's TTL and the throughput
settings; changing anything else replaces the group. And by AWS design,
the offline store's S3 objects survive group deletion.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
