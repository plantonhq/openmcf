---
title: "SageMaker Model"
description: "SageMaker Model deployment documentation"
icon: "package"
order: 100
componentName: "awssagemakermodel"
---

# AWS SageMaker Model

The serving definition behind every SageMaker deployment — a container
image (or model package) paired with S3 model artifacts, an execution
role, and optional VPC confinement, declared once and deployed by
endpoints, batch transform jobs, and inference components.

## What Gets Created

- A model as a single container (`primary_container`) or an inference
  pipeline of 2–15 containers (`containers`) run in `Serial` or
  addressed by hostname in `Direct` mode — exactly one form.
- Per-container artifact wiring: compressed `model_data_url`, or
  uncompressed `model_data_source` with EULA acceptance for gated
  models, plus adapter channels, MultiModel serving, and
  private-registry image configuration.
- Optional VPC attachment and full network isolation for private
  serving.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker control-plane permissions
  (`sagemaker:CreateModel` and its siblings).

### AWS Account

- An IAM role trusting `sagemaker.amazonaws.com` with ECR pull and S3
  read access to the image and artifacts (`execution_role_arn`).
- Model artifacts staged in S3 (a `.tar.gz` for the compressed form, a
  prefix for uncompressed or MultiModel serving).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region and execution
role, define the container (or pipeline), and deploy.

### CLI

```bash
planton apply -f model.yaml
```

## After Deploy

- `model_name` / `model_arn` identify the model; endpoint variants
  reference it by `model_name`.
- The model is immutable — any change replaces it. Roll a new model
  and repoint the endpoint (AWS's own contract).
- The model itself is free; billing starts when an endpoint or batch
  job deploys it.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
