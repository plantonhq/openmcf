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

  # The cloud-side name defaults to metadata.name when the spec leaves
  # network_endpoint_group_name empty — the same naming basis every kind uses.
  network_endpoint_group_name = (
    var.spec.network_endpoint_group_name != null && var.spec.network_endpoint_group_name != ""
    ? var.spec.network_endpoint_group_name
    : var.metadata.name
  )

  # Normalize "" -> null for optional strings the provider treats as
  # meaningfully absent. network_endpoint_type has a GCP API default
  # (SERVERLESS) that matches the spec's proto default, so null and the
  # middleware-applied default are behaviorally identical.
  network_endpoint_type = var.spec.network_endpoint_type != "" ? var.spec.network_endpoint_type : null
  network               = var.spec.network != "" ? var.spec.network : null
  subnetwork            = var.spec.subnetwork != "" ? var.spec.subnetwork : null
  psc_target_service    = var.spec.psc_target_service != "" ? var.spec.psc_target_service : null

  # psc_data / serverless blocks: pass through only when present; "" -> null on
  # every leaf so the provider sees only the fields that were actually set.
  psc_data = var.spec.psc_data == null ? null : {
    producer_port = try(var.spec.psc_data.producer_port, "") != "" ? var.spec.psc_data.producer_port : null
  }

  cloud_run = var.spec.cloud_run == null ? null : {
    service  = try(var.spec.cloud_run.service, "") != "" ? var.spec.cloud_run.service : null
    tag      = try(var.spec.cloud_run.tag, "") != "" ? var.spec.cloud_run.tag : null
    url_mask = try(var.spec.cloud_run.url_mask, "") != "" ? var.spec.cloud_run.url_mask : null
  }

  cloud_function = var.spec.cloud_function == null ? null : {
    function = try(var.spec.cloud_function.function, "") != "" ? var.spec.cloud_function.function : null
    url_mask = try(var.spec.cloud_function.url_mask, "") != "" ? var.spec.cloud_function.url_mask : null
  }

  # App Engine may be an empty block (routes to the default app), so it passes
  # through whenever present even with all leaves null.
  app_engine = var.spec.app_engine == null ? null : {
    service  = try(var.spec.app_engine.service, "") != "" ? var.spec.app_engine.service : null
    version  = try(var.spec.app_engine.version, "") != "" ? var.spec.app_engine.version : null
    url_mask = try(var.spec.app_engine.url_mask, "") != "" ? var.spec.app_engine.url_mask : null
  }
}
