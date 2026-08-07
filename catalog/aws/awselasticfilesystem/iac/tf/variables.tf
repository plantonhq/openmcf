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
  description = "AwsElasticFileSystem specification"
  type = object({
    region = string
    encrypted = optional(bool, false)
    kms_key_id = optional(string, "")
    performance_mode = optional(string, "")
    throughput_mode = optional(string, "")
    provisioned_throughput_in_mibps = optional(number, 0)
    availability_zone_name = optional(string, "")
    transition_to_ia = optional(string, "")
    transition_to_archive = optional(string, "")
    transition_to_primary_storage_class = optional(string, "")
    backup_enabled = optional(bool, false)
    replication_overwrite_protection = optional(string, "")
    mount_targets = list(object({
      subnet_id = string
      ip_address = optional(string, "")
      ip_address_type = optional(string, "")
      ipv6_address = optional(string, "")
    }))
    security_group_ids = optional(list(string), [])
    policy = optional(any)
    bypass_policy_lockout_safety_check = optional(bool, false)
    replication = optional(object({
      destination_region = optional(string, "")
      destination_availability_zone_name = optional(string, "")
      destination_kms_key_id = optional(string, "")
      destination_file_system_id = optional(string, "")
    }))
  })
}
