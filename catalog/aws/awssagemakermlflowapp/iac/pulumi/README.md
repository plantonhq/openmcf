# AwsSagemakerMlflowApp — Pulumi module (Go)

Deploys an Amazon SageMaker AI MLflow app (`sagemaker.MlflowApp`) —
the serverless MLflow 3.x deployment.

Module facts worth knowing before editing:

- **The ARN is the app's identity** — all API operations key on it;
  the name updates in place.
- **`RoleArn` is the ONE replace-on-change argument**
  (provider-enforced); everything else updates in place.
- **The app is standalone** — it associates with SageMaker DOMAINS,
  not with tracking servers; the bridged field name
  (`DefaultDomainIdLists`) pluralizes the provider's
  `default_domain_id_list`.
- **A soft-deleted app (status DELETED) reads as absent** upstream —
  a delete that leaves the app visible in a raw listing still
  succeeded.

Outputs mirror the Terraform module key-for-key: `app_arn`, `app_name`.
