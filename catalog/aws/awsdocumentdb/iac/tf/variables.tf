variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name        = string
    id          = optional(string, "")
    org         = optional(string, "")
    env         = optional(string, "")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    tags        = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsDocumentDb specification"
  type = object({
    region               = string
    subnet_ids           = optional(list(string), [])
    db_subnet_group_name = optional(string, "")
    security_group_ids   = optional(list(string), [])
    availability_zones   = optional(list(string), [])
    network_type         = optional(string, "")
    port                 = optional(number, 0)
    engine_version       = optional(string, "")
    storage_type         = optional(string, "")
    instances = optional(list(object({
      name                            = string
      instance_class                  = string
      promotion_tier                  = optional(number, 0)
      availability_zone               = optional(string, "")
      auto_minor_version_upgrade      = optional(bool)
      performance_insights_enabled    = optional(bool, false)
      performance_insights_kms_key_id = optional(string, "")
      preferred_maintenance_window    = optional(string, "")
      ca_cert_identifier              = optional(string, "")
      copy_tags_to_snapshot           = optional(bool, false)
      certificate_rotation_restart    = optional(bool)
      apply_immediately               = optional(bool, false)
    })), [])
    serverless_v2_scaling = optional(object({
      min_capacity = number
      max_capacity = number
    }))
    master_username                 = optional(string, "")
    manage_master_user_password     = optional(bool, false)
    master_password                 = optional(string, "")
    storage_encrypted               = optional(bool, false)
    kms_key_id                      = optional(string, "")
    backup_retention_period         = optional(number, 0)
    preferred_backup_window         = optional(string, "")
    preferred_maintenance_window    = optional(string, "")
    skip_final_snapshot             = optional(bool, false)
    final_snapshot_identifier       = optional(string, "")
    deletion_protection             = optional(bool, false)
    enabled_cloudwatch_logs_exports = optional(list(string), [])
    snapshot_identifier             = optional(string, "")
    restore_to_point_in_time = optional(object({
      source_cluster_identifier  = string
      restore_to_time            = optional(string, "")
      use_latest_restorable_time = optional(bool, false)
      restore_type               = optional(string, "")
    }))
    global_cluster_identifier       = optional(string, "")
    db_cluster_parameter_group_name = optional(string, "")
    parameters = optional(list(object({
      name         = string
      value        = string
      apply_method = optional(string, "")
    })), [])
    apply_immediately           = optional(bool, false)
    allow_major_version_upgrade = optional(bool, false)
  })
}