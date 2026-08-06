resource "aws_mwaa_environment" "this" {
  # The environment name is metadata.name -- create-only in AWS (ForceNew),
  # and the basis both engines share so a manifest deploys identically on
  # either.
  name               = local.environment_name
  dag_s3_path        = var.spec.dag_s3_path
  execution_role_arn = var.spec.execution_role_arn
  source_bucket_arn  = var.spec.source_bucket_arn
  tags               = local.aws_tags

  # Ingress is composed, never embedded: the referenced security groups carry
  # the self-referencing all-traffic rule MWAA components need to talk to each
  # other, plus HTTPS (443) ingress for UI access. Subnets are ForceNew; the
  # attached security groups are the one part of network_configuration AWS
  # allows changing in place.
  network_configuration {
    subnet_ids         = var.spec.subnet_ids
    security_group_ids = var.spec.security_group_ids
  }

  # Airflow version: omitted means AWS picks the latest supported release.
  airflow_version = var.spec.airflow_version != "" ? var.spec.airflow_version : null

  # Airflow config overrides ("section.property" keys). May carry secrets
  # (connection URIs) -- AWS marks the whole map sensitive.
  airflow_configuration_options = length(var.spec.airflow_configuration_options) > 0 ? var.spec.airflow_configuration_options : null

  # S3 artifacts: each optional path can be pinned to an S3 object version for
  # deterministic deployments.
  plugins_s3_path                  = var.spec.plugins_s3_path != "" ? var.spec.plugins_s3_path : null
  plugins_s3_object_version        = var.spec.plugins_s3_object_version != "" ? var.spec.plugins_s3_object_version : null
  requirements_s3_path             = var.spec.requirements_s3_path != "" ? var.spec.requirements_s3_path : null
  requirements_s3_object_version   = var.spec.requirements_s3_object_version != "" ? var.spec.requirements_s3_object_version : null
  startup_script_s3_path           = var.spec.startup_script_s3_path != "" ? var.spec.startup_script_s3_path : null
  startup_script_s3_object_version = var.spec.startup_script_s3_object_version != "" ? var.spec.startup_script_s3_object_version : null

  # KMS encryption at rest (metadata DB, logs, SQS). ForceNew.
  kms_key = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  # Environment sizing. Zero means "not set" for these optional numerics --
  # AWS then applies its own defaults (class mw1.small, 1-10 workers, 2
  # webservers, 2 schedulers).
  environment_class = var.spec.environment_class != "" ? var.spec.environment_class : null
  min_workers       = var.spec.min_workers > 0 ? var.spec.min_workers : null
  max_workers       = var.spec.max_workers > 0 ? var.spec.max_workers : null
  min_webservers    = var.spec.min_webservers > 0 ? var.spec.min_webservers : null
  max_webservers    = var.spec.max_webservers > 0 ? var.spec.max_webservers : null
  schedulers        = var.spec.schedulers > 0 ? var.spec.schedulers : null

  # Access and networking. endpoint_management is ForceNew; with CUSTOMER the
  # operator creates VPC endpoints against the exported endpoint service names.
  webserver_access_mode = var.spec.webserver_access_mode
  endpoint_management   = var.spec.endpoint_management != "" ? var.spec.endpoint_management : null

  # Maintenance
  weekly_maintenance_window_start = var.spec.weekly_maintenance_window_start != "" ? var.spec.weekly_maintenance_window_start : null

  # Worker replacement strategy during environment updates (FORCED = fast,
  # may interrupt running tasks; GRACEFUL = drain first).
  worker_replacement_strategy = var.spec.worker_replacement_strategy != "" ? var.spec.worker_replacement_strategy : null

  # Per-module log delivery to CloudWatch. Each of the five Airflow components
  # is independently switchable with its own level.
  dynamic "logging_configuration" {
    for_each = var.spec.logging_configuration != null ? [var.spec.logging_configuration] : []
    content {
      dynamic "dag_processing_logs" {
        for_each = logging_configuration.value.dag_processing_logs != null ? [logging_configuration.value.dag_processing_logs] : []
        content {
          enabled   = dag_processing_logs.value.enabled
          log_level = dag_processing_logs.value.log_level
        }
      }

      dynamic "scheduler_logs" {
        for_each = logging_configuration.value.scheduler_logs != null ? [logging_configuration.value.scheduler_logs] : []
        content {
          enabled   = scheduler_logs.value.enabled
          log_level = scheduler_logs.value.log_level
        }
      }

      dynamic "task_logs" {
        for_each = logging_configuration.value.task_logs != null ? [logging_configuration.value.task_logs] : []
        content {
          enabled   = task_logs.value.enabled
          log_level = task_logs.value.log_level
        }
      }

      dynamic "webserver_logs" {
        for_each = logging_configuration.value.webserver_logs != null ? [logging_configuration.value.webserver_logs] : []
        content {
          enabled   = webserver_logs.value.enabled
          log_level = webserver_logs.value.log_level
        }
      }

      dynamic "worker_logs" {
        for_each = logging_configuration.value.worker_logs != null ? [logging_configuration.value.worker_logs] : []
        content {
          enabled   = worker_logs.value.enabled
          log_level = worker_logs.value.log_level
        }
      }
    }
  }
}
