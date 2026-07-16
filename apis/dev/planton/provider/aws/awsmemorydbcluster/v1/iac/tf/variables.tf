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
  description = "AwsMemorydbCluster specification"
  type = object({
    region = string
    engine = string
    engine_version = optional(string, "")
    description = optional(string, "")
    node_type = string
    port = optional(number)
    num_shards = optional(number)
    num_replicas_per_shard = optional(number)
    acl_name = string
    subnet_ids = optional(list(string), [])
    subnet_group_name = optional(string, "")
    security_group_ids = optional(list(string), [])
    network_type = optional(string, "")
    ip_discovery = optional(string, "")
    tls_enabled = optional(bool)
    kms_key_arn = optional(string, "")
    maintenance_window = optional(string, "")
    snapshot_retention_limit = optional(number, 0)
    snapshot_window = optional(string, "")
    final_snapshot_name = optional(string, "")
    snapshot_arns = optional(list(string), [])
    snapshot_name = optional(string, "")
    parameter_group_family = optional(string, "")
    parameters = optional(list(object({
      name = string
      value = string
    })), [])
    parameter_group_name = optional(string, "")
    multi_region_cluster_name = optional(string, "")
    sns_topic_arn = optional(string, "")
    auto_minor_version_upgrade = optional(bool)
    data_tiering = optional(bool, false)
  })
}
