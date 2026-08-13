locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Empty optional strings become null so the provider omits them from the
  # API payload instead of sending empty values it would reject or diff on.
  redis_version           = var.spec.redis_version != "" ? var.spec.redis_version : null
  display_name            = var.spec.display_name != "" ? var.spec.display_name : null
  location_id             = var.spec.location_id != "" ? var.spec.location_id : null
  alternative_location_id = var.spec.alternative_location_id != "" ? var.spec.alternative_location_id : null
  authorized_network      = var.spec.authorized_network != "" ? var.spec.authorized_network : null
  connect_mode            = var.spec.connect_mode != "" ? var.spec.connect_mode : null
  reserved_ip_range       = var.spec.reserved_ip_range != "" ? var.spec.reserved_ip_range : null
  secondary_ip_range      = var.spec.secondary_ip_range != "" ? var.spec.secondary_ip_range : null
  transit_encryption_mode = var.spec.transit_encryption_mode != "" ? var.spec.transit_encryption_mode : null
  maintenance_version     = var.spec.maintenance_version != "" ? var.spec.maintenance_version : null
  read_replicas_mode      = var.spec.read_replicas_mode != "" ? var.spec.read_replicas_mode : null
  replica_count           = var.spec.replica_count > 0 ? var.spec.replica_count : null
  customer_managed_key    = var.spec.customer_managed_key != "" ? var.spec.customer_managed_key : null

  # The same planton-ai_* label set the Pulumi module applies, so an instance
  # is attributable to its Planton object regardless of the engine that
  # created it. The instance name (not metadata.name) keys the name label so
  # the label matches what is visible in the GCP console.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.spec.instance_name
    "planton-ai_kind"     = "gcpredisinstance"
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
  # a user label can never mask the platform's ownership markers.
  final_labels = merge(
    var.spec.labels,
    local.base_labels,
    local.org_label,
    local.env_label,
    local.id_label,
  )
}
