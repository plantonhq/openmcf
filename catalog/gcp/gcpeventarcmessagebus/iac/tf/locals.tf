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

  # The bus ID defaults to metadata.name when the spec leaves
  # message_bus_id empty — the same naming basis every kind uses.
  message_bus_id = (
    var.spec.message_bus_id != null && var.spec.message_bus_id != ""
    ? var.spec.message_bus_id
    : var.metadata.name
  )

  # The same planton-ai_* label set the Pulumi module applies, so a resource
  # is attributable to its Planton object regardless of the engine that
  # created it.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.metadata.name
    "planton-ai_kind"     = "gcpeventarcmessagebus"
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
  # the bus; satellites layer this shared set over their own labels (so
  # platform attribution can never be shadowed) — identical merge order to
  # the Pulumi module.
  final_labels = merge(var.spec.labels, local.base_labels, local.org_label, local.env_label, local.id_label)

  # Empty deletion_policy defers to the provider default (DELETE). One spec
  # lever wired to the bus and every satellite.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null
}
