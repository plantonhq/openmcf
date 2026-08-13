# Amazon SageMaker AI MLflow app: the serverless MLflow 3.x deployment
# (billed per use, nothing to size).
#
# Lifecycle facts the renders below depend on:
#   - the ARN is the app's identity (all API operations key on it); the
#     name updates in place;
#   - role_arn is the ONE replace-on-change argument
#     (provider-enforced);
#   - the app is standalone - it associates with SageMaker DOMAINS, not
#     with tracking servers;
#   - a soft-deleted app (status DELETED) reads as absent upstream.

resource "aws_sagemaker_mlflow_app" "this" {
  # The component's name IS the app name.
  name = local.app_name

  artifact_store_uri = var.spec.artifact_store_uri

  # The one replace-on-change argument.
  role_arn = var.spec.role_arn

  account_default_status = var.spec.account_default_status != "" ? var.spec.account_default_status : null

  default_domain_id_list = length(var.spec.default_domain_ids) > 0 ? var.spec.default_domain_ids : null

  model_registration_mode = var.spec.model_registration_mode != "" ? var.spec.model_registration_mode : null

  weekly_maintenance_window_start = var.spec.weekly_maintenance_window_start != "" ? var.spec.weekly_maintenance_window_start : null

  tags = local.aws_tags
}
