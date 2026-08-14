---
title: "Shared Model Registry"
description: "This preset is a cross-account registry group: a resource policy on the group lets another AWS account discover and list the model versions registered into it."
type: "preset"
rank: "02"
presetSlug: "02-shared-model-registry"
componentSlug: "sagemaker-model-registry"
componentTitle: "SageMaker Model Registry"
provider: "aws"
icon: "package"
order: 2
---

# Shared Model Registry

This preset is a cross-account registry group: a resource policy on the
group lets another AWS account discover and list the model versions
registered into it.

## When to Use

- A training account that publishes models a production account
  deploys
- Central registries consumed by multiple workload accounts

## What You Get

- A model package group with a folded resource policy granting the
  `210987654321` account `sagemaker:DescribeModelPackage` and
  `sagemaker:ListModelPackages`
- In-place policy updates — sharing evolves without replacing the group

## Customize

- Replace the principal with your consumer account's root ARN, or list
  several
- Add `sagemaker:CreateModelPackage` to the actions when the other
  account should register versions into the group, not just read them
- Remove `resourcePolicy` to withdraw sharing — the policy deletes,
  leaving the group open to its own account only
