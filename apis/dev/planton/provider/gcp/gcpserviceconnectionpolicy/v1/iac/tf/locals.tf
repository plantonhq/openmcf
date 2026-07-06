locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Policy name defaults to metadata.name when policy_name is omitted.
  policy_name = var.spec.policy_name != "" ? var.spec.policy_name : var.metadata.name

  # The Service Connectivity API rejects full https:// self-link URLs for
  # the network and subnets — it requires relative resource paths. A
  # GcpVpcNetwork reference already resolves to the resource path, but
  # literals (and GcpSubnetwork references, whose canonical output is the
  # self-link) may arrive as URLs, so both are normalized here. Stripping
  # is a no-op for values already in relative form — identical
  # normalization to the Pulumi module.
  network = replace(var.spec.network, "/^https://www\\.googleapis\\.com/compute/v1//", "")

  subnetworks = var.spec.psc_config != null ? [
    for s in var.spec.psc_config.subnetworks :
    replace(s, "/^https://www\\.googleapis\\.com/compute/v1//", "")
  ] : []

  description = var.spec.description != "" ? var.spec.description : null

  # The API types the connection limit as a string-encoded integer; 0 in
  # the spec means "leave GCP's default" and is translated to null.
  psc_limit = (
    var.spec.psc_config != null && try(var.spec.psc_config.limit, 0) > 0
  ) ? tostring(var.spec.psc_config.limit) : null

  producer_instance_location = (
    var.spec.psc_config != null && try(var.spec.psc_config.producer_instance_location, "") != ""
  ) ? var.spec.psc_config.producer_instance_location : null

  allowed_hierarchy_levels = (
    var.spec.psc_config != null && length(try(var.spec.psc_config.allowed_google_producers_resource_hierarchy_levels, [])) > 0
  ) ? var.spec.psc_config.allowed_google_producers_resource_hierarchy_levels : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.policy_name
    "planton-ai_kind"     = "gcpserviceconnectionpolicy"
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

  # User labels first so platform attribution labels win on key conflicts —
  # identical merge order to the Pulumi module.
  final_labels = merge(
    var.spec.labels,
    local.base_labels,
    local.org_label,
    local.env_label,
    local.id_label,
  )
}
