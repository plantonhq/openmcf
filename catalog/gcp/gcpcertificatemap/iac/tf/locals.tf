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

  # The map name defaults to metadata.name when the spec leaves map_name
  # empty — the same naming basis every kind uses.
  map_name = (
    var.spec.map_name != null && var.spec.map_name != ""
    ? var.spec.map_name
    : var.metadata.name
  )

  # The same planton-ai_* label set the Pulumi module applies, so a resource
  # is attributable to its Planton object regardless of the engine that
  # created it.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.metadata.name
    "planton-ai_kind"     = "gcpcertificatemap"
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
  # the map; entries layer this shared set over their own labels (so
  # platform attribution can never be shadowed) — identical merge order to
  # the Pulumi module.
  final_labels = merge(var.spec.labels, local.base_labels, local.org_label, local.env_label, local.id_label)

  # Empty deletion_policy defers to the provider default (DELETE). One spec
  # lever wired to the map and every entry.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  # The full resource name — assembled from the map's own computed project
  # attribute so the ambient-project case renders correctly.
  map_id = "projects/${google_certificate_manager_certificate_map.this.project}/locations/global/certificateMaps/${google_certificate_manager_certificate_map.this.name}"
}
