# AwsSagemakerModel — Pulumi module (Go)

Deploys an Amazon SageMaker AI model (`sagemaker.Model`) — the
immutable serving definition endpoints deploy.

Module facts worth knowing before editing:

- **Every argument is create-time only** — the provider's update is
  tags-only, so any spec change replaces the model. That is AWS's own
  contract (roll a new model, repoint the endpoint); the module adds
  nothing on top.
- **Primary and pipeline containers share one schema upstream but two
  bridged types here** (`ModelPrimaryContainerArgs` vs
  `ModelContainerArgs`) — the shared container message maps onto each
  through parallel builders; keep them in lockstep.
- **The `s3_data_source` wrapper is single-valued upstream** (the
  expander reads index 0 only) — rendered as a one-element array.
- **EULA acceptance renders only when true** — `ModelAccessConfig`
  appears exactly when `accept_eula` is set, on both the main and
  additional data sources.
- **Container-level pairing rules live in the spec**
  (image-or-package, one artifact form, cache-requires-MultiModel) —
  the module renders what validation already admitted.

Outputs mirror the Terraform module key-for-key: `model_name`,
`model_arn`.
