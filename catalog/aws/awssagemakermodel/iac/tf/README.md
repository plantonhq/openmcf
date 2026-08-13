# AwsSagemakerModel — Terraform/OpenTofu module

Deploys an Amazon SageMaker AI model (`aws_sagemaker_model`) — the
immutable serving definition endpoints deploy.

Module facts worth knowing before editing:

- **Every argument is create-time only** — the provider's update is
  tags-only, so any spec change replaces the model. That is AWS's own
  contract (roll a new model, repoint the endpoint); the module adds
  nothing on top.
- **`primary_container` and `container` share one schema upstream** —
  the spec's exactly-one rule decides which form renders; the pipeline
  renders as repeated `container` blocks in spec order.
- **The `s3_data_source` wrapper is single-valued upstream** (the
  expander reads index 0 only) — the spec flattens it to one message
  and the module renders one block.
- **EULA acceptance renders only when true** — `model_access_config`
  appears exactly when `accept_eula` is set, on both the main and
  additional data sources.
- **Container-level pairing rules live in the spec**
  (image-or-package, one artifact form, cache-requires-MultiModel) —
  the module renders what validation already admitted.

Outputs mirror the Pulumi module key-for-key: `model_name`,
`model_arn`.
