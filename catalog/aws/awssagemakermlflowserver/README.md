<p align="center">
  <img src="logo.svg" alt="AWS SageMaker MLflow Server" width="80"/>
</p>

# AWS SageMaker MLflow Server

Create and manage [Amazon SageMaker AI MLflow tracking servers](https://docs.aws.amazon.com/sagemaker/latest/dg/mlflow.html)
— the classic managed MLflow deployment for experiments, runs, and
model tracking, running on dedicated capacity and billed per hour. For
the serverless successor, see `AwsSagemakerMlflowApp`.

## What Gets Created

- **A tracking server** sized `Small` (~25 users, ~$0.6/hour),
  `Medium` (~50), or `Large` (~100) — size resizes in place.
- **An artifact store** wiring: `artifact_store_uri` points at your S3
  location for model files and run outputs, accessed through
  `role_arn`.
- Optional: an `mlflow_version` pin as `major.minor` (AWS normalizes
  away the patch; changing the pin replaces the server), automatic
  model registration into the SageMaker Model Registry, and a weekly
  maintenance window.

## Operations Take ~25 Minutes

AWS provisions dedicated capacity: creation takes ~25 minutes and
deletion is similar (the provider's own timeouts are 45m per
operation), and the server bills hourly from Created onward. One more
trap: `automatic_model_registration` cannot be turned back OFF through
the provider — a true-to-false change is silently not transmitted.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
