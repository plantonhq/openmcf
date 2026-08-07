# The Azure SQL Database: the unit of compute and billing on its logical
# server. The azurerm resource internally orchestrates the TDE,
# short/long-term retention, and threat-detection ARM sub-APIs, so those
# are plain blocks here rather than separate resources.
resource "azurerm_mssql_database" "main" {
  name      = var.spec.database_name
  server_id = var.spec.server_id

  # The SKU (or "ElasticPool" when pooled). Null lets Azure compute its
  # default (serverless GP_S_Gen5_2).
  sku_name        = local.sku_name
  elastic_pool_id = var.spec.elastic_pool_id

  # Fractional sizes are legal ARM values (Basic tops out at 2 GB, S0 at
  # 250 GB); Hyperscale grows elastically and ignores the ceiling.
  max_size_gb = var.spec.max_size_gb

  collation    = var.spec.collation
  license_type = local.license_type

  # Serverless dials -- spec validation gates them to GP_S_/HS_S_ skus
  # before the plan ever runs.
  auto_pause_delay_in_minutes = var.spec.auto_pause_delay_in_minutes
  min_capacity                = var.spec.min_capacity

  # Availability: Hyperscale readable replicas, Premium/BC read scale-out,
  # and zone spreading.
  read_replica_count = var.spec.read_replica_count
  read_scale         = var.spec.read_scale
  zone_redundant     = var.spec.zone_redundant

  # Integrity and confidential computing. Changing enclave_type replaces
  # the database (ARM's contract).
  ledger_enabled = var.spec.ledger_enabled
  enclave_type   = local.enclave_type

  # Pooled databases inherit the pool's window -- spec validation forces
  # this null when elastic_pool_id is set.
  maintenance_configuration_name = var.spec.maintenance_configuration_name

  # Lifecycle: how the database comes into existence. Each mode consumes
  # its matching source; the pairings are spec-validated.
  create_mode                           = local.create_mode
  creation_source_database_id           = var.spec.creation_source_database_id
  secondary_type                        = local.secondary_type
  restore_point_in_time                 = var.spec.restore_point_in_time
  recover_database_id                   = var.spec.recover_database_id
  recovery_point_id                     = var.spec.recovery_point_id
  restore_dropped_database_id           = var.spec.restore_dropped_database_id
  restore_long_term_retention_backup_id = var.spec.restore_long_term_retention_backup_id

  # Backup redundancy and the DW-only geo-backup dial.
  storage_account_type = local.storage_account_type
  geo_backup_enabled   = var.spec.geo_backup_enabled

  sample_name = var.spec.sample_name

  # Database-scoped identities for the database-scoped CMK.
  dynamic "identity" {
    for_each = length(var.spec.user_assigned_identity_ids) > 0 ? [1] : []
    content {
      type         = "UserAssigned"
      identity_ids = var.spec.user_assigned_identity_ids
    }
  }

  # Transparent data encryption: the database-scoped CMK overrides the
  # server's key for this database; rotation re-encrypts automatically
  # when enabled. The rotation flag rides WITH the key -- the provider
  # requires the two set together (configuring the flag alone, even
  # false, is rejected), so it is null unless a CMK key exists.
  transparent_data_encryption_enabled                         = var.spec.transparent_data_encryption_enabled
  transparent_data_encryption_key_vault_key_id                = var.spec.transparent_data_encryption_key_vault_key_id
  transparent_data_encryption_key_automatic_rotation_enabled = var.spec.transparent_data_encryption_key_vault_key_id != null ? var.spec.transparent_data_encryption_key_automatic_rotation_enabled : null

  # A bacpac import applied right after creation (fresh databases only --
  # spec-validated).
  dynamic "import" {
    for_each = var.spec.import != null ? [var.spec.import] : []
    content {
      storage_uri                  = import.value.storage_uri
      storage_key                  = import.value.storage_key
      storage_key_type             = local.import_storage_key_type_map[import.value.storage_key_type]
      administrator_login          = import.value.administrator_login
      administrator_login_password = import.value.administrator_login_password
      authentication_type          = local.import_authentication_type_map[import.value.authentication_type]
      storage_account_id           = import.value.storage_account_id
    }
  }

  # The point-in-time-restore horizon and differential-backup cadence.
  dynamic "short_term_retention_policy" {
    for_each = var.spec.short_term_retention_policy != null ? [var.spec.short_term_retention_policy] : []
    content {
      retention_days           = short_term_retention_policy.value.retention_days
      backup_interval_in_hours = short_term_retention_policy.value.backup_interval_in_hours
    }
  }

  # Long-term full-backup retention. ARM's "PT0S" means a horizon keeps
  # nothing, so unset horizons are simply not sent.
  dynamic "long_term_retention_policy" {
    for_each = var.spec.long_term_retention_policy != null ? [var.spec.long_term_retention_policy] : []
    content {
      weekly_retention  = long_term_retention_policy.value.weekly_retention
      monthly_retention = long_term_retention_policy.value.monthly_retention
      yearly_retention  = long_term_retention_policy.value.yearly_retention
      week_of_year      = long_term_retention_policy.value.week_of_year
    }
  }

  # Database-scoped Microsoft Defender threat detection (overrides the
  # server-scope policy for this database).
  dynamic "threat_detection_policy" {
    for_each = var.spec.threat_detection_policy != null ? [var.spec.threat_detection_policy] : []
    content {
      state                      = local.threat_detection_state_map[threat_detection_policy.value.state]
      disabled_alerts            = threat_detection_policy.value.disabled_alerts
      email_account_admins       = threat_detection_policy.value.email_account_admins ? "Enabled" : "Disabled"
      email_addresses            = threat_detection_policy.value.email_addresses
      retention_days             = threat_detection_policy.value.retention_days
      storage_endpoint           = threat_detection_policy.value.storage_endpoint
      storage_account_access_key = threat_detection_policy.value.storage_account_access_key
    }
  }

  tags = local.final_tags
}
