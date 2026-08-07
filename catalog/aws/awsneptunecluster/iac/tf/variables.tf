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
  description = "AwsNeptuneCluster specification"
  type = object({
    region = string
    subnet_ids = optional(list(string), [])
    neptune_subnet_group_name = optional(string, "")
    security_group_ids = optional(list(string), [])
    availability_zones = optional(list(string), [])
    port = optional(number, 0)
    engine_version = optional(string, "")
    storage_type = optional(string, "")
    instances = optional(list(object({
      name = string
      instance_class = string
      promotion_tier = optional(number, 0)
      availability_zone = optional(string, "")
      publicly_accessible = optional(bool, false)
      neptune_parameter_group_name = optional(string, "")
      auto_minor_version_upgrade = optional(bool)
      preferred_maintenance_window = optional(string, "")
    })), [])
    serverless_v2_scaling = optional(object({
      min_capacity = number
      max_capacity = number
    }))
    storage_encrypted = optional(bool, false)
    kms_key_id = optional(string, "")
    iam_database_authentication_enabled = optional(bool, false)
    iam_roles = optional(list(string), [])
    backup_retention_period = optional(number, 0)
    preferred_backup_window = optional(string, "")
    preferred_maintenance_window = optional(string, "")
    copy_tags_to_snapshot = optional(bool, false)
    skip_final_snapshot = optional(bool, false)
    final_snapshot_identifier = optional(string, "")
    deletion_protection = optional(bool, false)
    enabled_cloudwatch_logs_exports = optional(list(string), [])
    snapshot_identifier = optional(string, "")
    replication_source_identifier = optional(string, "")
    global_cluster_identifier = optional(string, "")
    neptune_cluster_parameter_group_name = optional(string, "")
    parameters = optional(list(object({
      name = string
      value = string
      apply_method = optional(string, "")
    })), [])
    neptune_instance_parameter_group_name = optional(string, "")
    apply_immediately = optional(bool, false)
    allow_major_version_upgrade = optional(bool, false)
  })
}
