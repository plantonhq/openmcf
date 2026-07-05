locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Job name defaults to metadata.name when job_name is omitted.
  job_name = var.spec.job_name != "" ? var.spec.job_name : var.metadata.name

  launch_stage   = var.spec.launch_stage != "" ? var.spec.launch_stage : null
  service_account = var.spec.template.service_account != "" ? var.spec.template.service_account : null
  encryption_key  = var.spec.template.encryption_key != "" ? var.spec.template.encryption_key : null

  execution_environment = (
    var.spec.template.execution_environment != "" && var.spec.template.execution_environment != "EXECUTION_ENVIRONMENT_UNSPECIFIED"
  ) ? var.spec.template.execution_environment : null

  # The API takes timeout as a duration string ("600s"); the spec keeps
  # the honest integer-seconds shape.
  timeout = var.spec.template.timeout_seconds != null ? "${var.spec.template.timeout_seconds}s" : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.job_name
    "planton-ai_kind"     = "gcpcloudrunjob"
  }

  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "planton-ai_organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "planton-ai_environment" = var.metadata.env } : {}

  id_label = (
    var.metadata.id != null && var.metadata.id != ""
  ) ? { "planton-ai_id" = var.metadata.id } : {}

  final_labels = merge(
    var.spec.labels,
    local.base_labels,
    local.org_label,
    local.env_label,
    local.id_label,
  )

  vpc_connector = try(var.spec.template.vpc_access.connector, "") != "" ? var.spec.template.vpc_access.connector : null
  vpc_egress    = try(var.spec.template.vpc_access.egress, "") != "" ? var.spec.template.vpc_access.egress : null
  vpc_interfaces = try(var.spec.template.vpc_access.network_interfaces, [])
  has_vpc_access = var.spec.template.vpc_access != null
}
