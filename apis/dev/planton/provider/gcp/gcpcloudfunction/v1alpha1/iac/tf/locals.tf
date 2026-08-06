locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Function name defaults to metadata.name.
  function_name = var.spec.function_name != "" ? var.spec.function_name : var.metadata.name

  # Platform attribution labels. User labels merge BENEATH them so the
  # platform's keys always win — the same order the Pulumi module applies.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.function_name
    "planton-ai_kind"     = "gcpcloudfunction"
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

  # Enum fields arrive from the converter as enum-name strings; empty means
  # "not set" and defers to the API default.
  is_http_trigger = var.spec.trigger == null || var.spec.trigger.trigger_type == "" || var.spec.trigger.trigger_type == "HTTP"

  service_config = var.spec.service_config

  # The Gen 2 API takes memory as a quantity string ("256M", "1Gi"); the
  # spec carries it verbatim. Empty defers to the API default (256M).
  available_memory = try(local.service_config.available_memory, "") != "" ? local.service_config.available_memory : null
  available_cpu    = try(local.service_config.available_cpu, "") != "" ? local.service_config.available_cpu : null

  vpc_connector = try(local.service_config.vpc_connector, "") != "" ? local.service_config.vpc_connector : null
  # Egress settings only make sense with a connector attached; sending them
  # without one is an API error.
  vpc_egress = local.vpc_connector != null && try(local.service_config.vpc_connector_egress_settings, "") != "" ? local.service_config.vpc_connector_egress_settings : null

  ingress_settings = try(local.service_config.ingress_settings, "") != "" ? local.service_config.ingress_settings : null

  build_update_policy = var.spec.build_config.update_policy
}
