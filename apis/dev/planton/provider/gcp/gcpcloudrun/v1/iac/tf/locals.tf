locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Service name defaults to metadata.name (the spec-level contract), so a
  # manifest only sets service_name when the cloud name must differ from the
  # Planton object name.
  service_name = var.spec.service_name != "" ? var.spec.service_name : var.metadata.name

  # Empty optional strings become null so the provider omits them from the
  # API payload instead of sending empty values it would reject or diff on.
  description     = var.spec.description != "" ? var.spec.description : null
  service_account = var.spec.service_account != "" ? var.spec.service_account : null
  encryption_key  = var.spec.encryption_key != "" ? var.spec.encryption_key : null
  revision        = var.spec.revision != "" ? var.spec.revision : null
  launch_stage    = var.spec.launch_stage != "" ? var.spec.launch_stage : null

  # Proto enum NAMES arrive as strings and match the provider's accepted
  # values exactly; the UNSPECIFIED zero value (or an empty string) means
  # "let GCP default" and must not reach the API.
  ingress = (
    var.spec.ingress != "" && var.spec.ingress != "INGRESS_TRAFFIC_UNSPECIFIED"
  ) ? var.spec.ingress : null
  execution_environment = (
    var.spec.execution_environment != "" && var.spec.execution_environment != "EXECUTION_ENVIRONMENT_UNSPECIFIED"
  ) ? var.spec.execution_environment : null

  # The API takes the timeout as a duration string ("300s"); the spec keeps
  # the honest integer-seconds shape.
  timeout = var.spec.timeout_seconds != null ? "${var.spec.timeout_seconds}s" : null

  # The same planton-ai_* label set the Pulumi module applies, so a service
  # is attributable to its Planton object regardless of the engine that
  # created it. User labels merge in first so the platform attribution
  # labels can never be clobbered by a spec label with the same key.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.service_name
    "planton-ai_kind"     = "gcpcloudrun"
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

  # vpc_access is optional end to end; these locals let the resource block
  # stay null-safe without repeating try() guards.
  vpc_connector = try(var.spec.vpc_access.connector, "") != "" ? var.spec.vpc_access.connector : null
  vpc_egress    = try(var.spec.vpc_access.egress, "") != "" ? var.spec.vpc_access.egress : null
  vpc_interfaces = try(var.spec.vpc_access.network_interfaces, [])
  has_vpc_access = var.spec.vpc_access != null
}
