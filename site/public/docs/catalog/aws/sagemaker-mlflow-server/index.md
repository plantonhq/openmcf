---
title: "SageMaker MLflow Server"
description: "SageMaker MLflow Server deployment documentation"
icon: "package"
order: 100
componentName: "awssagemakermlflowserver"
---

# AWS SageMaker MLflow Server

Managed MLflow as dedicated infrastructure — the classic tracking
server for experiments, runs, and model tracking, sized for your team,
storing artifacts in your S3 bucket, and billed per hour while it
runs. For the serverless successor, see `AwsSagemakerMlflowApp`.

## What Gets Created

- A tracking server sized `Small` (~25 users, ~$0.6/hour), `Medium`
  (~50), or `Large` (~100) — resizes in place.
- The S3 artifact store wiring (`artifact_store_uri`) accessed through
  the server's IAM role.
- Optional: an `mlflow_version` pin (`major.minor`), automatic model
  registration into the SageMaker Model Registry, and a weekly
  maintenance window.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker control-plane permissions
  (`sagemaker:CreateMlflowTrackingServer` and its siblings).

### AWS Account

- An S3 bucket (or prefix) for the artifact store.
- An IAM role trusting `sagemaker.amazonaws.com` with read/write
  access to that bucket (reference an `AwsIamRole` or pass the ARN).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, artifact
store, and role, choose a size, and deploy — then budget ~25 minutes
for the server to reach Created.

### CLI

```bash
planton apply -f mlflow-server.yaml
```

## After Deploy

- `tracking_server_name` / `tracking_server_arn` identify the server;
  `tracking_server_url` is the MLflow UI and the tracking URI your
  training code points at.
- Billing runs hourly from Created onward — delete servers you are not
  using (deletion also takes ~25 minutes).
- Turning `automatic_model_registration` on is effectively one-way
  through the provider — a true-to-false change is silently not
  transmitted; disabling requires replacing the server or an
  out-of-band API call.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
