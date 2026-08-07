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
  description = "AwsFsxOpenzfsFileSystem specification"
  type = object({
    region = string
    deployment_type = optional(string)
    storage_capacity_gib = optional(number)
    storage_type = optional(string)
    throughput_capacity = optional(number, 0)
    read_cache_configuration = optional(object({
      sizing_mode = string
      size_gib = optional(number)
    }))
    subnet_ids = list(string)
    security_group_ids = optional(list(string), [])
    preferred_subnet_id = optional(string, "")
    endpoint_ip_address_range = optional(string, "")
    route_table_ids = optional(list(string), [])
    kms_key_id = optional(string, "")
    backup_id = optional(string, "")
    disk_iops_configuration = optional(object({
      mode = optional(string)
      iops = optional(number)
    }))
    root_volume_configuration = optional(object({
      data_compression_type = optional(string)
      nfs_exports = optional(object({
        client_configurations = list(object({
          clients = string
          options = list(string)
        }))
      }))
      read_only = optional(bool, false)
      record_size_kib = optional(number)
      user_and_group_quotas = optional(list(object({
        id = optional(number, 0)
        storage_capacity_quota_gib = optional(number, 0)
        type = string
      })), [])
      copy_tags_to_snapshots = optional(bool, false)
    }))
    automatic_backup_retention_days = optional(number)
    daily_automatic_backup_start_time = optional(string, "")
    copy_tags_to_backups = optional(bool, false)
    copy_tags_to_volumes = optional(bool, false)
    skip_final_backup = optional(bool)
    final_backup_tags = optional(map(string), {})
    delete_options = optional(list(string), [])
    weekly_maintenance_start_time = optional(string, "")
  })
}
