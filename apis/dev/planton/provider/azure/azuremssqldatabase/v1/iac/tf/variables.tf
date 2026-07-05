variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure SQL Database specification"
  type = object({
    # The logical server the database lives on, by resolved ARM ID.
    server_id = string

    # The database name, unique within the server.
    database_name = string

    # The compute SKU (DTU, vCore, serverless, Hyperscale, DW/DS,
    # ElasticPool, or Free). Empty lets Azure apply its default.
    sku_name = optional(string)

    # The elastic pool the database joins, by resolved ARM ID. Requires
    # sku_name = "ElasticPool".
    elastic_pool_id = optional(string)

    # The storage ceiling in GB (fractional sizes are legal ARM values).
    max_size_gb = optional(number)

    # The database collation; fixed at creation.
    collation = optional(string, "SQL_Latin1_General_CP1_CI_AS")

    # Azure Hybrid Benefit, as the spec enum's name string
    # (BASE_PRICE/LICENSE_INCLUDED). vCore tiers only.
    license_type = optional(string)

    # Serverless dials (GP_S_/HS_S_ skus only): minutes before
    # auto-pause (-1 disables) and the always-warm vCore floor.
    auto_pause_delay_in_minutes = optional(number)
    min_capacity                = optional(number)

    # Hyperscale readable HA replicas (0-4).
    read_replica_count = optional(number)

    # Premium/Business Critical: route read-intent connections to a
    # readable secondary.
    read_scale = optional(bool, false)

    # Spread replicas across availability zones.
    zone_redundant = optional(bool, false)

    # Cryptographically verifiable (tamper-evident) tables; fixed at
    # creation.
    ledger_enabled = optional(bool, false)

    # The confidential-computing enclave, as the spec enum's name string
    # (VBS/DEFAULT_ENCLAVE).
    enclave_type = optional(string)

    # The maintenance window (e.g. SQL_Default, SQL_EastUS_DB_1). Unset
    # on pooled databases (they inherit the pool's).
    maintenance_configuration_name = optional(string)

    # How the database comes into existence, as the spec enum's name
    # string; unset means DEFAULT.
    create_mode = optional(string)

    # Mode-paired sources: the source database for
    # COPY/SECONDARY/ONLINE_SECONDARY/POINT_IN_TIME_RESTORE, the restore
    # timestamp for POINT_IN_TIME_RESTORE, and the recovery/restore ids
    # for RECOVERY/RESTORE/RESTORE_LONG_TERM_RETENTION_BACKUP.
    creation_source_database_id           = optional(string)
    secondary_type                        = optional(string)
    restore_point_in_time                 = optional(string)
    recover_database_id                   = optional(string)
    recovery_point_id                     = optional(string)
    restore_dropped_database_id           = optional(string)
    restore_long_term_retention_backup_id = optional(string)

    # Backup storage redundancy, as the spec enum's name string
    # (GEO_REDUNDANT/GEO_ZONE_REDUNDANT/LOCALLY_REDUNDANT/ZONE_REDUNDANT).
    storage_account_type = optional(string)

    # DW SKUs only: geo-redundant backups.
    geo_backup_enabled = optional(bool, true)

    # Seed with a sample schema ("AdventureWorksLT").
    sample_name = optional(string)

    # User-assigned identities attached to the database (database-scoped
    # CMK unwraps through them), by resolved ARM ID.
    user_assigned_identity_ids = optional(list(string), [])

    # Transparent data encryption: the on/off dial (Azure default true),
    # the database-scoped CMK (VERSIONED Key Vault key id), and automatic
    # re-encryption on key rotation.
    transparent_data_encryption_enabled                         = optional(bool, true)
    transparent_data_encryption_key_vault_key_id                = optional(string)
    transparent_data_encryption_key_automatic_rotation_enabled = optional(bool, false)

    # A bacpac import applied right after creation. storage_key_type and
    # authentication_type arrive as spec enum name strings.
    import = optional(object({
      storage_uri                  = string
      storage_key                  = string
      storage_key_type             = string
      administrator_login          = string
      administrator_login_password = string
      authentication_type          = string
      storage_account_id           = optional(string)
    }))

    # The point-in-time-restore horizon (1-35 days) and differential
    # backup cadence (12 or 24 hours).
    short_term_retention_policy = optional(object({
      retention_days           = number
      backup_interval_in_hours = optional(number)
    }))

    # Long-term full-backup retention, each horizon an ISO-8601 duration.
    long_term_retention_policy = optional(object({
      weekly_retention  = optional(string)
      monthly_retention = optional(string)
      yearly_retention  = optional(string)
      week_of_year      = optional(number)
    }))

    # Database-scoped Microsoft Defender threat detection. state arrives
    # as the spec enum's name string; disabled_alerts carry ARM's wire
    # vocabulary verbatim (e.g. Sql_Injection).
    threat_detection_policy = optional(object({
      state                      = string
      disabled_alerts            = optional(list(string), [])
      email_account_admins       = optional(bool, false)
      email_addresses            = optional(list(string), [])
      retention_days             = optional(number, 0)
      storage_endpoint           = optional(string)
      storage_account_access_key = optional(string)
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
