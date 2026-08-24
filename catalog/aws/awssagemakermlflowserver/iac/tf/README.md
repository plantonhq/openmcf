# AwsSagemakerMlflowServer — Terraform/OpenTofu module

Deploys an Amazon SageMaker AI MLflow tracking server
(`aws_sagemaker_mlflow_tracking_server`).

Module facts worth knowing before editing:

- **Create and delete each take ~25 minutes** — the provider's own
  timeouts are 45m per operation and are not user-configurable; budget
  accordingly.
- **The server bills hourly from Created onward** — size sets the rate.
- **`automatic_model_registration` cannot be turned back off** — a
  true-to-false change is silently not transmitted (an upstream
  update-guard gap taught on the spec field); the module always
  renders the spec value so the intent is visible in the plan.
- **`role_arn` and `mlflow_version` replace the server**
  (provider-enforced ForceNew); size resizes in place.
- **AWS normalizes `mlflow_version` to `major.minor`** — the spec's
  pattern already forbids patch-level values, so no drift.

Outputs mirror the Pulumi module key-for-key: `tracking_server_name`,
`tracking_server_arn`, `tracking_server_url`.
