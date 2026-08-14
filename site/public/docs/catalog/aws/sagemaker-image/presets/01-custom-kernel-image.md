---
title: "Custom Kernel Image"
description: "This preset creates the registry entry for a team's custom Studio kernels — the named shelf with display metadata, ready before any container image exists. Versions land later, appended as images are..."
type: "preset"
rank: "01"
presetSlug: "01-custom-kernel-image"
componentSlug: "sagemaker-image"
componentTitle: "SageMaker Image"
provider: "aws"
icon: "package"
order: 1
---

# Custom Kernel Image

This preset creates the registry entry for a team's custom Studio
kernels — the named shelf with display metadata, ready before any
container image exists. Versions land later, appended as images are
pushed to ECR.

## When to Use

- Standing up the registry entry ahead of the first image build
- CI pipelines that append versions as builds complete

## What You Get

- A named SageMaker image with Studio display metadata
  (`displayName`, `description`) that updates in place
- The ECR-pull role wired, so the first appended version registers
  without further IAM work

## Customize

- Append entries to `versions` as images are pushed — always at the
  end, never reordered; entries are keyed by position
- Point `roleArn` at a composed `AwsIamRole` with `valueFrom` instead
  of a literal ARN
- Set `jobType`, `mlFramework`, and `processor` on each version so
  Studio surfaces the right compatibility hints
