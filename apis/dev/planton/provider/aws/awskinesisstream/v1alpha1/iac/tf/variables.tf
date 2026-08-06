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
  description = "AwsKinesisStream specification"
  type = object({
    region = string
    stream_mode = optional(string, "")
    shard_count = optional(number, 0)
    retention_period_hours = optional(number, 0)
    kms_key_id = optional(string, "")
    max_record_size_in_kib = optional(number, 0)
    shard_level_metrics = optional(list(string), [])
    enforce_consumer_deletion = optional(bool, false)
    warm_throughput_mib_ps = optional(number, 0)
    resource_policy = optional(any)
  })
}
