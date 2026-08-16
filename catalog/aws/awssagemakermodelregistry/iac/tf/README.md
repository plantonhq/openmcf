# AwsSagemakerModelRegistry — Terraform/OpenTofu module

Deploys an Amazon SageMaker model package group
(`aws_sagemaker_model_package_group`) with its folded resource policy
(`aws_sagemaker_model_package_group_policy`).

Module facts worth knowing before editing:

- **The group is immutable except tags** — even the description is
  ForceNew upstream; a description edit replaces the group.
- **The component's name IS the group name** — derived from
  `metadata.name`, no override.
- **The policy is an idempotent upsert** (`PutModelPackageGroupPolicy`)
  that updates in place; the resource is rendered exactly when the spec
  carries `resource_policy`, and removing it deletes the policy —
  leaving the group open to its own account only.
- **Model package VERSIONS register into the group imperatively**
  (training pipelines, SDK) — never through this module.

Outputs mirror the Pulumi module key-for-key:
`model_package_group_name`, `model_package_group_arn`.
