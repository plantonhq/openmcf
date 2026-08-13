---
title: "SageMaker Model Registry"
description: "SageMaker Model Registry deployment documentation"
icon: "package"
order: 100
componentName: "awssagemakermodelregistry"
---

# AWS SageMaker Model Registry

A model registry group as managed infrastructure — the named shell your
training pipelines register versioned model packages into, with an
optional cross-account resource policy so other accounts can discover
and register models on the group.

## What Gets Created

- A model package group named after `metadata.name`, with an optional
  description.
- Optionally, a resource policy attached to the group — structured JSON
  granting other accounts access (e.g.
  `sagemaker:DescribeModelPackage`, `sagemaker:CreateModelPackage`).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker control-plane permissions
  (`sagemaker:CreateModelPackageGroup` and its siblings).

### AWS Account

- Nothing beyond the connection — the group is free-standing; creating
  it costs nothing.
- For cross-account sharing: the account IDs you will name as
  principals in `resource_policy`.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, write the
description once (it replaces the group when changed), and deploy.

### CLI

```bash
planton apply -f model-registry.yaml
```

## After Deploy

- `model_package_group_name` / `model_package_group_arn` identify the
  group; training pipelines register model package versions into it by
  name — imperatively, never through this resource.
- The resource policy is the one arm that updates in place — iterate on
  sharing freely; removing it from the spec deletes the policy, leaving
  the group open to its own account only.
- Everything else, tags aside, replaces the group.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
