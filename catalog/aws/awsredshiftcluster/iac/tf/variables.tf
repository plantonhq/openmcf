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
  description = "AwsRedshiftCluster specification"
  type = object({
    region = string
    subnet_ids = optional(list(string), [])
    cluster_subnet_group_name = optional(string, "")
    security_group_ids = optional(list(string), [])
    availability_zone = optional(string, "")
    availability_zone_relocation_enabled = optional(bool, false)
    publicly_accessible = optional(bool, false)
    elastic_ip = optional(string, "")
    enhanced_vpc_routing = optional(bool, false)
    port = optional(number, 0)
    node_type = string
    number_of_nodes = optional(number, 0)
    cluster_version = optional(string, "")
    database_name = optional(string, "")
    master_username = optional(string, "")
    manage_master_password = optional(bool, false)
    master_password = optional(string, "")
    master_password_secret_kms_key_id = optional(string, "")
    encrypted = optional(bool)
    kms_key_id = optional(string, "")
    multi_az = optional(bool, false)
    iam_roles = optional(list(string), [])
    default_iam_role_arn = optional(string, "")
    automated_snapshot_retention_period = optional(number)
    manual_snapshot_retention_period = optional(number, 0)
    preferred_maintenance_window = optional(string, "")
    maintenance_track_name = optional(string, "")
    allow_version_upgrade = optional(bool)
    apply_immediately = optional(bool, false)
    skip_final_snapshot = optional(bool, false)
    final_snapshot_identifier = optional(string, "")
    snapshot_identifier = optional(string, "")
    snapshot_arn = optional(string, "")
    snapshot_cluster_identifier = optional(string, "")
    owner_account = optional(string, "")
    logging = optional(object({
      log_destination_type = string
      s3_bucket_name = optional(string, "")
      s3_key_prefix = optional(string, "")
      log_exports = optional(list(string), [])
    }))
    snapshot_copy = optional(object({
      destination_region = string
      retention_period = optional(number, 0)
      manual_snapshot_retention_period = optional(number, 0)
      snapshot_copy_grant_name = optional(string, "")
    }))
    cluster_parameter_group_name = optional(string, "")
    parameters = optional(list(object({
      name = string
      value = string
    })), [])
    parameter_group_family = optional(string, "")
    snapshot_schedule_identifier = optional(string, "")
    usage_limits = optional(list(object({
      feature_type = string
      limit_type = string
      amount = optional(number, 0)
      period = optional(string, "")
      breach_action = optional(string, "")
    })), [])
    scheduled_actions = optional(list(object({
      name = string
      description = optional(string, "")
      disabled = optional(bool, false)
      schedule = string
      start_time = optional(string, "")
      end_time = optional(string, "")
      iam_role_arn = string
      pause_cluster = optional(bool, false)
      resume_cluster = optional(bool, false)
      resize_cluster = optional(object({
        classic = optional(bool, false)
        cluster_type = optional(string, "")
        node_type = optional(string, "")
        number_of_nodes = optional(number, 0)
      }))
    })), [])
    endpoint_accesses = optional(list(object({
      endpoint_name = string
      subnet_group_name = optional(string, "")
      vpc_security_group_ids = optional(list(string), [])
    })), [])
    endpoint_authorizations = optional(list(object({
      account = string
      vpc_ids = optional(list(string), [])
      force_delete = optional(bool, false)
    })), [])
  })
}
