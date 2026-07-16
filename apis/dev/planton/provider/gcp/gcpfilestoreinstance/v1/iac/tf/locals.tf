locals {
  # Honor the spec contract: an empty project_id falls back to the
  # provider's default project.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # instance_name falls back to metadata.name — explicit conditional, so
  # both engines derive the identical cloud-side name.
  instance_name = var.spec.instance_name != "" ? var.spec.instance_name : var.metadata.name

  location = var.spec.location
  tier     = var.spec.tier

  # Empty optional strings become null so the provider applies its own
  # defaults (NFS_V3 protocol, DIRECT_PEERING connect mode, auto-picked
  # reserved range) instead of receiving an empty string it would reject.
  description  = var.spec.description != "" ? var.spec.description : null
  protocol     = var.spec.protocol != "" ? var.spec.protocol : null
  kms_key_name = var.spec.kms_key_name != "" ? var.spec.kms_key_name : null

  network           = var.spec.network_config.network
  connect_mode      = var.spec.network_config.connect_mode != "" ? var.spec.network_config.connect_mode : null
  reserved_ip_range = var.spec.network_config.reserved_ip_range != "" ? var.spec.network_config.reserved_ip_range : null

  # Empty modes follow the spec's documented default: IPv4 service.
  modes = length(var.spec.network_config.modes) > 0 ? var.spec.network_config.modes : ["MODE_IPV4"]

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.instance_name
    "planton-ai_kind"     = "gcpfilestoreinstance"
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
