---
title: "SageMaker Image"
description: "SageMaker Image deployment documentation"
icon: "package"
order: 100
componentName: "awssagemakerimage"
---

# AWS SageMaker Image

Your container images as first-class Studio citizens — a named registry
entry that makes custom kernels and training environments selectable in
SageMaker Studio and notebooks, with folded versions each pointing at a
concrete ECR image, numbered by AWS and aliased by you.

## What Gets Created

- An image: the named registry entry with the ECR-pull role and
  in-place Studio display metadata (`display_name`, `description`).
- Versions: one per ECR image, AWS-numbered sequentially, with movable
  aliases and compatibility metadata (`job_type`, `ml_framework`,
  `processor`, `programming_lang`, `vendor_guidance`).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker control-plane permissions
  (`sagemaker:CreateImage`, `sagemaker:CreateImageVersion` and their
  siblings).

### AWS Account

- An IAM role trusting `sagemaker.amazonaws.com` with pull access to
  the version images (reference an `AwsIamRole` or pass the ARN).
- For versions: the container images pushed to ECR in the SAME account
  and region as the image.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region and the
ECR-pull role, add display metadata and versions, and deploy.

### CLI

```bash
planton apply -f image.yaml
```

## After Deploy

- `image_name` / `image_arn` identify the image; `version_numbers`
  maps each `versions` entry (by position) to its AWS-assigned number.
- Attach the image to a SageMaker domain or space to make it
  selectable in Studio; jobs reference versions by number or alias.
- Append new versions at the end of `versions` and never reorder —
  entries are keyed by position, and a changed `base_image` replaces
  the version under a new number.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
