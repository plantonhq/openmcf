<p align="center">
  <img src="logo.svg" alt="AWS SageMaker Model Registry" width="80"/>
</p>

# AWS SageMaker Model Registry

Create and manage [Amazon SageMaker model package groups](https://docs.aws.amazon.com/sagemaker/latest/dg/model-registry.html)
— the model registry's unit of organization, holding the versioned
model packages a team registers, approves, and deploys from.

## What Gets Created

- **A model package group** whose AWS name derives from
  `metadata.name`, with an optional description (max 1024 characters —
  write it once, well: changing it replaces the group).
- Optional: a **folded resource policy** for cross-account model
  sharing (grant other accounts `sagemaker:DescribeModelPackage` /
  `sagemaker:CreateModelPackage` on the group) — the one arm that
  updates in place.

## The Group Is a Named Shell

Model package VERSIONS register into the group imperatively — from
training pipelines and SDK calls, never declaratively. The group itself
is immutable except tags; even a description edit replaces it
(provider-enforced).

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
