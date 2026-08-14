# Amazon SageMaker AI MLflow tracking server: the classic hourly-billed
# managed MLflow deployment.
#
# Lifecycle facts the renders below depend on:
#   - create and delete each take ~25 minutes (the provider's own
#     timeouts are 45m per operation, not user-configurable) - budget
#     accordingly;
#   - the server bills hourly from Created onward (Small ~$0.6/hour);
#   - automatic_model_registration CANNOT be turned back off through the
#     provider (a true-to-false change is silently not transmitted - an
#     upstream update-guard gap taught on the spec field); the module
#     always renders the spec value so the intent is visible in the
#     plan;
#   - AWS normalizes mlflow_version to major.minor (the spec's pattern
#     already forbids patch-level values, so no drift).

resource "aws_sagemaker_mlflow_tracking_server" "this" {
  # The component's name IS the tracking server name.
  tracking_server_name = local.tracking_server_name

  artifact_store_uri = var.spec.artifact_store_uri

  # Changing the role replaces the server (provider-enforced).
  role_arn = var.spec.role_arn

  tracking_server_size = var.spec.size != "" ? var.spec.size : null

  # Changing the version replaces the server (provider-enforced).
  mlflow_version = var.spec.mlflow_version != "" ? var.spec.mlflow_version : null

  automatic_model_registration = var.spec.automatic_model_registration

  weekly_maintenance_window_start = var.spec.weekly_maintenance_window_start != "" ? var.spec.weekly_maintenance_window_start : null

  tags = local.aws_tags
}
