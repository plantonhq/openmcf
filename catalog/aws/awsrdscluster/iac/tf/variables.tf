variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsRdsCluster specification"
  type = object({
    region = string
    subnet_ids = optional(list(string), [])
    db_subnet_group_name = optional(string, "")
    security_group_ids = optional(list(string), [])
    availability_zones = optional(list(string), [])
    network_type = optional(string, "")
    port = optional(number, 0)
    engine = string
    engine_version = optional(string, "")
    engine_mode = optional(string, "")
    engine_lifecycle_support = optional(string, "")
    instances = optional(list(object({
      name = string
      instance_class = string
      promotion_tier = optional(number, 0)
      availability_zone = optional(string, "")
      publicly_accessible = optional(bool, false)
      db_parameter_group_name = optional(string, "")
      auto_minor_version_upgrade = optional(bool)
      performance_insights_enabled = optional(bool)
      monitoring_interval = optional(number, 0)
      ca_cert_identifier = optional(string, "")
    })), [])
    serverless_v2_scaling = optional(object({
      min_capacity = optional(number, 0)
      max_capacity = number
      seconds_until_auto_pause = optional(number, 0)
    }))
    serverless_v1_scaling = optional(object({
      auto_pause = optional(bool)
      min_capacity = optional(number, 0)
      max_capacity = optional(number, 0)
      seconds_until_auto_pause = optional(number, 0)
      seconds_before_timeout = optional(number, 0)
      timeout_action = optional(string, "")
    }))
    db_cluster_instance_class = optional(string, "")
    allocated_storage_gb = optional(number, 0)
    iops = optional(number, 0)
    storage_type = optional(string, "")
    database_name = optional(string, "")
    master_username = optional(string, "")
    manage_master_user_password = optional(bool, false)
    master_user_secret_kms_key_id = optional(string, "")
    master_password = optional(string, "")
    storage_encrypted = optional(bool, false)
    kms_key_id = optional(string, "")
    backup_retention_period = optional(number, 0)
    preferred_backup_window = optional(string, "")
    preferred_maintenance_window = optional(string, "")
    copy_tags_to_snapshot = optional(bool, false)
    delete_automated_backups = optional(bool)
    skip_final_snapshot = optional(bool, false)
    final_snapshot_identifier = optional(string, "")
    deletion_protection = optional(bool, false)
    backtrack_window_seconds = optional(number, 0)
    iam_database_authentication_enabled = optional(bool, false)
    iam_roles = optional(list(string), [])
    enable_http_endpoint = optional(bool, false)
    enabled_cloudwatch_logs_exports = optional(list(string), [])
    performance_insights_enabled = optional(bool, false)
    performance_insights_kms_key_id = optional(string, "")
    performance_insights_retention_period = optional(number, 0)
    monitoring_interval = optional(number, 0)
    monitoring_role_arn = optional(string, "")
    database_insights_mode = optional(string, "")
    snapshot_identifier = optional(string, "")
    restore_to_point_in_time = optional(object({
      source_cluster_identifier = optional(string, "")
      source_cluster_resource_id = optional(string, "")
      restore_to_time = optional(string, "")
      use_latest_restorable_time = optional(bool, false)
      restore_type = optional(string, "")
    }))
    replication_source_identifier = optional(string, "")
    source_region = optional(string, "")
    global_cluster_identifier = optional(string, "")
    enable_global_write_forwarding = optional(bool, false)
    enable_local_write_forwarding = optional(bool, false)
    db_cluster_parameter_group_name = optional(string, "")
    parameters = optional(list(object({
      name = string
      value = string
      apply_method = optional(string, "")
    })), [])
    db_instance_parameter_group_name = optional(string, "")
    ca_certificate_identifier = optional(string, "")
    apply_immediately = optional(bool, false)
    allow_major_version_upgrade = optional(bool, false)
  })
}
