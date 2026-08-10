locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The SLO display name defaults to metadata.name when the spec leaves
  # display_name empty — the same naming basis every kind uses.
  display_name = (
    var.spec.display_name != null && var.spec.display_name != ""
    ? var.spec.display_name
    : var.metadata.name
  )

  # Which service arm is set (exactly one — spec-validated).
  create_custom_service = var.spec.service.custom_service != null
  create_basic_service  = var.spec.service.basic_service != null

  # The service ID for a service this kind CREATES — the arm's service_id,
  # defaulting to metadata.name.
  created_service_id = (
    local.create_custom_service && var.spec.service.custom_service.service_id != ""
    ? var.spec.service.custom_service.service_id
    : local.create_basic_service && var.spec.service.basic_service.service_id != ""
    ? var.spec.service.basic_service.service_id
    : var.metadata.name
  )

  # The service the SLO measures: the existing service's ID, or the one the
  # module creates.
  slo_service_id = (
    local.create_custom_service
    ? google_monitoring_custom_service.this[0].service_id
    : local.create_basic_service
    ? google_monitoring_service.this[0].service_id
    : var.spec.service.service_id
  )

  # The created custom service's display name: the arm's own display_name,
  # else the kind's naming basis.
  custom_service_display_name = (
    local.create_custom_service && var.spec.service.custom_service.display_name != ""
    ? var.spec.service.custom_service.display_name
    : local.display_name
  )

  # The same planton-ai_* label set the Pulumi module applies, so a resource
  # is attributable to its Planton object regardless of the engine that
  # created it.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.metadata.name
    "planton-ai_kind"     = "gcpmonitoringslo"
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

  # User labels first: the platform labels win on key conflicts. Applied to
  # the SLO and to any service this kind creates.
  final_labels = merge(var.spec.labels, local.base_labels, local.org_label, local.env_label, local.id_label)
}
