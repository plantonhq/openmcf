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
  description = "AwsRedisElasticache specification"
  type = object({
    region = string
    engine = optional(string, "")
    engine_version = optional(string, "")
    description = string
    node_type = optional(string, "")
    port = optional(number)
    num_cache_clusters = optional(number, 0)
    preferred_cache_cluster_azs = optional(list(string), [])
    num_node_groups = optional(number, 0)
    replicas_per_node_group = optional(number, 0)
    node_group_configurations = optional(list(object({
      node_group_id = optional(string, "")
      primary_availability_zone = optional(string, "")
      replica_availability_zones = optional(list(string), [])
      replica_count = optional(number, 0)
      slots = optional(string, "")
    })), [])
    automatic_failover_enabled = optional(bool, false)
    multi_az_enabled = optional(bool, false)
    durability = optional(string, "")
    global_replication_group_id = optional(string, "")
    subnet_ids = optional(list(string), [])
    subnet_group_name = optional(string, "")
    security_group_ids = optional(list(string), [])
    network_type = optional(string, "")
    ip_discovery = optional(string, "")
    at_rest_encryption_enabled = optional(bool, false)
    transit_encryption_enabled = optional(bool, false)
    transit_encryption_mode = optional(string, "")
    kms_key_id = optional(string, "")
    auth_token = optional(string, "")
    auth_token_update_strategy = optional(string, "")
    user_group_ids = optional(list(string), [])
    snapshot_arns = optional(list(string), [])
    snapshot_name = optional(string, "")
    maintenance_window = optional(string, "")
    snapshot_retention_limit = optional(number, 0)
    snapshot_window = optional(string, "")
    final_snapshot_identifier = optional(string, "")
    apply_immediately = optional(bool, false)
    parameter_group_family = optional(string, "")
    parameters = optional(list(object({
      name = string
      value = string
    })), [])
    parameter_group_name = optional(string, "")
    log_delivery_configurations = optional(list(object({
      destination_type = string
      destination = string
      log_format = string
      log_type = string
    })), [])
    notification_topic_arn = optional(string, "")
    auto_minor_version_upgrade = optional(bool, false)
    data_tiering_enabled = optional(bool, false)
    cluster_mode = optional(string, "")
  })
}
