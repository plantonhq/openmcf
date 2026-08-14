---
title: "SageMaker MLflow App"
description: "SageMaker MLflow App deployment documentation"
icon: "package"
order: 100
componentName: "awssagemakermlflowapp"
---

# AWS SageMaker MLflow App

Serverless MLflow as managed infrastructure — the MLflow 3.x successor
to the hourly-billed tracking server, with nothing to size, billed per
use and $0 when idle, storing artifacts in your S3 bucket and serving
as the default MLflow for the SageMaker domains you associate.

## What Gets Created

- A serverless MLflow 3.x app with your S3 artifact store
  (`artifact_store_uri`) accessed through the app's IAM role.
- Domain associations (`default_domain_ids`): Studio users in those
  domains track to this app automatically.
- Optional: account-default status (one app per account), automatic
  model registration into the SageMaker Model Registry, and a weekly
  maintenance window.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker control-plane permissions
  (`sagemaker:CreateMlflowApp` and its siblings).

### AWS Account

- An S3 bucket (or prefix) for the artifact store.
- An IAM role trusting `sagemaker.amazonaws.com` with read/write
  access to that bucket (reference an `AwsIamRole` or pass the ARN).
- For domain associations: the SageMaker domain IDs (reference
  `AwsSagemakerDomain` resources or pass literal IDs).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, artifact
store, and role, associate domains if wanted, and deploy.

### CLI

```bash
planton apply -f mlflow-app.yaml
```

## After Deploy

- `app_arn` is the app's AWS identity (all API operations key on it);
  `app_name` is the name, which updates in place.
- Everything updates in place except `role_arn` — the one
  replace-on-change field.
- There is no idle cost — the app bills per use, so it can sit ready
  without a meter running.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
