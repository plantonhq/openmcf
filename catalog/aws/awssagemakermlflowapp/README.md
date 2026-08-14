<p align="center">
  <img src="logo.svg" alt="AWS SageMaker MLflow App" width="80"/>
</p>

# AWS SageMaker MLflow App

Create and manage [Amazon SageMaker AI MLflow apps](https://docs.aws.amazon.com/sagemaker/latest/dg/mlflow.html)
— the serverless managed MLflow 3.x deployment that succeeds the
hourly-billed tracking server: no capacity to size, billed per use, $0
when idle.

## What Gets Created

- **An MLflow app** — serverless MLflow 3.x, storing artifacts in your
  S3 location (`artifact_store_uri`) through `role_arn`.
- **Domain associations** (`default_domain_ids`) — Studio users in
  those SageMaker domains track to this app automatically.
- Optional: account-default status (`account_default_status`; only one
  app per account can hold it), automatic model registration into the
  SageMaker Model Registry (`model_registration_mode`), and a weekly
  maintenance window.

## The ARN Is the Identity

All API operations key on the app's ARN — the name updates in place,
as does everything else except `role_arn`, the one replace-on-change
field. The app is standalone: it does NOT attach to a tracking server;
it associates with SageMaker domains. A soft-deleted app (status
DELETED) reads as absent.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
