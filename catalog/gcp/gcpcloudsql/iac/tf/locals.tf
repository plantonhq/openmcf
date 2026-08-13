locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Empty optional strings become null so the provider omits them from the
  # API payload instead of sending empty values it would reject or diff on.
  root_password        = var.spec.root_password != "" ? var.spec.root_password : null
  encryption_key_name  = var.spec.encryption_key_name != "" ? var.spec.encryption_key_name : null
  master_instance_name = var.spec.master_instance_name != "" ? var.spec.master_instance_name : null
  time_zone            = var.spec.time_zone != "" ? var.spec.time_zone : null
  collation            = var.spec.collation != "" ? var.spec.collation : null
  connector_enforcement = (
    var.spec.connector_enforcement != "" ? var.spec.connector_enforcement : null
  )

  # Disk defaults apply when the block is omitted entirely: 10 GB PD_SSD
  # with auto-resize — the same shape the variables' per-field defaults
  # produce for a present-but-sparse block, so both paths behave alike.
  disk_type              = var.spec.disk != null ? var.spec.disk.type : "PD_SSD"
  disk_size_gb           = var.spec.disk != null ? var.spec.disk.size_gb : 10
  disk_auto_resize       = var.spec.disk != null ? var.spec.disk.auto_resize : true
  disk_auto_resize_limit = var.spec.disk != null ? var.spec.disk.auto_resize_limit : 0

  # Connectivity contract: an omitted network block means public IPv4 with
  # no authorized networks — reachable only through the Auth Proxy or
  # connectors (IAM-authenticated). A present block states connectivity
  # explicitly (spec CEL guarantees at least one path is enabled).
  #
  # HCL's && does NOT short-circuit, so nested attributes on a nullable
  # block are read through try() — a bare `x != null && x.attr` would error
  # on the null block.
  ipv4_enabled = var.spec.network != null ? var.spec.network.ipv4_enabled : true
  private_network = (
    try(var.spec.network.private_network, "") != ""
    ? var.spec.network.private_network
    : null
  )
  allocated_ip_range = (
    try(var.spec.network.allocated_ip_range, "") != ""
    ? var.spec.network.allocated_ip_range
    : null
  )
  enable_private_path = try(var.spec.network.enable_private_path_for_google_cloud_services, false)
  ssl_mode = (
    try(var.spec.network.ssl_mode, "") != ""
    ? var.spec.network.ssl_mode
    : null
  )
  server_ca_mode = (
    try(var.spec.network.server_ca_mode, "") != ""
    ? var.spec.network.server_ca_mode
    : null
  )
  server_ca_pool = (
    try(var.spec.network.server_ca_pool, "") != ""
    ? var.spec.network.server_ca_pool
    : null
  )
  custom_sans         = try(var.spec.network.custom_subject_alternative_names, [])
  authorized_networks = try(var.spec.network.authorized_networks, [])
  psc = (
    try(var.spec.network.psc.enabled, false)
    ? var.spec.network.psc
    : null
  )

  server_certificate_rotation_mode = (
    try(var.spec.network.server_certificate_rotation_mode, "") != ""
    ? var.spec.network.server_certificate_rotation_mode
    : null
  )

  # The Enterprise Plus data cache block is emitted only when enabled so
  # ENTERPRISE instances carry no cache stanza at all (the API rejects it).
  data_cache_enabled = var.spec.data_cache_enabled

  backup_enabled = try(var.spec.backup.enabled, false)

  # Empty-string → null normalization for the new optional scalars.
  deletion_policy          = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null
  instance_type            = var.spec.instance_type != "" ? var.spec.instance_type : null
  backupdr_backup          = var.spec.backupdr_backup != "" ? var.spec.backupdr_backup : null
  maintenance_version      = var.spec.maintenance_version != "" ? var.spec.maintenance_version : null
  data_api_access          = var.spec.data_api_access != "" ? var.spec.data_api_access : null
  failover_dr_replica_name = var.spec.failover_dr_replica_name != "" ? var.spec.failover_dr_replica_name : null
  replica_names            = length(var.spec.replica_names) > 0 ? var.spec.replica_names : null
  final_backup_description = try(var.spec.final_backup.description, "") != "" ? var.spec.final_backup.description : null

  # Retention settings emit when either dial is present (the provider
  # defaults retention_unit to COUNT when only the count is given).
  backup_retention_settings = (
    try(var.spec.backup.retained_backups, null) != null || try(var.spec.backup.retention_unit, "") != ""
  ) ? [1] : []

  location_preference = (
    try(var.spec.location_preference.zone, "") != "" || try(var.spec.location_preference.secondary_zone, "") != ""
    ? var.spec.location_preference
    : null
  )

  # Convert the flags map to the provider's list-of-objects shape.
  database_flags_list = [
    for name, value in var.spec.database_flags : {
      name  = name
      value = value
    }
  ]

  connection_pool_flags_list = [
    for name, value in try(var.spec.connection_pooling.flags, {}) : {
      name  = name
      value = value
    }
  ]

  # The same planton-ai_* label set the Pulumi module applies, so a resource
  # is attributable to its Planton object regardless of the engine that
  # created it. Conditional labels appear under the same conditions on both
  # sides.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.spec.instance_name
    "planton-ai_kind"     = "gcpcloudsql"
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

  final_labels = merge(local.base_labels, local.org_label, local.env_label, local.id_label)
}
