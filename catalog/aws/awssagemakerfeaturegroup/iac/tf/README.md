# AwsSagemakerFeatureGroup — Terraform/OpenTofu module

Deploys an Amazon SageMaker Feature Store feature group
(`aws_sagemaker_feature_group`) with its online and/or offline store.

Module facts worth knowing before editing:

- **Almost everything is create-time only** (schema, stores, role) —
  the ONLY in-place updates are the online store's TTL and the
  throughput settings.
- **The component's name IS the feature group name** — derived from
  `metadata.name`, no override.
- **The provider's `collection_config` has exactly one member**
  (`vector_config`) — rendered exactly when the dimension is set (the
  spec pairs it with the Vector collection type).
- **Provisioned capacity units are sent only in Provisioned mode** —
  the throughput expander mirrors the provider's create behavior
  (spec-validated pairing).
- **Collection types require InMemory online storage** —
  server-enforced by AWS, not rendered by the module.
- **Destroy leaves the offline store's S3 objects in place** — AWS
  design, not a module choice.

Outputs mirror the Pulumi module key-for-key: `feature_group_name`,
`feature_group_arn`.
