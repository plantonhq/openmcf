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
  description = "AwsServerlessElasticache specification"
  type = object({
    region = string
    engine = string
    major_engine_version = optional(string, "")
    description = optional(string, "")
    data_storage_max_gb = optional(number, 0)
    data_storage_min_gb = optional(number, 0)
    ecpu_max = optional(number, 0)
    ecpu_min = optional(number, 0)
    subnet_ids = optional(list(string), [])
    security_group_ids = optional(list(string), [])
    network_type = optional(string, "")
    kms_key_id = optional(string, "")
    daily_snapshot_time = optional(string, "")
    snapshot_retention_limit = optional(number, 0)
    snapshot_arns_to_restore = optional(list(string), [])
    user_group_id = optional(string, "")
  })
}
